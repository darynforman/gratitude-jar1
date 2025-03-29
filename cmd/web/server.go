package main

import (
	"log"
	"net/http"
)

func startServer() {
	// Initialize the server with routes
	mux := routes()

	// Start the server
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
