package main

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/godoc_vfs"
	"golang.org/x/tools/godoc/vfs"

	_ "github.com/lib/pq"
)

//go:embed migrations
var migrations embed.FS
var database *sql.DB

func connectToDB() error {
	// Create the database driver
	db, err := sql.Open("postgres", "postgres://dom@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return err
	}

	// Run migrations
	if err := migrateDatabase(db); err != nil {
		return err
	}

	database = db
	return nil
}

func migrateDatabase(db *sql.DB) error {

	// Create the driver for the source and database
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database migration driver: %w", err)
	}
	source, err := godoc_vfs.WithInstance(vfs.FromFS(migrations), "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source driver: %w", err)
	}

	// initialise the migrate instance
	m, err := migrate.NewWithInstance("migrations", source, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create database migration instance: %w", err)
	}

	// Perform any up migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("failed to create migrate database: %w", err)
	}
	return nil
}
