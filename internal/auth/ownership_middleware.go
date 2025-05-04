package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/darynforman/gratitude-jar1/internal/config"
	"github.com/darynforman/gratitude-jar1/internal/data"
	"github.com/darynforman/gratitude-jar1/internal/session"
)

// RequireOwnership ensures the user owns the requested resource
func RequireOwnership(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from session
		userID := session.Manager.GetInt(r, "userID")
		if userID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Extract resource ID from URL
		// Assuming URLs like /gratitude/edit/123 or /notes/123
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			next(w, r)
			return
		}

		resourceIDStr := parts[len(parts)-1]
		resourceID, err := strconv.Atoi(resourceIDStr)
		if err != nil {
			// If we can't parse the ID, just continue
			next(w, r)
			return
		}

		// Check if this is a gratitude note
		if strings.Contains(r.URL.Path, "/gratitude/") || strings.Contains(r.URL.Path, "/notes/") {
			// Get the note
			gratitudeModel := data.NewGratitudeModel(config.DB)
			note, err := gratitudeModel.Get(resourceID)
			if err != nil || note == nil {
				http.NotFound(w, r)
				return
			}

			// Check ownership
			if note.UserID != userID {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}

		// If we get here, the user owns the resource or it's not a resource that needs checking
		next(w, r)
	}
}
