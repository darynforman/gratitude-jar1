// Package main contains template rendering functionality for the Gratitude Jar application.
// This package handles the loading, parsing, and execution of HTML templates.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// templateDir is the directory containing HTML templates
	templateDir = "ui/html"
)

// SetTemplateDir sets the template directory for testing
func SetTemplateDir(dir string) {
	templateDir = dir
}

// render renders an HTML template with the provided data.
// It performs the following tasks:
// 1. Checks if the required template files exist
// 2. Loads and parses the base template and page template
// 3. Sets the appropriate content type header
// 4. Executes the template with the provided data
//
// Parameters:
// - w: The HTTP response writer
// - tmpl: The name of the template file to render (without path)
// - data: The data to pass to the template
func render(w http.ResponseWriter, tmpl string, data interface{}) {
	log.Printf("Starting to render template: %s", tmpl)

	// Get the absolute path to the templates directory
	absTemplateDir, err := filepath.Abs(templateDir)
	if err != nil {
		log.Printf("Error getting template directory: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Construct paths to template files
	basePath := filepath.Join(absTemplateDir, "base.tmpl")
	templatePath := filepath.Join(absTemplateDir, tmpl)

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

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute template
	err = templates.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
