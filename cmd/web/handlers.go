package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// Handles the home page request
func home(w http.ResponseWriter, r *http.Request) {
	// Create a PageData struct with a title for the homepage
	data := PageData{
		Title: "Welcome to Gratitude Jar",
	}
	// Render the home template with the given data
	render(w, "home.tmpl", data)
}

// Handles the gratitude page request
func gratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling gratitude page request")
	// Create a PageData struct with a title and all gratitude notes
	data := PageData{
		Title: "Gratitude Notes",
		Notes: GetNotes(),
	}
	log.Printf("Created PageData with %d notes", len(data.Notes))

	// Check if the request is from HTMX (for partial updates)
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		// Load and parse only the notes-list section from gratitude.tmpl
		tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	log.Printf("Rendering full gratitude template")
	// Otherwise, render the full gratitude template
	render(w, "gratitude.tmpl", data)
}

// Handles form submissions for adding gratitude notes
func createGratitude(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data from the request
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")
	emoji := r.FormValue("emoji")

	// Validate required fields
	if title == "" || content == "" || emoji == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create a new note
	note := GratitudeNote{
		ID:        len(GetNotes()) + 1, // Simple ID generation for now
		Title:     title,
		Content:   content,
		Category:  category,
		Emoji:     emoji,
		CreatedAt: time.Now().Format("2006-01-02"),
	}

	// Add the note to our in-memory storage
	AddNote(note)

	// Redirect back to the gratitude page
	http.Redirect(w, r, "/gratitude", http.StatusSeeOther)
}

// Handles the deletion of gratitude notes
func deleteGratitude(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Placeholder: No actual deletion logic implemented yet
	// Respond with HTTP 200 OK to indicate success
	w.WriteHeader(http.StatusOK)
}
