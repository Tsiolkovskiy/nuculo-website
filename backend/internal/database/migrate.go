package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs database migrations
func RunMigrations(config *Config) error {
	// Create migration instance
	m, err := migrate.New(
		"file://migrations",
		config.ConnectionString(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back database migrations
func RollbackMigrations(config *Config, steps int) error {
	m, err := migrate.New(
		"file://migrations",
		config.ConnectionString(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Steps(-steps); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Printf("✅ Rolled back %d migration(s) successfully", steps)
	return nil
}