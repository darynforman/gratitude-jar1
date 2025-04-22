package data

import "errors"

var (
	// ErrRecordNotFound is returned when a requested record is not found in the database
	ErrRecordNotFound = errors.New("record not found")
)
