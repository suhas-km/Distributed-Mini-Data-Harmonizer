package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/config"
	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/handler"
	"github.com/suhas-km/distributed-mini-data-harmonizer/go-worker/internal/worker"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configPath != "" {
		// Load from file if provided
		cfg, err = config.LoadConfigFromFile(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config from file: %v", err)
		}
	} else {
		// Load from environment variables
		cfg, err = config.LoadConfig()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	// Create worker pool
	pool := worker.NewPool(cfg.WorkerCount, cfg.QueueSize)

	// Start worker pool
	pool.Start()
	defer pool.Stop()

	// Create HTTP handler
	jobHandler := handler.NewJobHandler(cfg, pool)

	// Create HTTP server
	mux := http.NewServeMux()
	jobHandler.RegisterRoutes(mux)

	// Add basic logging middleware
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log.Printf("Request: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
			log.Printf("Response: %s %s - %v", r.Method, r.URL.Path, time.Since(start))
		})
	}

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: loggingMiddleware(mux),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%d", cfg.Host, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	
	// Graceful shutdown
	pool.Stop()
	log.Println("Worker pool stopped")
	
	log.Println("Server shutdown complete")
}
