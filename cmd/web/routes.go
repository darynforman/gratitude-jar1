// Package main contains the HTTP route configuration for the Gratitude Jar application.
package main

import (
	"net/http"

	"github.com/darynforman/gratitude-jar1/internal/auth"
	"github.com/justinas/nosurf"
)

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
	mux.HandleFunc("/contact", contact) // Contact page
	mux.HandleFunc("/", home)           // Home page
	// Protected routes
	mux.Handle("/gratitude", auth.RequireLogin(http.HandlerFunc(gratitude)))
	mux.Handle("/notes", auth.RequireLogin(http.HandlerFunc(viewNotes)))
	mux.Handle("/gratitude/create", auth.RequireLogin(http.HandlerFunc(createGratitude)))
	mux.Handle("/gratitude/edit/", auth.RequireLogin(auth.RequireOwnership(http.HandlerFunc(getNoteForEdit))))
	mux.Handle("/notes/", auth.RequireLogin(auth.RequireOwnership(http.HandlerFunc(updateGratitude))))

	// Auth routes
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/user/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	// Chain middleware in the correct order
	// The order is important as each middleware wraps the next one
	handler := LoggingMiddleware(mux)                // Log all requests
	handler = RateLimitMiddleware(handler)           // Rate limiting
	handler = SecureHeadersMiddleware(handler)       // Add security headers
	handler = auth.SessionTimeoutMiddleware(handler) // Check session timeout
	handler = RecoverPanicMiddleware(handler)        // Recover from panics
	handler = nosurf.New(handler)                   // Add CSRF protection

	return handler
}
