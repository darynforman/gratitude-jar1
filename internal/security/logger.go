package security

import (
	"log"
	"net/http"
	"time"
)

// EventType represents the type of security event
type EventType string

const (
	// EventLogin represents a login attempt
	EventLogin EventType = "LOGIN"
	// EventLogout represents a logout event
	EventLogout EventType = "LOGOUT"
	// EventRegistration represents a user registration
	EventRegistration EventType = "REGISTRATION"
	// EventPasswordChange represents a password change
	EventPasswordChange EventType = "PASSWORD_CHANGE"
	// EventAccessDenied represents an access denied event
	EventAccessDenied EventType = "ACCESS_DENIED"
	// EventCSRFFailure represents a CSRF token validation failure
	EventCSRFFailure EventType = "CSRF_FAILURE"
	// EventRateLimitExceeded represents a rate limit exceeded event
	EventRateLimitExceeded EventType = "RATE_LIMIT_EXCEEDED"
)

// LogSecurityEvent logs a security event with the given details
func LogSecurityEvent(eventType EventType, userID int, username, ipAddress, details string, success bool) {
	outcome := "SUCCESS"
	if !success {
		outcome = "FAILURE"
	}
	
	log.Printf("[SECURITY] %s | %s | User ID: %d | Username: %s | IP: %s | Details: %s | Time: %s",
		eventType,
		outcome,
		userID,
		username,
		ipAddress,
		details,
		time.Now().Format(time.RFC3339),
	)
}

// GetClientIP gets the client IP address from the request
func GetClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	
	// Otherwise, use RemoteAddr
	return r.RemoteAddr
}
