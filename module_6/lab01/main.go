package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"load-balancer/configs"
	"load-balancer/db"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
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
	backends    []*Backend
	current     uint64
	mutex       sync.RWMutex
	redisClient *redis.Client
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
	db.InitializeRedis(c)

	if err != nil {
		log.Printf("Load env error: %v\n", err)
		return
	}

	var port int
	flag.IntVar(&port, "port", 3000, "Port to serve")
	flag.Parse()

	if len(c.Servers) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	if c.LogEnabled {
		logOutput, err := setupLogFile(c.LogFile)
		if err != nil {
			log.Fatalf("Failed to setup log file: %v", err)
		}
		logger = log.New(logOutput, "", log.LstdFlags)
	}

	fmt.Println(c.Servers)

	serverPool = ServerPool{
		redisClient: db.RedisClient,
	}

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
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			backend := serverPool.GetNextBackend()

			if backend != nil {
				if c.RateLimit > 0 {
					serverPool.mutex.Lock()
					defer serverPool.mutex.Unlock()

					ip := r.RemoteAddr
					key := string(fmt.Sprintf("caches:ips:%s", ip))
					value, err := serverPool.redisClient.Get(context.Background(), key).Bytes()

					var data struct {
						Counters  int    `json:"counters"`
						TimeLimit string `json:"time_limit"`
					}

					json.Unmarshal(value, &data)

					if err != nil || data.TimeLimit < time.Now().String() {
						data.Counters = 0
						data.TimeLimit = time.Now().Add(time.Second * 60).String()
					}
					data.Counters++

					if data.Counters > c.RateLimit && data.TimeLimit > time.Now().String() {
						http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
						return
					} else {
						jsonData, err := json.Marshal(&data)
						if err != nil {
							panic(err)
						}

						err = serverPool.redisClient.SetEx(context.Background(), key, jsonData, time.Minute*10).Err()

						if err != nil {
							panic(err)
						}
					}
				}

				start := time.Now()
				backend.ReverseProxy.ServeHTTP(w, r)
				elapsed := time.Since(start)

				if logger != nil {
					logger.Printf("[%s] %s %s - %dms", r.Method, r.Host, r.URL.Path, elapsed.Milliseconds())
				}

				return
			}

			http.Error(w, "Service not available", http.StatusServiceUnavailable)
		}),
	}

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func setupLogFile(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}
