package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/gratitude", gratitude)
	mux.HandleFunc("/gratitude/create", createGratitude)
	mux.HandleFunc("/gratitude/", deleteGratitude)

	// Chain middleware
	handler := LoggingMiddleware(mux)
	handler = SecureHeadersMiddleware(handler)
	handler = RecoverPanicMiddleware(handler)

	return handler
}
