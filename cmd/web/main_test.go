package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupTestEnvironment sets up the test environment
func setupTestEnvironment(t *testing.T) {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Set the template directory relative to the project root
	SetTemplateDir(filepath.Join(projectRoot, "ui/html"))
}

func TestMain(m *testing.M) {
	// Setup code
	setupTestEnvironment(nil)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// TestHomeHandler tests the home page handler
func TestHomeHandler(t *testing.T) {
	// Create a test request
	req := httptest.NewRequest("GET", "/", nil)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Create the handler
	handler := http.HandlerFunc(home)

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the status code
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Welcome to Gratitude Jar"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("got %q, expected to contain %q", rr.Body.String(), expected)
	}
}
