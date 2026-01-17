package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"gofi/internal/config"
	"time"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/maxrichie5/go-sqlfmt/sqlfmt"
)

type BaseRepository struct {
	DB        *sql.DB
	TableName string
	Config    *config.ConfigApp
}

func (r BaseRepository) countExec(exc Executor) (int64, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM "%s";
	`, r.TableName)

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := exc.QueryRowContext(ctx, query)
	if row == nil {
		return 0, errtrace.New("error scanning row: no next row")
	}

	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, errtrace.Errorf("error scanning row: %w", err)
	}

	return count, nil
}

func (r BaseRepository) deleteExec(exc Executor, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM "%s" 
		WHERE "id" = $1;
	`, r.TableName)

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return errtrace.Wrap(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (r BaseRepository) softDeleteExec(exc Executor, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE "%s" 
		SET "deleted_at" = now() 
		WHERE "id" = $1;
	`, r.TableName)

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return errtrace.Wrap(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (r BaseRepository) restoreExec(exc Executor, id uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE "%s" 
		SET "deleted_at" = NULL 
		WHERE "id" = $1;
	`, r.TableName)

	args := []any{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return errtrace.Wrap(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
