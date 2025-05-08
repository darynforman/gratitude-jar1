package data

import (
	"context"
	"database/sql"
	"time"
)

// GratitudeNote represents a gratitude note in the database
type GratitudeNote struct {
	ID        int
	Title     string
	Content   string
	Category  string
	Emoji     string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GratitudeModel wraps a database connection pool
type GratitudeModel struct {
	DB *sql.DB
}

// NewGratitudeModel creates a new GratitudeModel instance
func NewGratitudeModel(db *sql.DB) *GratitudeModel {
	return &GratitudeModel{DB: db}
}

// GetAll returns all gratitude notes for the current user
func (m *GratitudeModel) GetAll(ctx context.Context, userID int) ([]GratitudeNote, error) {
	query := `SELECT id, title, content, category, emoji, created_at, updated_at 
	          FROM gratitude_notes 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC`
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []GratitudeNote
	for rows.Next() {
		var note GratitudeNote
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Category, &note.Emoji, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

// Get returns a single gratitude note by ID
func (m *GratitudeModel) Get(ctx context.Context, id int) (*GratitudeNote, error) {
	query := `SELECT id, title, content, category, emoji, user_id, created_at, updated_at 
	          FROM gratitude_notes 
	          WHERE id = $1`
	note := &GratitudeNote{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&note.ID, &note.Title, &note.Content, &note.Category, &note.Emoji, &note.UserID, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return note, nil
}

// Insert creates a new gratitude note
func (m *GratitudeModel) Insert(ctx context.Context, note *GratitudeNote) error {
	query := `INSERT INTO gratitude_notes (title, content, category, emoji, user_id, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) 
	          RETURNING id`
	err := m.DB.QueryRowContext(
		ctx,
		query,
		note.Title,
		note.Content,
		note.Category,
		note.Emoji,
		note.UserID,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(&note.ID)
	return err
}

// Update modifies an existing gratitude note
func (m *GratitudeModel) Update(ctx context.Context, note *GratitudeNote) error {
	query := `UPDATE gratitude_notes 
	          SET title = $1, content = $2, category = $3, emoji = $4, updated_at = $5 
	          WHERE id = $6 AND user_id = $7`
	_, err := m.DB.ExecContext(
		ctx,
		query,
		note.Title,
		note.Content,
		note.Category,
		note.Emoji,
		note.UpdatedAt,
		note.ID,
		note.UserID,
	)
	return err
}

// Delete removes a gratitude note
func (m *GratitudeModel) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM gratitude_notes WHERE id = $1 AND user_id = $2`
	_, err := m.DB.ExecContext(ctx, query, id, userID)
	return err
}
