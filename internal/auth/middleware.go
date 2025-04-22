package auth

import (
	"net/http"
	"github.com/darynforman/gratitude-jar1/internal/session"
)

// RequireLogin is middleware that ensures a user is logged in
func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserIDFromSession(r)
		if userID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// RequireRole is middleware that ensures a user has a specific role
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := GetRoleFromSession(r)
		if userRole != role {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// GetUserIDFromSession retrieves the user ID from the session
func GetUserIDFromSession(r *http.Request) int {
	userID, _ := session.GetLoggedInUser(r)
	return userID
}

// GetRoleFromSession retrieves the user role from the session
func GetRoleFromSession(r *http.Request) string {
	_, role := session.GetLoggedInUser(r)
	return role
}
