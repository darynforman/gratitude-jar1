package auth

import (
	"net/http"
	"time"

	"github.com/darynforman/gratitude-jar1/internal/session"
)

// SessionTimeoutMiddleware checks if the session has timed out due to inactivity
func SessionTimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip for unauthenticated users
		userID := session.Manager.GetInt(r, "userID")
		if userID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Get last activity time from session
		lastActivityVal := session.Manager.Get(r, "last_activity")
		
		// Check if last_activity exists and is a valid time
		var lastActivity time.Time
		if lastActivityVal != nil {
			if lastTime, ok := lastActivityVal.(time.Time); ok {
				lastActivity = lastTime
			}
		}

		now := time.Now()
		
		// If session has been inactive for too long (30 minutes), log out
		if !lastActivity.IsZero() && now.Sub(lastActivity) > 30*time.Minute {
			session.LogoutUser(w, r)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		
		// Update last activity time
		session.Manager.Put(r, "last_activity", now)
		
		next.ServeHTTP(w, r)
	})
}
