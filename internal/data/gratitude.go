package data

import (
	"database/sql"
	"log"
	"time"

	"github.com/darynforman/gratitude-jar/internal/config"
)

type GratitudeNote struct {
	ID        int
	Title     string
	Content   string
	Category  string
	Emoji     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GratitudeModel struct {
	DB *sql.DB
}

// NewGratitudeModel creates a new GratitudeModel instance
func NewGratitudeModel() *GratitudeModel {
	return &GratitudeModel{
		DB: config.DB,
	}
}

// Insert adds a new gratitude note to the database
func (m *GratitudeModel) Insert(note *GratitudeNote) (int, error) {
	log.Printf("Attempting to insert note with title: %s, category: %s, emoji: %s", note.Title, note.Category, note.Emoji)

	query := `
		INSERT INTO gratitude_notes (title, content, category, emoji, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var id int
	err := m.DB.QueryRow(query, note.Title, note.Content, note.Category, note.Emoji, note.CreatedAt, note.UpdatedAt).Scan(&id)
	if err != nil {
		log.Printf("Error inserting note: %v", err)
		return 0, err
	}
	log.Printf("Successfully inserted note with ID: %d", id)
	return id, nil
}

// Get retrieves a gratitude note by ID
func (m *GratitudeModel) Get(id int) (*GratitudeNote, error) {
	query := `
		SELECT id, title, content, category, emoji, created_at, updated_at
		FROM gratitude_notes
		WHERE id = $1`

	note := &GratitudeNote{}
	err := m.DB.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.Category,
		&note.Emoji,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

// List retrieves all gratitude notes, ordered by creation date
func (m *GratitudeModel) List() ([]*GratitudeNote, error) {
	log.Printf("Attempting to list all notes")

	query := `
		SELECT id, title, content, category, emoji, created_at, updated_at
		FROM gratitude_notes
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		log.Printf("Error querying notes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notes []*GratitudeNote
	for rows.Next() {
		note := &GratitudeNote{}
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.Category,
			&note.Emoji,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning note: %v", err)
			return nil, err
		}
		notes = append(notes, note)
	}
	log.Printf("Successfully retrieved %d notes", len(notes))
	return notes, nil
}

// Delete removes a gratitude note by ID
func (m *GratitudeModel) Delete(id int) error {
	query := `DELETE FROM gratitude_notes WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

// Update modifies an existing gratitude note
func (m *GratitudeModel) Update(note *GratitudeNote) error {
	log.Printf("Attempting to update note with ID: %d", note.ID)

	query := `
		UPDATE gratitude_notes 
		SET title = $1, content = $2, category = $3, emoji = $4, updated_at = $5
		WHERE id = $6`

	result, err := m.DB.Exec(query, note.Title, note.Content, note.Category, note.Emoji, note.UpdatedAt, note.ID)
	if err != nil {
		log.Printf("Error updating note: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	log.Printf("Successfully updated note with ID: %d", note.ID)
	return nil
}
