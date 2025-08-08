# Makefile for GraphQL TypeScript-Go Backend

.PHONY: help setup generate run test clean

# Default target
help:
	@echo "Available commands:"
	@echo "  setup     - Install dependencies and generate code"
	@echo "  generate  - Generate GraphQL code from schema"
	@echo "  run       - Run the development server"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean generated files"

# Setup project
setup:
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Installing gqlgen..."
	go install github.com/99designs/gqlgen@latest
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate
	@echo "Setup complete!"

# Generate GraphQL code
generate:
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate

# Run development server
run:
	@echo "Starting GraphQL server..."
	go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf internal/graph/generated/
	rm -f internal/graph/model/models_gen.go

# Install dependencies
deps:
	go mod tidy
	go mod download