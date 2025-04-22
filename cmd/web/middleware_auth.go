package main

import (
	"net/http"

	"github.com/darynforman/gratitude-jar1/internal/auth"
)

// RequirePermission creates middleware that checks if the user has the required permission
func RequirePermission(permission auth.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from session
			_, role := GetLoggedInUser(r)
			
			// Check if user has the required permission
			if !auth.HasPermission(role, permission) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyPermission creates middleware that checks if the user has any of the required permissions
func RequireAnyPermission(permissions ...auth.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from session
			_, role := GetLoggedInUser(r)
			
			// Check if user has any of the required permissions
			for _, permission := range permissions {
				if auth.HasPermission(role, permission) {
					next.ServeHTTP(w, r)
					return
				}
			}
			
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
