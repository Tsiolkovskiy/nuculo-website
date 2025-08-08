#!/bin/bash

# Setup script for Go GraphQL server

echo "Installing Go dependencies..."
go mod tidy

echo "Installing gqlgen..."
go install github.com/99designs/gqlgen@latest

echo "Generating GraphQL code..."
go run github.com/99designs/gqlgen generate

echo "Setup complete!"
echo "To run the server: go run cmd/server/main.go"