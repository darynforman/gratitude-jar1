# Gratitude Jar Application

A Go web application for managing gratitude entries with PostgreSQL database integration.

## Overview

The Gratitude Jar application is a web application that provides the following features:
- Create and view gratitude entries
- Persistent storage in PostgreSQL
- Simple and intuitive user interface

## Technologies

- Go (v1.23+)
- PostgreSQL database (v17)
- migrate tool for database migrations
- HTML templates with Go's template package
- CSS for styling

## Project Structure

```
gratitude-jar/
├── .gitignore
├── Makefile
├── go.mod
├── go.sum
├── cmd/
│   └── web/
│       ├── main.go              # Application entry point
│       ├── handlers.go          # HTTP handlers
│       ├── middleware.go        # HTTP middleware
│       ├── render.go            # Template rendering
│       ├── routes.go            # HTTP routing
│       ├── server.go            # HTTP server configuration
│       ├── templates.go         # Template loading
│       └── template_data.go     # Template data structures
├── internal/
│   ├── data/
│   │   ├── models.go            # Database models
│   │   └── gratitude.go         # Gratitude model and validation
│   └── validator/
│       └── validator.go         # Input validation
├── migrations/                  # Database migrations
│   ├── 000001_create_gratitude_table.up.sql
│   └── 000001_create_gratitude_table.down.sql
└── ui/
    ├── html/                    # HTML templates
    │   ├── base.tmpl
    │   ├── home.tmpl
    │   ├── gratitude.tmpl
    │   └── gratitudes.tmpl
    └── static/                  # Static assets
        ├── css/
        │   └── styles.css
        └── js/
            └── script.js
```

## Setup Instructions

### Prerequisites

1. Go 1.23 or later
2. PostgreSQL database
3. migrate CLI tool for database migrations

### Database Setup

1. Create a new PostgreSQL database:

```sql
CREATE DATABASE gratitude_jar;
CREATE USER gratitude_user WITH PASSWORD 'gratitude123';
GRANT ALL PRIVILEGES ON DATABASE gratitude_jar TO gratitude_user;
```

2. Run the database migrations:

```bash
make migrate-up
```

### Running the Application

1. Set up your development environment:

```bash
make dev-setup
```

2. Install dependencies:

```bash
make deps
```

3. Run the application:

```bash
make run
```

4. Open your browser and navigate to http://localhost:4000

## Development

### Available Makefile Commands

- `make run` - Run the application
- `make build` - Build the application
- `make test` - Run tests
- `make fmt` - Format Go code
- `make lint` - Run linter
- `make vet` - Vet Go code
- `make clean` - Clean build artifacts
- `make migrate name=migration_name` - Create a new migration
- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Revert the last migration
- `make dev-setup` - Set up development environment
- `make check` - Run all checks (fmt, lint, vet, test)

## Database Configuration

The application uses the following database configuration:
- Database name: `gratitude_jar`
- Username: `gratitude_user`
- Password: `gratitude123`
- Host: `localhost`
- Port: `5432`
- SSL Mode: `disable`

Connection string: `postgres://gratitude_user:gratitude123@localhost:5432/gratitude_jar?sslmode=disable`

## License

This project is developed for CMPS3162 Homework 2. 