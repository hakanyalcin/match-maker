package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"matchmaking-httpapi/pkg/api"
	"matchmaking-httpapi/pkg/matchmaker"
	"matchmaking-httpapi/pkg/metrics"

	"github.com/gorilla/mux"
)

const (
	serverPort = ":8080"
	shutdownTimeout = 15 * time.Second
)

func main() {
	// Set up logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting matchmaking service...")

	// Create components
	mm := matchmaker.NewMatchmaker()
	met := metrics.NewMetrics()
	handler := api.NewHandler(mm, met)

	// Set up router
	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	
	// Add middleware for request logging
	router.Use(loggingMiddleware)

	// Create server
	srv := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server listening on %s", serverPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// Middleware for logging HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		duration := time.Since(start)
		log.Printf("Response: %s %s - took %v", r.Method, r.URL.Path, duration)
	})
} 