package migrations

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed upgrade
var migrationsUpgradeFS embed.FS

func RunMigrations(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	// 从 embed 的文件系统读取
	sourceDriver, err := iofs.New(migrationsUpgradeFS, "upgrade")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs", sourceDriver,
		"sqlite3", driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}
