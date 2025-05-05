// Package main contains middleware functions for the Gratitude Jar application.
// Middleware functions wrap HTTP handlers to provide additional functionality
// such as logging, security headers, and panic recovery.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/darynforman/gratitude-jar1/internal/session"
)

// RequireLogin ensures the user is logged in, otherwise redirects to login
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := session.GetLoggedInUser(r)
		if userID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin ensures the user is an admin, otherwise returns 403
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, role := session.GetLoggedInUser(r)
		if role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware creates a middleware that logs information about each HTTP request.
// It records:
// - HTTP method (GET, POST, etc.)
// - Request URI
// - Client IP address
// - Time taken to process the request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record the start time

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// Log the request method, URI, client IP, and response time
		log.Printf(
			"%s %s %s %v",
			r.Method,          // HTTP method (e.g., GET, POST)
			r.RequestURI,      // Requested URI
			r.RemoteAddr,      // Client's IP address
			time.Since(start), // Time taken to process the request
		)
	})
}

// SecureHeadersMiddleware adds security-related headers to all HTTP responses.
// These headers help protect against common web vulnerabilities:
// - X-Content-Type-Options: Prevents MIME type sniffing
// - X-Frame-Options: Prevents clickjacking attacks
// - X-XSS-Protection: Enables basic XSS protection
// - Referrer-Policy: Controls referrer information exposure
func SecureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers to enhance protection against common attacks
		w.Header().Set("X-Content-Type-Options", "nosniff")                  // Prevent MIME type sniffing
		w.Header().Set("X-Frame-Options", "deny")                            // Prevent clickjacking attacks
		w.Header().Set("X-XSS-Protection", "1; mode=block")                  // Enable basic XSS protection
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin") // Control referrer information exposure

		// Pass request to the next handler
		next.ServeHTTP(w, r)
	})
}

// RecoverPanicMiddleware recovers from any panics that occur during request handling.
// If a panic occurs:
// 1. The error is logged
// 2. A 500 Internal Server Error response is sent to the client
// 3. The application continues running
func RecoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Recover from panic and log the error
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				// Send a generic 500 error response
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
