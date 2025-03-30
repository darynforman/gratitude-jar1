package main

// PageData holds the data that will be passed to the HTML templates
// It includes the title of the page and a slice of GratitudeNote structs.
type PageData struct {
	Title string          // The title of the page to be displayed in the template
	Notes []GratitudeNote // A slice of GratitudeNote that will be displayed in the template
}

// GratitudeNote represents a single gratitude note in the templates
// Each note has an ID, title, content, category, and creation date.
type GratitudeNote struct {
	ID        int    // Unique identifier for the gratitude note
	Title     string // Title of the gratitude note
	Content   string // Content or the body of the gratitude note
	Category  string // Category of the gratitude note (e.g., personal, work)
	CreatedAt string // The creation date of the gratitude note
}

// GetSampleNotes returns a slice of sample gratitude notes for testing
// This function provides mock data that can be used for testing or development purposes.
func GetSampleNotes() []GratitudeNote {
	// Returning two sample gratitude notes
	return []GratitudeNote{
		{
			ID:        1,
			Title:     "First Gratitude Note",                          // Title of the first note
			Content:   "I'm grateful for the beautiful weather today.", // Content of the first note
			Category:  "personal",                                      // Category of the first note
			CreatedAt: "2024-03-29",                                    // Creation date of the first note
		},
		{
			ID:        2,
			Title:     "Work Achievement",                         // Title of the second note
			Content:   "Completed the project ahead of schedule.", // Content of the second note
			Category:  "work",                                     // Category of the second note
			CreatedAt: "2024-03-28",                               // Creation date of the second note
		},
	}
}
