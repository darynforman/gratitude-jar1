package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// render renders an HTML template with the provided data
func render(w http.ResponseWriter, tmpl string, data interface{}) {
	log.Printf("Starting to render template: %s", tmpl)

	basePath := filepath.Join("ui", "html", "base.tmpl")
	templatePath := filepath.Join("ui", "html", tmpl)

	// Check if template files exist
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		log.Printf("Base template does not exist: %s", basePath)
		http.Error(w, "Template Not Found", http.StatusInternalServerError)
		return
	}
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Page template does not exist: %s", templatePath)
		http.Error(w, "Template Not Found", http.StatusInternalServerError)
		return
	}

	log.Printf("Loading templates from: base=%s, page=%s", basePath, templatePath)

	// Parse templates
	templates, err := template.ParseFiles(basePath, templatePath)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Templates parsed successfully")

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the base template
	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Template executed successfully")
}
