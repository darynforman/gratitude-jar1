package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var Manager *scs.SessionManager

func init() {
	Manager = scs.New()
	Manager.Lifetime = 24 * time.Hour
	Manager.Cookie.HttpOnly = true
	Manager.Cookie.Secure = false // Set to true in production with HTTPS
}

// GetLoggedInUser returns the user ID and role from the session, or 0, "" if not logged in
func GetLoggedInUser(r *http.Request) (int, string) {
	userID := Manager.GetInt(r.Context(), "userID")
	role := Manager.GetString(r.Context(), "role")
	return userID, role
}
