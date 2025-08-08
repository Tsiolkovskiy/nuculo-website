package database

import (
	"context"
	"fmt"
	"log"
)

// Initialize sets up the database connection and runs migrations
func Initialize() (*DB, error) {
	// Load configuration
	config := NewConfig()
	
	// Create database connection
	db, err := NewConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Test connection
	ctx := context.Background()
	if err := db.Health(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database health check failed: %w", err)
	}
	
	log.Println("✅ Database initialization completed successfully")
	return db, nil
}

// InitializeWithMigrations sets up the database and runs migrations
func InitializeWithMigrations() (*DB, error) {
	config := NewConfig()
	
	// Run migrations first
	if err := RunMigrations(config); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	
	// Then initialize connection
	return Initialize()
}