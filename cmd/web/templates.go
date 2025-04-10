package main

import (
	"errors"
	"html/template"
	"path/filepath"
	"time"
)

// ErrTemplateNotFound is returned when a template is not found in the cache
var ErrTemplateNotFound = errors.New("template not found")

// humanDate formats a time.Time value to a human-readable string
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// formatDateShort formats a time.Time value to a short date string
func formatDateShort(t time.Time) string {
	return t.Format("Jan 02, 2006")
}

// functions is a map of template functions that can be used in templates
var functions = template.FuncMap{
	"humanDate":       humanDate,
	"formatDateShort": formatDateShort,
}

// templateCache holds the parsed templates
var templateCache map[string]*template.Template

// initTemplateCache initializes the template cache
func initTemplateCache() error {
	cache := map[string]*template.Template{}

	// Get all page templates
	pages, err := filepath.Glob("ui/html/*.tmpl")
	if err != nil {
		return err
	}

	// Loop through page templates
	for _, page := range pages {
		name := filepath.Base(page)

		// Skip base template
		if name == "base.tmpl" {
			continue
		}

		// Parse the base template first
		ts, err := template.New(name).Funcs(functions).ParseFiles("ui/html/base.tmpl")
		if err != nil {
			return err
		}

		// Parse the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return err
		}

		// Add to cache
		cache[name] = ts
	}

	templateCache = cache
	return nil
}

// getTemplate returns a cached template by name
func getTemplate(name string) (*template.Template, error) {
	if templateCache == nil {
		if err := initTemplateCache(); err != nil {
			return nil, err
		}
	}

	tmpl, ok := templateCache[name]
	if !ok {
		return nil, ErrTemplateNotFound
	}

	return tmpl, nil
}
