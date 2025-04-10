// Package main contains the server configuration and initialization code.
package main

import (
	"log"
	"net/http"
)

// startServer initializes and starts the HTTP server for the Gratitude Jar application.
// It performs the following tasks:
// 1. Sets up the HTTP router with all application routes
// 2. Configures the server to listen on port 4000
// 3. Handles server startup errors
//
// Note: This function is currently not used as the server initialization
// is handled directly in main.go. It's kept here for reference and potential
// future use if server configuration becomes more complex.
func startServer() {
	// Initialize the HTTP router with all application routes
	// The routes function is defined in routes.go
	mux := routes()

	// Start the HTTP server on port 4000
	// The server will handle all incoming HTTP requests using the configured routes
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		// Log the detailed error message for debugging
		log.Printf("Server error: %v", err)
		// Terminate the application if the server fails to start
		log.Fatal(err)
	}
}
