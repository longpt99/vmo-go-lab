package main

import (
	"album-manager/src/configs"
	"album-manager/src/configs/database"
	"album-manager/src/configs/repository"
	"album-manager/src/configs/router"
	"album-manager/src/utils/validate"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := configs.InitConfig()
	if err != nil {
		log.Fatalf("InitConfig error occurred. Err: %s", err)
	}

	err = validate.RegisterValidation()
	if err != nil {
		log.Fatalf("RegisterValidation error occurred. Err: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load database
	store, err := database.InitDatabase(ctx)
	if err != nil {
		log.Printf("InitDatabase error occurred. Err: %s", err)
	} else {
		store.InitializeFunction()
	}
	defer store.Close()

	repo := repository.InitRepositories(store)
	p := configs.Env.Port
	r := router.New(repo)
	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", p),
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	// Create a channel to receive signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start the server in a separate goroutine
	go func() {
		log.Printf("Server is listening on %s\n", s.Addr)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for a signal to shutdown the server
	sig := <-signalCh
	log.Printf("[Server] Received signal: %v\n", sig)

	// Shutdown the server gracefully
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("[Server] Shutdown Failed: %v\n", err)
		return
	}

	log.Println("[Server] Shutdown Gracefully! 🚀")
}
