// Package main is the entry point for the Gratitude Jar application.
// It initializes the database connection, sets up routes, and starts the HTTP server.
package main

import (
	"log"
	"net/http"

	"github.com/darynforman/gratitude-jar/internal/config"
)

// main is the entry point of the application.
// It performs the following tasks:
// 1. Initializes the database connection
// 2. Sets up the HTTP routes
// 3. Starts the HTTP server on port 4000
func main() {
	// Initialize database connection using the configuration from internal/config
	// If initialization fails, the application will terminate with a fatal error
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize the HTTP router with all application routes
	// The routes function is defined in routes.go
	mux := routes()

	// Start the HTTP server on port 4000
	// The server will handle all incoming HTTP requests using the configured routes
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
