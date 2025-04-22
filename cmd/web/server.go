// Package main contains the server configuration and initialization code.
package main

import (
	"log"
	"net/http"

	"github.com/darynforman/gratitude-jar1/internal/session"
)

// startServer initializes and starts the HTTP server for the Gratitude Jar application.
// It performs the following tasks:
// 1. Initializes the template cache
// 2. Sets up the HTTP router with all application routes
// 3. Configures the server to listen on port 4000
// 4. Handles server startup errors
//
// Note: This function is currently not used as the server initialization
// is handled directly in main.go. It's kept here for reference and potential
// future use if server configuration becomes more complex.
func startServer(app *application) {
	// Initialize the template cache
	if err := initTemplateCache(); err != nil {
		log.Fatalf("Failed to initialize template cache: %v", err)
	}
	log.Println("Template cache initialized")

	// Initialize the HTTP router with all application routes
	mux := routes()

	// Wrap mux with SCS session middleware
	handler := session.Manager.LoadAndSave(mux)

	// Start the HTTP server on port 4000
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", handler)
	if err != nil {
		log.Printf("Server error: %v", err)
		log.Fatal(err)
	}
}
