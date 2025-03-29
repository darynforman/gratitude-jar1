package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func render(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	// Parse both templates together
	templates := template.Must(template.ParseFiles(
		filepath.Join("html", "base.tmpl"),
		filepath.Join("html", tmpl),
	))

	// Execute the templates
	err := templates.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
