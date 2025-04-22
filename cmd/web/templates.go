package main

import (
	"errors"
	"html/template"
	"path/filepath"
	"sync"
)

// ErrTemplateNotFound is returned when a template is not found in the cache
var ErrTemplateNotFound = errors.New("template not found")

// templateCache holds the parsed templates
var templateCache map[string]*template.Template
var templateMutex sync.RWMutex

// testTemplateDir is used in tests to override the default template directory
var testTemplateDir string

// getTemplatePath returns the path to the template directory
func getTemplatePath(name string) string {
	if testTemplateDir != "" {
		return filepath.Join(testTemplateDir, name)
	}
	return filepath.Join("ui/html", name)
}

// initTemplateCache initializes the template cache
func initTemplateCache() error {
	cache := map[string]*template.Template{}

	// Get all page and partial templates
	pages, err := filepath.Glob(getTemplatePath("*.tmpl"))
	if err != nil {
		return err
	}
	partials, err := filepath.Glob(getTemplatePath("partials/*.tmpl"))
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
		ts, err := template.New(name).ParseFiles(getTemplatePath("base.tmpl"))
		if err != nil {
			return err
		}

		// Parse all partial templates
		ts, err = ts.ParseFiles(partials...)
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

	// Also cache partial templates individually
	for _, partial := range partials {
		name := filepath.Base(partial)
		ts, err := template.ParseFiles(partial)
		if err != nil {
			return err
		}
		cache["partials/"+name] = ts
	}

	templateCache = cache
	return nil
}

// getTemplate returns a template from the cache
func getTemplate(name string) (*template.Template, error) {
	templateMutex.RLock()
	defer templateMutex.RUnlock()

	tmpl, ok := templateCache[name]
	if !ok {
		return nil, ErrTemplateNotFound
	}

	return tmpl, nil
}
