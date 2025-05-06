package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/darynforman/gratitude-jar1/internal/security"
	"github.com/justinas/nosurf"
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

// CSRFMiddleware adds CSRF protection to all non-GET requests using nosurf
func CSRFMiddleware(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	
	// Determine if we're in production
	isProduction := os.Getenv("ENVIRONMENT") == "production"

	// Configure the CSRF cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})

	// Custom error handler for CSRF failures
	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CSRF error for %s %s", r.Method, r.URL.String())
		
		// For HTMX requests, send a more specific error
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
			"CSRF validation failed",
			false,
		)

		// Redirect to login page for regular requests
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}))

	return csrfHandler
}

// GetCSRFToken is a helper function to get the CSRF token for a request
func GetCSRFToken(r *http.Request) string {
	return nosurf.Token(r)
}
