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
	userID, ok := Manager.Get(r, "userID").(int)
	if !ok {
		userID = 0
	}
	role, ok := Manager.Get(r, "role").(string)
	if !ok {
		role = ""
	}
	return userID, role
}
