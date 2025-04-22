package main

import (
	"net/http"

	"github.com/darynforman/gratitude-jar1/internal/session"
)

// GetLoggedInUser returns the user ID and role from the session, or 0, "" if not logged in
func GetLoggedInUser(r *http.Request) (int, string) {
	return session.GetLoggedInUser(r)
}
