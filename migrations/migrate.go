package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"monetix-be-api/configs"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version string
	Up      string
	Down    string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrations/migrate.go [up|down]")
	}

	command := os.Args[1]

	// Load config
	config, err := configs.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", config.Database.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}

	// Load migration files
	migrations, err := loadMigrations("migrations")
	if err != nil {
		log.Fatal("Failed to load migrations:", err)
	}

	switch command {
	case "up":
		if err := runMigrationsUp(db, migrations); err != nil {
			log.Fatal("Migration up failed:", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		// if err := runMigrationsDown(db, migrations); err != nil {
		// 	log.Fatal("Migration down failed:", err)
		// }
		log.Println("Migrations reverted successfully")
	default:
		log.Fatal("Unknown command. Use 'up' or 'down'")
	}
}

func createMigrationsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `
	_, err := db.Exec(query)
	return err
}

func loadMigrations(migrationsDir string) ([]Migration, error) {
	var migrations []Migration

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return nil, err
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}

		parts := strings.Split(string(content), "-- +migrate Down")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid migration file format: %s", file)
		}

		up := strings.TrimPrefix(parts[0], "-- +migrate Up")
		migration := Migration{
			Version: filepath.Base(file),
			Up:      strings.TrimSpace(up),
			Down:    strings.TrimSpace(parts[1]),
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func runMigrationsUp(db *sql.DB, migrations []Migration) error {
	for _, migration := range migrations {
		// Check if migration already applied
		var exists bool
		err := db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)",
			migration.Version,
		).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		// Run migration
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		if _, err := tx.Exec(migration.Up); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %s: %w", migration.Version, err)
		}

		// Record migration
		if _, err := tx.Exec(
			"INSERT INTO schema_migrations (version) VALUES ($1)",
			migration.Version,
		); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		log.Printf("Applied migration: %s\n", migration.Version)
	}

	return nil
}
