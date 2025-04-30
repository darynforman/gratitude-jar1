package session

import (
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
)

var Manager *sessions.Session

func init() {
	Manager = sessions.New([]byte("your-secret-key"))
	Manager.Lifetime = 24 * time.Hour
	Manager.Secure = false // Set to true in production with HTTPS
}

// GetLoggedInUser returns the user ID and role from the session, or 0, "" if not logged in
func GetLoggedInUser(r *http.Request) (int, string) {
	// First check if the values exist at all
	userIDVal := Manager.Get(r, "userID")
	if userIDVal == nil {
		return 0, ""
	}

	// Then try to convert to int
	userID, ok := userIDVal.(int)
	if !ok {
		return 0, ""
	}

	// Get role
	roleVal := Manager.Get(r, "role")
	if roleVal == nil {
		return 0, ""
	}

	// Convert to string
	role, ok := roleVal.(string)
	if !ok {
		return 0, ""
	}

	return userID, role
}

// LogoutUser properly cleans up the session and ensures the user is logged out
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	// Clear all session data
	Manager.Put(r, "userID", nil)
	Manager.Put(r, "role", nil)
	Manager.Put(r, "flash", nil)

	// Force the session to be saved with the cleared values
	Manager.Put(r, "_cleared", time.Now().Unix())

	// Destroy the session
	Manager.Destroy(r)

	// Explicitly remove the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
