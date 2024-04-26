package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/env"
	"github.com/WieseChristoph/go-oauth2-backend/internal/utils/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type DB struct {
	*sql.DB
}

var (
	host       = env.GetEnvOrFatal("DB_HOST")
	port       = env.GetEnvOrFatal("DB_PORT")
	database   = env.GetEnvOrFatal("DB_DATABASE")
	password   = env.GetEnvOrFatal("DB_PASSWORD")
	username   = env.GetEnvOrFatal("DB_USERNAME")
	dbInstance *DB
)

func New() *DB {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database. Err: %v", err)
	}

	dbInstance = &DB{db}
	return dbInstance
}

func (db *DB) Migrate() error {
	log.Infoln("Migrating database...")

	// Create migrations table if it doesn't exist
	createMigrationsTableQuery := `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`
	_, err := db.Exec(createMigrationsTableQuery)
	if err != nil {
		return err
	}

	// Get all migrations from database (migration that have already been run)
	rows, err := db.Query(`
		SELECT name FROM migration;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Put migration names in slice
	var migrations []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}

		migrations = append(migrations, name)
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// Get all migration files in migrations folder
	files, err := os.ReadDir(path + "/migrations")
	if err != nil {
		return err
	}

	migrationsDone := 0
	for _, file := range files {
		// Skip file if it has already been run
		if slices.Contains(migrations, file.Name()) {
			continue
		}

		log.Infof("Running migration: %s", file.Name())

		// Read migration file
		migration, err := os.ReadFile("./migrations/" + file.Name())
		if err != nil {
			return err
		}

		// Run migration
		_, err = db.Exec(string(migration))
		if err != nil {
			return err
		}

		// Insert migration into database
		insertMigrationQuery := "INSERT INTO migration (name) VALUES ($1)"
		_, err = db.Exec(insertMigrationQuery, file.Name())
		if err != nil {
			return err
		}

		migrationsDone++
	}

	log.Infof("Migrated %d files", migrationsDone)

	return nil
}

func (db *DB) Close() error {
	if db != nil {
		return db.Close()
	}

	return nil
}

func (db *DB) IsHealthy() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}
