package repository

import (
	"context"
	"fmt"
	"gofi/database/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (repo *SessionRepository) CreateSession(ctx context.Context, s *entity.Session) (*entity.Session, error) {
	var (
		lastInsertID uuid.UUID
		createdAt    time.Time
		updatedAt    time.Time
	)

	const query_insert = `
		INSERT INTO "session" (user_id, token, expired_at) 
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := repo.db.QueryRowContext(ctx, query_insert, s.UserID, s.Token, s.ExpiredAt).
		Scan(&lastInsertID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	s.ID = lastInsertID
	s.CreatedAt = createdAt
	s.UpdatedAt = updatedAt

	return s, nil
}

func (repo *SessionRepository) GetSession(ctx context.Context, id uuid.UUID, token string) (*entity.Session, error) {
	var s entity.Session

	const query_find_one = `
		SELECT * FROM "session" 
		WHERE user_id=$1 AND token=$2
	`

	err := repo.db.GetContext(ctx, &s, query_find_one, id, token)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %v", err)
	}

	return &s, nil
}

func (repo *SessionRepository) ListSessions(ctx context.Context) ([]entity.Session, error) {
	var sessions []entity.Session

	const query_find_all = `
		SELECT * FROM "session"
	`

	err := repo.db.SelectContext(ctx, &sessions, query_find_all)
	if err != nil {
		return nil, fmt.Errorf("error listing session: %v", err)
	}

	return sessions, nil
}

func (repo *SessionRepository) UpdateSession(ctx context.Context, s *entity.Session) (*entity.Session, error) {
	const query_update = `
		UPDATE "session" SET user_id=:user_id, token=:token, expired_at=:expired_at, updated_at=:updated_at 
		WHERE id=:id
	`

	_, err := repo.db.NamedExecContext(ctx, query_update, s)
	if err != nil {
		return nil, fmt.Errorf("error updating session: %v", err)
	}

	return s, nil
}

func (repo *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	const query_delete = `
		DELETE FROM "session" 
		WHERE id=$1
	`

	_, err := repo.db.ExecContext(ctx, query_delete, id)
	if err != nil {
		return fmt.Errorf("error deleting session: %v", err)
	}

	return nil
}
