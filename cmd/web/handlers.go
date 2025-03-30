package main

import (
	"html/template"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Welcome to Gratitude Jar",
	}
	render(w, r, "home.tmpl", data)
}

func gratitude(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Gratitude Notes",
		Notes: GetSampleNotes(),
	}

	// If it's an HTMX request, only return the notes list
	if r.Header.Get("HX-Request") == "true" {
		tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	render(w, r, "gratitude.tmpl", data)
}

func createGratitude(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// For now, just redirect back to the gratitude page
	// Later we'll add database functionality
	http.Redirect(w, r, "/gratitude", http.StatusSeeOther)
}

func deleteGratitude(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For now, just return success
	// Later we'll add database functionality
	w.WriteHeader(http.StatusOK)
}
