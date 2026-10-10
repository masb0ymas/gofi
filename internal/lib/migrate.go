package lib

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"gofi/migrations"
)

// NewMigrator builds a golang-migrate instance over the embedded SQL files.
//
// Closing the returned migrator also closes db, so callers must pass a
// connection dedicated to migrations rather than the application pool.
func NewMigrator(db *sql.DB) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("load embedded migrations: %w", err)
	}

	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return nil, fmt.Errorf("build migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "pgx5", driver)
	if err != nil {
		return nil, fmt.Errorf("init migrator: %w", err)
	}

	return m, nil
}
