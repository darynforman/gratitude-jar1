// Package main is the entry point for the Gratitude Jar application.
// It initializes the database connection, sets up routes, and starts the HTTP server.
package main

import (
	"database/sql"
	"log"

	"github.com/darynforman/gratitude-jar1/internal/config"
	"github.com/darynforman/gratitude-jar1/internal/data"
)

// application holds the application-wide dependencies and configuration
type application struct {
	config *config.Config
	models *data.Models
	DB     *sql.DB
}

var app *application

// main is the entry point of the application.
// It performs the following tasks:
// 1. Initializes the database connection and configuration
// 2. Creates an application instance
// 3. Starts the HTTP server
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create application instance
	app = &application{
		config: cfg,
		models: data.NewModels(config.DB),
		DB:     config.DB,
	}

	// Start the server
	startServer()
}
