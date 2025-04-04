package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/darynforman/gratitude-jar/internal/data"
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

// getGratitudeModel returns a new instance of GratitudeModel
func getGratitudeModel() *data.GratitudeModel {
	return data.NewGratitudeModel()
}

// Handles the gratitude page request
func gratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling gratitude page request")

	// Get notes from database
	notes, err := getGratitudeModel().List()
	if err != nil {
		log.Printf("Error fetching notes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert database notes to template notes
	var templateNotes []GratitudeNote
	for _, note := range notes {
		templateNotes = append(templateNotes, GratitudeNote{
			ID:        note.ID,
			Title:     note.Title,
			Content:   note.Content,
			Category:  note.Category,
			Emoji:     note.Emoji,
			CreatedAt: note.CreatedAt.Format("2006-01-02"),
		})
	}

	data := PageData{
		Title: "Gratitude Notes",
		Notes: templateNotes,
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
	log.Printf("Handling create gratitude note request")

	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	note := &data.GratitudeNote{
		Title:     r.PostForm.Get("title"),
		Content:   r.PostForm.Get("content"),
		Category:  r.PostForm.Get("category"),
		Emoji:     r.PostForm.Get("emoji"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := getGratitudeModel().Insert(note)
	if err != nil {
		log.Printf("Error inserting note into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the newly created note
	note, err = getGratitudeModel().Get(id)
	if err != nil {
		log.Printf("Error fetching new note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Render the new note as HTML
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
	if err := tmpl.ExecuteTemplate(w, "note-card", note); err != nil {
		log.Printf("Error rendering new note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handles both updating and deleting gratitude notes
func updateGratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling update/delete request with method: %s", r.Method)
	log.Printf("Request URL: %s", r.URL.Path)

	// Extract ID from URL
	idStr := r.URL.Path[len("/gratitude/update/"):]
	log.Printf("Extracted ID string: %s", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed ID: %d", id)

	// Handle DELETE request
	if r.Method == http.MethodDelete {
		log.Printf("Processing DELETE request for note ID: %d", id)
		err = getGratitudeModel().Delete(id)
		if err != nil {
			log.Printf("Error deleting note: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// Handle PUT request
	if r.Method != http.MethodPut {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Processing PUT request for note ID: %d", id)

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Log all form values
	log.Printf("Form values received:")
	for key, values := range r.PostForm {
		log.Printf("  %s: %v", key, values)
	}

	// Create updated note
	note := &data.GratitudeNote{
		ID:        id,
		Title:     r.PostForm.Get("title"),
		Content:   r.PostForm.Get("content"),
		Category:  r.PostForm.Get("category"),
		Emoji:     r.PostForm.Get("emoji"),
		UpdatedAt: time.Now(),
	}
	log.Printf("Created note object: %+v", note)

	// Update note in database
	err = getGratitudeModel().Update(note)
	if err != nil {
		log.Printf("Error updating note in database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully updated note in database")

	// Fetch the updated note
	updatedNote, err := getGratitudeModel().Get(id)
	if err != nil {
		log.Printf("Error fetching updated note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched updated note: %+v", updatedNote)

	// Render the updated note as HTML
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
	if err := tmpl.ExecuteTemplate(w, "note-card", updatedNote); err != nil {
		log.Printf("Error rendering updated note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Printf("Successfully rendered updated note")
}

// Handles getting a note for editing
func getNoteForEdit(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling get note for edit request")

	// Extract ID from URL
	idStr := r.URL.Path[len("/gratitude/edit/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get note from database
	note, err := getGratitudeModel().Get(id)
	if err != nil {
		log.Printf("Error fetching note: %v", err)
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	// Render edit form
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("ui/html/edit-form.tmpl"))
	if err := tmpl.ExecuteTemplate(w, "edit-form", note); err != nil {
		log.Printf("Error rendering edit form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
