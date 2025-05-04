// Package main contains template rendering functionality for the Gratitude Jar application.
// This package handles the loading, parsing, and execution of HTML templates.
package main

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/darynforman/gratitude-jar1/internal/auth"
	"github.com/darynforman/gratitude-jar1/internal/session"
)

// render renders a template with the given data
func render(w http.ResponseWriter, r *http.Request, name string, data PageData) {
	// Get the template from the cache
	tmpl, err := getTemplate(name)
	if err != nil {
		log.Printf("Template %s not found in cache: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Add session data to the template data
	userID := session.Manager.GetInt(r, "userID")
	role := session.Manager.GetString(r, "role")
	flash := session.Manager.PopString(r, "flash")

	// Create a new data struct that includes session data
	templateData := struct {
		PageData
		IsAuthenticated bool
		UserID          int
		UserRole        string
		Flash           string
		CurrentYear     int
		CSRFToken       string
	}{
		PageData:        data,
		IsAuthenticated: userID > 0,
		UserID:          userID,
		UserRole:        role,
		Flash:           flash,
		CurrentYear:     time.Now().Year(),
		CSRFToken:       auth.GetCSRFToken(r),
	}

	// For partial templates, execute without base template
	if filepath.Dir(name) == "partials" {
		err := tmpl.Execute(w, templateData)
		if err != nil {
			log.Printf("Error executing partial template %s: %v", name, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}

	// For full pages, execute with base template
	err = tmpl.ExecuteTemplate(w, "base", templateData)
	if err != nil {
		log.Printf("Error executing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
