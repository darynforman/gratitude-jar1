package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the server
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Start the server
	log.Println("Starting server on :4000...")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
