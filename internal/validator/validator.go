// Package validator provides validation utilities for the Gratitude Jar application.
// It includes validation rules for gratitude notes and other user inputs.
package validator

import (
	"strings"
	"unicode/utf8"
)

// Validator holds validation errors for form fields.
// It uses a map to store field-specific error messages.
type Validator struct {
	Errors map[string]string
}

// NewValidator creates and returns a new Validator instance
// with an initialized errors map.
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// ValidData checks if the validator has any errors.
// Returns true if there are no validation errors.
func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a specific field.
// If the field already has an error, it won't be overwritten.
func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

// Check performs a validation check and adds an error if the check fails.
// It takes a boolean condition, field name, and error message as parameters.
func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

// NotBlank checks if a string value is not empty or just whitespace.
// Returns true if the string contains non-whitespace characters.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MinLength checks if a string meets the minimum length requirement.
// The length is measured in Unicode code points, not bytes.
func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// MaxLength checks if a string does not exceed the maximum length.
// The length is measured in Unicode code points, not bytes.
func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// ValidCategory checks if a category is valid.
// Currently valid categories are: personal, work, health, relationships, other
func ValidCategory(category string) bool {
	validCategories := map[string]bool{
		"personal":     true,
		"work":         true,
		"family":       true,
		"achievements": true,
		"health":       true,
		"experiences":  true,
	}
	return validCategories[category]
}

// ValidEmoji checks if the emoji string is valid.
// This is a simple check that verifies the string is not too long
// and contains at least one emoji character.
func (v *Validator) ValidEmoji(emoji string) bool {
	// Emojis can be 1-8 characters long (including variation selectors)
	if len(emoji) < 1 || len(emoji) > 8 {
		return false
	}
	// Check if the string contains at least one emoji character
	for _, r := range emoji {
		if r >= 0x1F300 && r <= 0x1F9FF { // Basic emoji range
			return true
		}
		if r >= 0x2600 && r <= 0x26FF { // Misc symbols
			return true
		}
		if r >= 0x2700 && r <= 0x27BF { // Dingbats
			return true
		}
	}
	return false
}

// ValidateGratitudeNote validates a gratitude note's fields.
// It checks:
// - Title is not blank and within length limits
// - Content is not blank and within length limits
// - Category is valid
// - Emoji is valid
func ValidateGratitudeNote(title, content, category, emoji string) *Validator {
	v := NewValidator()

	// Validate title
	v.Check(NotBlank(title), "title", "Title cannot be blank")
	v.Check(MinLength(title, 3), "title", "Title must be at least 3 characters long")
	v.Check(MaxLength(title, 100), "title", "Title cannot be more than 100 characters long")

	// Validate content
	v.Check(NotBlank(content), "content", "Content cannot be blank")
	v.Check(MinLength(content, 10), "content", "Content must be at least 10 characters long")
	v.Check(MaxLength(content, 1000), "content", "Content cannot be more than 1000 characters long")

	// Validate category
	v.Check(ValidCategory(category), "category", "Please select a valid category")

	// Validate emoji
	v.Check(v.ValidEmoji(emoji), "emoji", "Please select a valid emoji")

	return v
}
