package data

import (
	"context"
	"database/sql"
	"errors"
)

// User represents a user in the system.
type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	Role         string
}

// UserModel wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// GetByEmail fetches a user by email
func (m *UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, role FROM users WHERE email = $1`
	user := &User{}
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Insert adds a new user to the database
func (m *UserModel) Insert(ctx context.Context, username, email, passwordHash, role string) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4)`
	_, err := m.DB.ExecContext(ctx, query, username, email, passwordHash, role)
	return err
}

// GetByUsername fetches a user by username.
func (m *UserModel) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password_hash, role FROM users WHERE username = $1`
	user := &User{}
	err := m.DB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// NewUserModel creates a new UserModel instance.
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{DB: db}
}
