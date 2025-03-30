package main

import (
	"net/http"
)

type TemplateData struct {
	Title string
	Notes []struct {
		Title     string
		Content   string
		Category  string
		CreatedAt string
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Welcome to Gratitude Jar",
	}
	render(w, r, "home.tmpl", data)
}

func gratitude(w http.ResponseWriter, r *http.Request) {
	// Sample data for testing
	sampleNotes := []struct {
		Title     string
		Content   string
		Category  string
		CreatedAt string
	}{
		{
			Title:     "First Gratitude Note",
			Content:   "I'm grateful for the beautiful weather today.",
			Category:  "personal",
			CreatedAt: "2024-03-29",
		},
		{
			Title:     "Work Achievement",
			Content:   "Completed the project ahead of schedule.",
			Category:  "work",
			CreatedAt: "2024-03-28",
		},
	}

	data := TemplateData{
		Title: "Gratitude Notes",
		Notes: sampleNotes,
	}
	render(w, r, "gratitude.tmpl", data)
}

func createGratitude(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/gratitude", http.StatusSeeOther)
}
