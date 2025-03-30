package main

import (
	"log"
	"net/http"
)

// startServer initializes and starts the HTTP server
func startServer() {
	// Initialize the server with routes
	mux := routes() // Assume 'routes' is a function that defines the routing logic

	// Start the server
	log.Println("Starting server on :4000...") // Using port 4000 consistently
	err := http.ListenAndServe(":4000", mux)   // Using port 4000 consistently
	if err != nil {                            // Check for any errors during server start
		log.Printf("Server error: %v", err) // Detailed error logging
		log.Fatal(err)                      // Log fatal error and stop execution if the server fails to start
	}
}
