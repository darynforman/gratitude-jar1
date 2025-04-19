// Package main contains the HTTP route configuration for the Gratitude Jar application.
package main

import "net/http"

// routes sets up all HTTP routes for the application and configures middleware.
// It returns an http.Handler that can be used to start the server.
//
// The function:
// 1. Creates a new ServeMux for routing
// 2. Configures static file serving
// 3. Defines all application routes
// 4. Chains middleware in the correct order
func routes() http.Handler {
	// Create a new ServeMux to handle routing
	mux := http.NewServeMux()

	// Configure static file serving from the ui/static directory
	// All static files will be served under the /static/ URL path
	fileServer := http.FileServer(http.Dir("ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Define application routes
	// Each route is mapped to its corresponding handler function
	mux.HandleFunc("/contact", contact)                  // Contact page
	mux.HandleFunc("/", home)                            // Home page
	mux.HandleFunc("/gratitude", gratitude)              // Gratitude list page
	mux.HandleFunc("/notes", viewNotes)                  // View all notes
	mux.HandleFunc("/gratitude/create", createGratitude) // Create new gratitude entry
	mux.HandleFunc("/gratitude/edit/", getNoteForEdit)   // Get edit form for a note
	mux.HandleFunc("/notes/", updateGratitude)           // Handle PUT and DELETE requests for notes

	// Chain middleware in the correct order
	// The order is important as each middleware wraps the next one
	handler := LoggingMiddleware(mux)          // Log all requests
	handler = SecureHeadersMiddleware(handler) // Add security headers
	handler = RecoverPanicMiddleware(handler)  // Recover from panics

	return handler
}
