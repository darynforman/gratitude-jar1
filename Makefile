# Variables
BINARY_NAME=gratitude-jar
GO=go
GOFMT=gofmt
GOLINT=golint
GOVET=go vet
GOBUILD=$(GO) build
GOTEST=$(GO) test
GORUN=$(GO) run
GOFILES=$(shell find . -name '*.go' -not -path "./vendor/*" -not -name "*_test.go")
GOFMTFILES=$(shell find . -name '*.go' -not -path "./vendor/*" -not -path "./node_modules/*")

# Build flags
LDFLAGS=-ldflags "-w -s"

.PHONY: all build run test clean fmt lint vet help

# Default target
all: build

# Build the application
build:
	@echo "Building..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/web

# Run the application
run:
	@echo "Running..."
	$(GORUN) $(shell find ./cmd/web -name '*.go' -not -name '*_test.go')

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)

# Format Go code
fmt:
	@echo "Formatting..."
	$(GOFMT) -w $(GOFMTFILES)

# Run linter
lint:
	@echo "Linting..."
	$(GOLINT) ./...

# Run vet
vet:
	@echo "Vetting..."
	$(GOVET) ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod tidy

# Generate database migrations
migrate:
	@echo "Generating migrations..."
	migrate create -ext sql -dir migrations -seq $(name)

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	migrate -path migrations -database "postgres://gratitude_user:gratitude123@localhost:5432/gratitude_jar?sslmode=disable" up

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	migrate -path migrations -database "postgres://gratitude_user:gratitude123@localhost:5432/gratitude_jar?sslmode=disable" down

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	$(GO) install golang.org/x/lint/golint@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all checks (fmt, lint, vet, test)
check: fmt lint vet test

# Show help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make fmt         - Format Go code"
	@echo "  make lint        - Run linter"
	@echo "  make vet         - Run vet"
	@echo "  make deps        - Install dependencies"
	@echo "  make migrate     - Generate new migration (use name=NAME)"
	@echo "  make migrate-up  - Run migrations up"
	@echo "  make migrate-down - Run migrations down"
	@echo "  make dev-setup   - Set up development environment"
	@echo "  make check       - Run all checks (fmt, lint, vet, test)"
	@echo "  make help        - Show this help message" 