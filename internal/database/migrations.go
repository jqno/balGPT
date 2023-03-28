package database

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(conn database.Driver, appBaseDir string) error {
	migrationsPath := filepath.Join(appBaseDir, "db", "migrations")

	migrationUrl := fmt.Sprintf("file://%s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(migrationUrl, "postgres", conn)
	if err != nil {
		return fmt.Errorf("unable to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")

	return nil
}
