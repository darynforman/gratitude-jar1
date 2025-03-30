package main

// PageData holds the data passed to HTML templates
type PageData struct {
	Title string
	Notes []GratitudeNote
}

// GratitudeNote represents a single gratitude note in the templates
type GratitudeNote struct {
	ID        int
	Title     string
	Content   string
	Category  string
	CreatedAt string
}

// GetSampleNotes returns sample gratitude notes for testing
func GetSampleNotes() []GratitudeNote {
	return []GratitudeNote{
		{
			ID:        1,
			Title:     "First Gratitude Note",
			Content:   "I'm grateful for the beautiful weather today.",
			Category:  "personal",
			CreatedAt: "2024-03-29",
		},
		{
			ID:        2,
			Title:     "Work Achievement",
			Content:   "Completed the project ahead of schedule.",
			Category:  "work",
			CreatedAt: "2024-03-28",
		},
	}
}
