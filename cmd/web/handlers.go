package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

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
			ID:        int(note.ID),
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

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.FormValue("category")
	emoji := r.FormValue("emoji")

	log.Printf("Received form data - Title: %s, Category: %s, Emoji: %s", title, category, emoji)

	if title == "" || content == "" || emoji == "" {
		log.Printf("Missing required fields - Title: %s, Content: %s, Emoji: %s", title, content, emoji)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Insert note into database
	id, err := getGratitudeModel().Insert(title, content, category, emoji)
	if err != nil {
		log.Printf("Error inserting note into database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created note with ID: %d", id)
	http.Redirect(w, r, "/gratitude", http.StatusSeeOther)
}

// Handles the deletion of gratitude notes
func deleteGratitude(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := r.URL.Path[len("/gratitude/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Delete note from database
	err = getGratitudeModel().Delete(id)
	if err != nil {
		log.Printf("Error deleting note: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
