package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bharat1Rajput/HoldUp/internal/config"
	"github.com/Bharat1Rajput/HoldUp/internal/handlers"
	"github.com/Bharat1Rajput/HoldUp/internal/middleware"
	"github.com/Bharat1Rajput/HoldUp/internal/ratelimiter"
)

func main() {

	cfg := config.Load()

	// Initialize rate limiter
	limiter := ratelimiter.NewLimiter(
		cfg.RateLimitCapacity,
		cfg.RateLimitRefillRate,
	)

	mux := http.NewServeMux()

	// Public endpoint (no rate limit)
	mux.HandleFunc("/health", handlers.HealthHandler)

	// Rate-limited endpoint
	rateLimitedResource := middleware.RateLimitMiddleware(limiter)(
		http.HandlerFunc(handlers.ResourceHandler),
	)

	mux.Handle("/api/resource", rateLimitedResource)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("server is running on port %s\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	log.Println("--Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
