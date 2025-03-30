package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// render renders an HTML template with the provided data
func render(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse both the base template and the specific page template together
	templates := template.Must(template.ParseFiles(
		filepath.Join("ui", "html", "base.tmpl"), // Base layout template (common structure)
		filepath.Join("ui", "html", tmpl),        // Specific page template
	))

	// Execute the parsed templates and send the output to the response writer
	err := templates.Execute(w, data)
	if err != nil {
		// Log the error if template execution fails
		log.Printf("Error executing template: %v", err)

		// Send a 500 Internal Server Error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
