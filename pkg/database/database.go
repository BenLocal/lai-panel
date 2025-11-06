package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB(dbPath string, migrationsPath string) error {
	var err error
	db, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db.MustExec("PRAGMA foreign_keys = ON")
	db.MustExec("PRAGMA journal_mode = WAL")

	err = runMigrations(db, migrationsPath)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("migrations up to date")
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	DB = db
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func GetDB() *sqlx.DB {
	return DB
}
