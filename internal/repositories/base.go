package repositories

import (
	"context"
	"database/sql"
	"log/slog"
)

type BaseRepository struct {
	DB *sql.DB
}

// execContext runs a write statement after logging it at debug level. Only
// the SQL text is logged, never its arguments.
func (r BaseRepository) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	slog.Debug("query", "sql", query)
	return r.DB.ExecContext(ctx, query, args...)
}

// queryContext runs a read statement after logging it.
func (r BaseRepository) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	slog.Debug("query", "sql", query)
	return r.DB.QueryContext(ctx, query, args...)
}

// queryRowContext runs a single-row read after logging it. A miss surfaces at
// Scan as sql.ErrNoRows, which translateNotFound maps onto a typed 404.
func (r BaseRepository) queryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	slog.Debug("query", "sql", query)
	return r.DB.QueryRowContext(ctx, query, args...)
}
