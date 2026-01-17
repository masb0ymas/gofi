package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gofi/internal/config"
	"gofi/internal/models"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/maxrichie5/go-sqlfmt/sqlfmt"
)

type UserVerifyAccountRepository struct {
	DB     *sql.DB
	Config *config.ConfigApp
}

func (r UserVerifyAccountRepository) Get(id uuid.UUID, token string) (*models.UserVerifyAccount, error) {
	return r.getExec(r.DB, id, token)
}

func (r UserVerifyAccountRepository) getExec(exc Executor, id uuid.UUID, token string) (*models.UserVerifyAccount, error) {
	query := `
		SELECT "id", "token", "expires_at"
		FROM "user_verify_accounts"
		WHERE "id" = $1 AND "token" = $2;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.UserVerifyAccount{}
	err := exc.QueryRowContext(ctx, query, id, token).Scan(&user.ID, &user.Token, &user.ExpiresAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return user, nil
}

func (r UserVerifyAccountRepository) Insert(users ...*models.UserVerifyAccount) error {
	return r.InsertExec(r.DB, users...)
}

func (r UserVerifyAccountRepository) InsertExec(exc Executor, users ...*models.UserVerifyAccount) error {
	if len(users) == 0 {
		return nil
	}

	columns := []string{"id", "token", "expires_at"}

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]any, 0, len(users)*len(columns))

	for i, user := range users {
		values := []any{user.ID, user.Token, user.ExpiresAt}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "user_verify_accounts" (%s)
		VALUES %s
		RETURNING "id";
	`, strings.Join(columns[:], ", "), strings.Join(valueStrings, ", "))

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

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

	for _, user := range users {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&user.ID); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}
