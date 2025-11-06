package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func runMigrations(db *sqlx.DB, migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	sourceDriver, err := (&file.File{}).Open("file://" + migrationsPath)
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"file", sourceDriver,
		"sqlite3", driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}
