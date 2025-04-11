// Package main is the entry point for the Gratitude Jar application.
// It initializes the database connection, sets up routes, and starts the HTTP server.
package main

import (
	"log"

	"github.com/darynforman/gratitude-jar/internal/config"
)

// main is the entry point of the application.
// It performs the following tasks:
// 1. Initializes the database connection
// 2. Starts the HTTP server using startServer()
func main() {
	// Initialize database connection using the configuration from internal/config
	// If initialization fails, the application will terminate with a fatal error
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start the server using the startServer function from server.go
	startServer()
}
