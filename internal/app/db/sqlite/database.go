package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/alexgolang/ishare-task/internal/app/db/sqlite/sqlc"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type Database struct {
	db      *sql.DB
	Queries *sqlc.Queries
}

func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := sqlc.New(db)

	return &Database{db: db, Queries: queries}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) RunMigrations() error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "migrations")

	if err := goose.Up(d.db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
