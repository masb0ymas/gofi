package seeders

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"gofi/internal/lib"
	"time"
)

// data holds the large seed sets exported from the old TypeORM backend
// (src/lib/constants/seed).
//
//go:embed data/*.json
var data embed.FS

// MinPasswordLen is the shortest default password accepted for seeded accounts.
const MinPasswordLen = 12

// ValidatePassword refuses the placeholder from older .env.example files and
// anything short, because the seeded admin account gets this password.
func ValidatePassword(pwd string) error {
	if len(pwd) < MinPasswordLen || pwd == "defaultPassword123" {
		return fmt.Errorf("must be a non-default password of at least %d characters when seeding", MinPasswordLen)
	}
	return nil
}

// All returns every seeder in the order of the old TypeORM seeder migrations:
// referenced rows first.
func All(db *sql.DB, defaultPassword string) []Seeder {
	return []Seeder{
		CategorySeeder{DB: db},
	}
}

// runSeed executes fn in one transaction, so a seeder either applies
// completely or not at all.
func runSeed(db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return lib.WithTransaction(ctx, db, func(tx *sql.Tx) error {
		// Serializes seeders across processes: replicas seeding on boot at the
		// same time would otherwise both pass the assignment NOT EXISTS check.
		if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext('kodingup:seed'))`); err != nil {
			return err
		}
		return fn(ctx, tx)
	})
}

func loadJSON[T any](name string) ([]T, error) {
	raw, err := data.ReadFile("data/" + name)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}

	var rows []T
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, fmt.Errorf("decode %s: %w", name, err)
	}
	return rows, nil
}
