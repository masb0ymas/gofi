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

type UserOAuthRepository struct {
	DB *sql.DB
}

func (r UserOAuthRepository) GetByUserProvider(userID, provider string) (*models.UserOAuth, error) {
	return r.getByUserProviderExec(r.DB, userID, provider)
}

func (r UserOAuthRepository) getByUserProviderExec(exc Executor, userID, provider string) (*models.UserOAuth, error) {
	query := `
		SELECT "id", "user_id", "provider", "access_token", "refresh_token", "expires_at"
		FROM "users_oauth"
		WHERE "user_id" = $1 AND "provider" = $2
		LIMIT 1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	userOAuth := &models.UserOAuth{}
	err := exc.QueryRowContext(ctx, query, userID, provider).Scan(
		&userOAuth.ID,
		&userOAuth.UserID,
		&userOAuth.Provider,
		&userOAuth.AccessToken,
		&userOAuth.RefreshToken,
		&userOAuth.ExpiresAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return userOAuth, nil
}

func (r UserOAuthRepository) Insert(usersOAuths ...*models.UserOAuth) error {
	return r.InsertExec(r.DB, usersOAuths...)
}

func (r UserOAuthRepository) InsertExec(exc Executor, usersOAuths ...*models.UserOAuth) error {
	if len(usersOAuths) == 0 {
		return nil
	}

	columns := []string{"id", "user_id", "provider", "access_token", "refresh_token", "expires_at"}

	valueStrings := make([]string, 0, len(usersOAuths))
	valueArgs := make([]any, 0, len(usersOAuths)*len(columns))

	for i, userOAuth := range usersOAuths {
		values := []any{userOAuth.ID, userOAuth.UserID, userOAuth.Provider, userOAuth.AccessToken, userOAuth.RefreshToken, userOAuth.ExpiresAt}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "users_oauth" (%s)
		VALUES %s
		RETURNING "id";
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

	for _, userOAuth := range usersOAuths {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&userOAuth.ID); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}
