// Package main contains data structures and functions for template data in the Gratitude Jar application.
// This package defines the data models used in HTML templates and provides functions for managing gratitude notes.
package main

import "github.com/darynforman/gratitude-jar1/internal/data"

// PageData holds data passed to templates
type PageData struct {
	Title           string               // The title of the page to be displayed in the template
	Notes           []data.GratitudeNote // A slice of GratitudeNote that will be displayed in the template
	Note            *data.GratitudeNote  // A single gratitude note for editing/viewing
	Errors          map[string]string    // Validation errors for form fields
	Emojis          []string             // Available emojis for gratitude note creation
	Form            map[string]string    // Form values for re-populating registration/login
	IsAuthenticated bool                 // Indicates whether the user is authenticated
	UserRole        string               // The role of the authenticated user
	Flash           string               // Flash messages for user feedback
	SuccessMessage  string               // Success message for form submissions
}

// GratitudeNote represents a single gratitude note in the templates.
// This structure is used both for displaying notes and for processing form submissions.
type GratitudeNote struct {
	ID        int    // Unique identifier for the gratitude note
	Title     string // Title of the gratitude note
	Content   string // Content or the body of the gratitude note
	Category  string // Category of the gratitude note (e.g., personal, work)
	Emoji     string // Emoji representing the mood/feeling
	CreatedAt string // The creation date of the gratitude note
}

// notes is our in-memory storage for gratitude notes.
// This slice holds all notes, including both sample data and user-created notes.
var notes []GratitudeNote

// init initializes the notes slice with sample data.
// This function is called automatically when the package is loaded.
func init() {
	notes = GetSampleNotes()
}

// GetNotes returns all notes from the in-memory storage.
// This includes both sample notes and any user-created notes that have been added.
func GetNotes() []GratitudeNote {
	return notes
}

// AddNote adds a new note to the in-memory storage.
// The note is appended to the existing slice of notes.
func AddNote(note GratitudeNote) {
	notes = append(notes, note)
}

// GetSampleNotes returns a slice of sample gratitude notes for testing and development.
// These notes are used to populate the application with initial data and for demonstration purposes.
func GetSampleNotes() []GratitudeNote {
	// Returning two sample gratitude notes with different categories and content
	return []GratitudeNote{
		{
			ID:        1,
			Title:     "First Gratitude Note",                          // Title of the first note
			Content:   "I'm grateful for the beautiful weather today.", // Content of the first note
			Category:  "personal",                                      // Category of the first note
			Emoji:     "😊",                                             // Emoji for the first note
			CreatedAt: "2024-03-29",                                    // Creation date of the first note
		},
		{
			ID:        2,
			Title:     "Work Achievement",                         // Title of the second note
			Content:   "Completed the project ahead of schedule.", // Content of the second note
			Category:  "work",                                     // Category of the second note
			Emoji:     "🤩",                                        // Emoji for the second note
			CreatedAt: "2024-03-28",                               // Creation date of the second note
		},
	}
}
