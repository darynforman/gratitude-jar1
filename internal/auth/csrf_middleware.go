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

	// Configure CSRF protection
	return csrf.Protect(
		CSRFKey(),
		csrf.Secure(isProduction), // Set to true in production
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reason := csrf.FailureReason(r)
			log.Printf("CSRF error: %v", reason)

			// Log CSRF failure as a security event
			security.LogSecurityEvent(
				security.EventCSRFFailure,
				0, // We don't know the user ID in this context
				"",
				security.GetClientIP(r),
				"CSRF validation failed: "+reason.Error(),
				false,
			)

			http.Error(w, "CSRF token validation failed", http.StatusForbidden)
		})),
	)(next)
}

// GetCSRFToken is a helper function to get the CSRF token for a request
func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}
