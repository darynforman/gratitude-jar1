package auth

import (
	"net/http"

	"github.com/darynforman/gratitude-jar1/internal/session"
)

// RequireLogin is middleware that ensures a user is logged in
func RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if session cookie exists
		_, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Check if user ID exists in session
		userID := GetUserIDFromSession(r)
		if userID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
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
