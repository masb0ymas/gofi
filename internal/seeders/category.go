package seeders

import (
	"context"
	"database/sql"

	"gofi/internal/lib/constant"
)

type CategorySeeder struct {
	DB *sql.DB
}

func (s CategorySeeder) Name() string {
	return "category"
}

// Seed inserts the fixed categories; existing ids are left untouched.
func (s CategorySeeder) Seed() error {
	categories := []struct{ id, name string }{
		{constant.CategoryPublic, "Public"},
		{constant.CategoryPrivate, "Private"},
	}

	return runSeed(s.DB, func(ctx context.Context, tx *sql.Tx) error {
		for _, category := range categories {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO "categories" ("id", "name") VALUES ($1, $2) ON CONFLICT ("id") DO NOTHING`,
				category.id, category.name,
			); err != nil {
				return err
			}
		}
		return nil
	})
}
