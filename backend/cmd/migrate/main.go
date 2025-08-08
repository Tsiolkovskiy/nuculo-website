package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"backend/internal/database"
)

func main() {
	var (
		up       = flag.Bool("up", false, "Run migrations up")
		down     = flag.Bool("down", false, "Run migrations down")
		steps    = flag.Int("steps", 1, "Number of steps to rollback (only for down)")
	)
	flag.Parse()

	config := database.NewConfig()

	if *up {
		log.Println("Running migrations up...")
		if err := database.RunMigrations(config); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("✅ Migrations completed successfully")
		return
	}

	if *down {
		log.Printf("Rolling back %d migration(s)...", *steps)
		if err := database.RollbackMigrations(config, *steps); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("✅ Rollback completed successfully")
		return
	}

	// Default: show usage
	log.Println("Usage:")
	log.Println("  go run cmd/migrate/main.go -up          # Run migrations up")
	log.Println("  go run cmd/migrate/main.go -down -steps=2  # Rollback 2 migrations")
	log.Println("")
	log.Println("Environment variables:")
	log.Println("  DB_HOST     - Database host (default: localhost)")
	log.Println("  DB_PORT     - Database port (default: 5432)")
	log.Println("  DB_USER     - Database user (default: postgres)")
	log.Println("  DB_PASSWORD - Database password (default: postgres)")
	log.Println("  DB_NAME     - Database name (default: graphql_typescript_go)")
	log.Println("  DB_SSL_MODE - SSL mode (default: disable)")
}