package repositories

import (
	"context"

	"gofi/internal/models"
)

// UserRepository manages the public mirror of Authula's users. Authula owns
// the source rows, so the only write here is the upsert that follows it.
type UserRepository struct {
	BaseRepository
}

// Upsert copies an Authula user into the mirror, inserting or refreshing it.
func (r *UserRepository) Upsert(ctx context.Context, u models.User) error {
	return r.execUpsert(ctx, r.DB, u)
}

func (r *UserRepository) execUpsert(ctx context.Context, ex Executor, u models.User) error {
	_, err := r.execContext(ctx, ex, `
		INSERT INTO users (id, name, email, image)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, email = EXCLUDED.email, image = EXCLUDED.image, updated_at = now()`,
		u.ID, u.Name, u.Email, u.Image,
	)
	return err
}
