package data

import (
	"database/sql"
	"time"
)

type GratitudeNote struct {
	ID        int64
	Title     string
	Content   string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GratitudeModel struct {
	DB *sql.DB
}

func (m *GratitudeModel) Insert(title, content, category string) (int64, error) {
	query := `
		INSERT INTO gratitude_notes (title, content, category, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id`

	var id int64
	err := m.DB.QueryRow(query, title, content, category).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *GratitudeModel) Get(id int64) (*GratitudeNote, error) {
	query := `
		SELECT id, title, content, category, created_at, updated_at
		FROM gratitude_notes
		WHERE id = $1`

	note := &GratitudeNote{}
	err := m.DB.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.Category,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (m *GratitudeModel) List() ([]*GratitudeNote, error) {
	query := `
		SELECT id, title, content, category, created_at, updated_at
		FROM gratitude_notes
		ORDER BY created_at DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
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
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}
