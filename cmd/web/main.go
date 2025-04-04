package main

import (
	"log"
	"net/http"

	"github.com/darynforman/gratitude-jar/internal/config"
)

func main() {
	// Initialize database connection
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize the server with routes
	mux := routes()

	// Start the server
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
