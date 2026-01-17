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

type RefreshTokenRepository struct {
	DB     *sql.DB
	Config *config.ConfigApp
}

func (r RefreshTokenRepository) Get(userID uuid.UUID, token string) (*models.RefreshToken, error) {
	return r.getExec(r.DB, userID, token)
}

func (r RefreshTokenRepository) getExec(exc Executor, userID uuid.UUID, token string) (*models.RefreshToken, error) {
	query := `
		SELECT "id", "user_id", "token", "expires_at", "created_at", "revoked_at"
		FROM "refresh_tokens"
		WHERE "user_id" = $1 AND "token" = $2 AND "expires_at" > now() AND "revoked_at" IS NULL;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rt := &models.RefreshToken{}
	err := exc.QueryRowContext(ctx, query, userID, token).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Token,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.RevokedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return rt, nil
}

func (r RefreshTokenRepository) Insert(refreshTokens ...*models.RefreshToken) error {
	return r.InsertExec(r.DB, refreshTokens...)
}

func (r RefreshTokenRepository) InsertExec(exc Executor, refreshTokens ...*models.RefreshToken) error {
	if len(refreshTokens) == 0 {
		return nil
	}

	columns := []string{"id", "user_id", "token", "expires_at", "created_at", "revoked_at"}

	valueStrings := make([]string, 0, len(refreshTokens))
	valueArgs := make([]any, 0, len(refreshTokens)*len(columns))

	for i, refreshToken := range refreshTokens {
		values := []any{refreshToken.ID, refreshToken.UserID, refreshToken.Token, refreshToken.ExpiresAt, refreshToken.CreatedAt, refreshToken.RevokedAt}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "refresh_tokens" (%s)
		VALUES %s
		RETURNING "id", "created_at";
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

	for _, refreshToken := range refreshTokens {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&refreshToken.ID, &refreshToken.CreatedAt); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}

func (r RefreshTokenRepository) Update(refreshToken *models.RefreshToken) error {
	return r.updateExec(r.DB, refreshToken)
}

func (r RefreshTokenRepository) updateExec(exc Executor, refreshToken *models.RefreshToken) error {
	query := `
		UPDATE "refresh_tokens"
		SET "token" = $1, "expires_at" = $2, "revoked_at" = $3
		WHERE "id" = $4;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	args := []any{
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.RevokedAt,
		refreshToken.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}
