package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gofi/internal/models"

	"braces.dev/errtrace"
	"github.com/lib/pq"
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

func (r SessionRepository) Insert(session ...*models.Session) error {
	return r.insertExec(r.DB, session...)
}

func (r SessionRepository) insertExec(exc Executor, session ...*models.Session) error {
	if len(session) == 0 {
		return nil
	}

	columns := []string{"id", "user_id", "token", "expires_at", "ip_address", "user_agent"}

	valueStrings := make([]string, 0, len(session))
	valueArgs := make([]any, 0, len(session)*len(columns))

	for i, s := range session {
		values := []any{s.ID, s.UserID, s.Token, s.ExpiresAt, s.IPAddress, s.UserAgent}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "sessions" (%s)
		VALUES %s
		RETURNING "id", "created_at", "updated_at";
	`, strings.Join(columns[:], ", "), strings.Join(valueStrings, ", "))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exc.QueryContext(ctx, query, valueArgs...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errtrace.Wrap(ErrInsertDuplicate)
			}
		}
		return errtrace.Wrap(err)
	}
	defer rows.Close()

	for _, s := range session {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}
