package main

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs the incoming HTTP request and its duration
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

// SecureHeadersMiddleware adds security-related headers to all responses
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

// RecoverPanicMiddleware recovers from any panics and returns a 500 Internal Server Error
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
