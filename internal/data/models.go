package data

import "database/sql"

// Models holds the models for the application
type Models struct {
	Users      *UserModel
	Gratitudes *GratitudeModel
}

// NewModels creates a new Models instance
func NewModels(db *sql.DB) *Models {
	return &Models{
		Users:      NewUserModel(db),
		Gratitudes: NewGratitudeModel(db),
	}
}
