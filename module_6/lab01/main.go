package main

import (
	"fmt"
	"load-balancer/configs"
	"load-balancer/configs/database"
	"load-balancer/helpers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

var (
	serverPool ServerPool
	logger     *log.Logger
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

type ServerPool struct {
	backends []*Backend
	current  uint64
	mutex    sync.RWMutex
}

// AddBackend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.backends = append(s.backends, backend)
}

func (s *ServerPool) GetNextBackend() *Backend {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.current++
	next := s.current % uint64(len(s.backends))
	return s.backends[next]
}

func main() {
	c, err := configs.LoadConfig()
	database.InitializeRedis(c)

	if err != nil {
		log.Printf("Load env error: %v\n", err)
		return
	}

	if len(c.Servers) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	if c.LogEnabled {
		logOutput, err := helpers.LogFile(c.LogFile)
		if err != nil {
			log.Fatalf("Failed to setup log file: %v", err)
		}
		logger = log.New(logOutput, "", log.LstdFlags)
	}

	serverPool = ServerPool{}
	lb := configs.NewLoadBalancerServer(c)

	for _, s := range c.Servers {
		serverUrl, err := url.Parse(s)

		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverUrl)
		serverPool.AddBackend(&Backend{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}

	server := http.Server{
		Addr: fmt.Sprintf(":%d", c.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			backend := serverPool.GetNextBackend()

			if backend == nil {
				http.Error(w, "Service not available", http.StatusServiceUnavailable)
				return
			}

			ip := r.RemoteAddr
			serverPool.mutex.Lock()
			err := lb.HandleRequest(ip)
			serverPool.mutex.Unlock()

			if err != nil {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
				return
			}

			start := time.Now()
			backend.ReverseProxy.ServeHTTP(w, r)
			elapsed := time.Since(start)

			if logger != nil {
				logger.Printf("[%s] %s %s - %dms", r.Method, r.Host, r.URL.Path, elapsed.Milliseconds())
			}
		}),
	}

	log.Printf("Load Balancer started at :%d\n", c.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
