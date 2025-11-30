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
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type SessionRepository struct {
	DB *sql.DB
}

func (r SessionRepository) Count() (int64, error) {
	return r.countExec(r.DB)
}

func (r SessionRepository) countExec(exc Executor) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM "sessions";
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int64
	err := exc.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, errtrace.Errorf("error scanning row: %w", err)
	}

	return count, nil
}

func (r SessionRepository) List(opts *QueryOptions) ([]*models.Session, PaginationMetadata, error) {
	return r.listExec(r.DB, opts)
}

func (r SessionRepository) listExec(exc Executor, opts *QueryOptions) ([]*models.Session, PaginationMetadata, error) {
	if opts == nil {
		opts = &QueryOptions{}
	}

	selectFields := `"s"."id", "s"."created_at", "s"."updated_at", "s"."user_id", "s"."expires_at", "s"."ip_address", "s"."user_agent"`
	baseQuery := fmt.Sprintf(`
		SELECT %s
		FROM "sessions" "s"
	`, selectFields)

	var args []any
	argIndex := 1

	var queryBuilder strings.Builder
	queryBuilder.WriteString(baseQuery)

	orderBy := `"s"."created_at"`
	order := "DESC"

	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}

	if opts.Order != "" {
		upperOrder := strings.ToUpper(opts.Order)
		if upperOrder != "ASC" && upperOrder != "DESC" {
			return nil, PaginationMetadata{}, errtrace.New("invalid order")
		}
		order = upperOrder
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", orderBy, order))

	if opts.Limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
		args = append(args, opts.Limit)
		argIndex++
	}

	if opts.Offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
		args = append(args, opts.Offset)
		argIndex++
	}

	query := queryBuilder.String()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exc.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Errorf("error querying rows: %w", err)
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		session := &models.Session{}
		if err := rows.Scan(
			&session.ID,
			&session.CreatedAt,
			&session.UpdatedAt,
			&session.UserID,
			&session.ExpiresAt,
			&session.IPAddress,
			&session.UserAgent,
		); err != nil {
			return nil, PaginationMetadata{}, errtrace.Errorf("error scanning row: %w", err)
		}
		sessions = append(sessions, session)
	}

	count, err := r.countExec(exc)
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Errorf("error counting rows: %w", err)
	}

	return sessions, PaginationMetadata{
		Total: count,
	}, nil
}

func (r SessionRepository) GetByUserID(userID uuid.UUID) (*models.Session, error) {
	return r.getByUserIDExec(r.DB, userID)
}

func (r SessionRepository) getByUserIDExec(exc Executor, userID uuid.UUID) (*models.Session, error) {
	query := `
		SELECT "id", "user_id", "token", "expires_at"
		FROM "sessions"
		WHERE "user_id" = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	session := &models.Session{}
	err := exc.QueryRowContext(ctx, query, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
	)
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

func (r SessionRepository) GetByToken(token string) (*models.Session, error) {
	return r.getByTokenExec(r.DB, token)
}

func (r SessionRepository) getByTokenExec(exc Executor, token string) (*models.Session, error) {
	query := `
		SELECT "s"."id", "s"."created_at", "s"."updated_at", "s"."user_id", "s"."token", "s"."expires_at", "s"."ip_address", "s"."user_agent"
		FROM "sessions" "s"
		WHERE "s"."token" = $1 AND "s"."expires_at" > now();
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	session := &models.Session{}
	err := exc.QueryRowContext(ctx, query, token).Scan(
		&session.ID,
		&session.CreatedAt,
		&session.UpdatedAt,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.IPAddress,
		&session.UserAgent,
	)
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
	return r.InsertExec(r.DB, session...)
}

func (r SessionRepository) InsertExec(exc Executor, session ...*models.Session) error {
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

func (r SessionRepository) Delete(userID uuid.UUID, token string) error {
	return r.deleteExec(r.DB, userID, token)
}

func (r SessionRepository) deleteExec(exc Executor, userID uuid.UUID, token string) error {
	query := `
		DELETE FROM "sessions"
		WHERE "user_id" = $1 AND "token" = $2;
	`

	args := []any{userID, token}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}
