package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"gofi/internal/models"

	"braces.dev/errtrace"
)

type SessionRepository struct {
	DB *sql.DB
}

func (r SessionRepository) GetByToken(token string) (*models.Session, error) {
	return r.getByTokenExec(r.DB, token)
}

func (r SessionRepository) getByTokenExec(exc Executor, token string) (*models.Session, error) {
	query := `
		SELECT "s"."id", "s"."created_at", "s"."updated_at", "s"."user_id", "s"."token"
		FROM "sessions" "s"
		WHERE "s"."token" = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	session := &models.Session{}
	err := exc.QueryRowContext(ctx, query, token).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt, &session.UserID, &session.Token)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return session, nil
}
