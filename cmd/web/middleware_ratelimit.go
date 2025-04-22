package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/darynforman/gratitude-jar1/internal/ratelimit"
)

var (
	// Global rate limiter instance
	globalLimiter = ratelimit.NewRateLimiter(10, 20) // 10 requests per second, burst of 20

	// Cleanup old rate limiters every hour
	cleanupInterval = 1 * time.Hour
)

func init() {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			globalLimiter.Cleanup(cleanupInterval)
		}
	}()
}

// RateLimitMiddleware limits the number of requests from each IP address
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip rate limiting for static files
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		// Get client IP
		ip := getClientIP(r)
		
		// Get rate limiter for this IP
		limiter := globalLimiter.GetLimiter(ip)
		
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getClientIP gets the client's real IP address, taking into account proxy headers
func getClientIP(r *http.Request) string {
	// Check X-Real-IP header
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Check X-Forwarded-For header
	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}
