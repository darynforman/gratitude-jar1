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
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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

// Handles viewing all gratitude notes
func viewNotes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling view notes request")

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
		Title: "My Gratitude Notes",
		Notes: templateNotes,
	}
	log.Printf("Created PageData with %d notes", len(data.Notes))

	// Check if the request is from HTMX (for partial updates)
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		// Load and parse only the notes-list section from view-notes.tmpl
		tmpl := template.Must(template.ParseFiles("ui/html/view-notes.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	log.Printf("Rendering view notes template")
	// Otherwise, render the view notes template
	render(w, "view-notes.tmpl", data)
}

// Handles the gratitude page request (for adding new notes)
func gratitude(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling gratitude page request")

	data := PageData{
		Title: "Add Gratitude Note",
	}

	// Check if the request is from HTMX (for partial updates)
	if r.Header.Get("HX-Request") == "true" {
		log.Printf("HTMX request detected, rendering partial template")
		// Load and parse only the notes-list section from gratitude.tmpl
		tmpl := template.Must(template.ParseFiles("ui/html/gratitude.tmpl"))
		tmpl.ExecuteTemplate(w, "notes-list", data)
		return
	}

	log.Printf("Rendering add note template")
	// Otherwise, render the add note template
	render(w, "add-note.tmpl", data)
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

	_, err = getGratitudeModel().Insert(note)
	if err != nil {
		log.Printf("Error inserting note into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If this is an HTMX request, return a redirect response
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/notes")
		w.WriteHeader(http.StatusOK)
		return
	}

	// For regular form submissions, redirect to the notes page
	http.Redirect(w, r, "/notes", http.StatusSeeOther)
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

	// Convert to template note
	templateNote := GratitudeNote{
		ID:        updatedNote.ID,
		Title:     updatedNote.Title,
		Content:   updatedNote.Content,
		Category:  updatedNote.Category,
		Emoji:     updatedNote.Emoji,
		CreatedAt: updatedNote.CreatedAt.Format("2006-01-02"),
	}

	// Render the updated note as HTML
	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("ui/html/view-notes.tmpl"))
	if err := tmpl.ExecuteTemplate(w, "note-card", templateNote); err != nil {
		log.Printf("Error rendering updated note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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

// contactHandler handles the contact page
func contact(w http.ResponseWriter, r *http.Request) {
	log.Printf("Contact handler called with path: %s and method: %s", r.URL.Path, r.Method)

	// Only handle exact /contact path
	if r.URL.Path != "/contact" {
		log.Printf("Invalid contact path: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	// Only handle GET requests
	if r.Method != http.MethodGet {
		log.Printf("Invalid method for contact: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := PageData{
		Title: "Contact Us",
	}

	// Render the contact template with the given data
	render(w, "contact.tmpl", data)
	log.Printf("Contact page rendered successfully")
}
