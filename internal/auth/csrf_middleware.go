package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/darynforman/gratitude-jar1/internal/security"
	"github.com/gorilla/csrf"
)

// CSRFKey returns the CSRF key from environment variable or a default for development
func CSRFKey() []byte {
	key := os.Getenv("CSRF_KEY")
	if key == "" {
		// Only use this default in development
		log.Println("WARNING: Using default CSRF key. Set CSRF_KEY environment variable in production.")
		key = "32-byte-long-auth-key-for-development"
	}
	return []byte(key)
}

// CSRFMiddleware adds CSRF protection to all non-GET requests
func CSRFMiddleware(next http.Handler) http.Handler {
	// Determine if we're in production
	isProduction := os.Getenv("ENVIRONMENT") == "production"

	log.Printf("Initializing CSRF middleware (secure mode: %v)", isProduction)

	// Basic options that work for both environments
	opts := []csrf.Option{
		csrf.Path("/"),
		csrf.CookieName("_gorilla_csrf"),
		csrf.FieldName("gorilla.csrf.Token"),
		csrf.Secure(isProduction),
		csrf.RequestHeader("X-CSRF-Token"),
	}

	if !isProduction {
		// In development mode, accept requests from localhost
		opts = append(opts, csrf.TrustedOrigins([]string{"http://localhost:4002"}))
	}

	// Add error handler to options
	opts = append(opts, csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reason := csrf.FailureReason(r)
		log.Printf("CSRF error: %v", reason)
		log.Printf("Request Method: %s", r.Method)
		log.Printf("Request URL: %s", r.URL.String())
		log.Printf("Request Headers: %v", r.Header)
		log.Printf("Cookie Header: %s", r.Header.Get("Cookie"))
		log.Printf("Origin Header: %s", r.Header.Get("Origin"))
		log.Printf("Referer Header: %s", r.Header.Get("Referer"))

		// For HTMX requests, send a more specific error with HTMX-specific headers
		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Retarget", "#error-container")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("<div class='error'>CSRF token validation failed. Please refresh the page and try again.</div>"))
			return
		}

		// Log CSRF failure as a security event
		security.LogSecurityEvent(
			security.EventCSRFFailure,
			0,
			"",
			security.GetClientIP(r),
			"CSRF validation failed: "+reason.Error(),
			false,
		)

		// Redirect to login page for regular requests
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	})))

	return csrf.Protect(CSRFKey(), opts...)(next)
}

// GetCSRFToken is a helper function to get the CSRF token for a request
func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}
