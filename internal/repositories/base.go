package repositories

import (
	"context"
	"database/sql"

	"github.com/maxrichie5/go-sqlfmt/sqlfmt"
)

type BaseRepository struct {
	DB    *sql.DB
	Debug bool
}

func (r BaseRepository) debugQuery(query string) {
	if !r.Debug {
		return
	}
	formatted := sqlfmt.PrettyFormat(query)
	println(formatted)
}

// execContext runs a write statement against ex — the pool or an open
// transaction — after logging it.
func (r BaseRepository) execContext(ctx context.Context, ex Executor, query string, args ...any) (sql.Result, error) {
	r.debugQuery(query)
	return ex.ExecContext(ctx, query, args...)
}

// queryContext runs a read statement against ex after logging it.
func (r BaseRepository) queryContext(ctx context.Context, ex Executor, query string, args ...any) (*sql.Rows, error) {
	r.debugQuery(query)
	return ex.QueryContext(ctx, query, args...)
}

// queryRowContext runs a single-row read against ex after logging it. A miss
// surfaces at Scan as sql.ErrNoRows, which translateNotFound maps onto
// apperr.ErrNotFound.
func (r BaseRepository) queryRowContext(ctx context.Context, ex Executor, query string, args ...any) *sql.Row {
	r.debugQuery(query)
	return ex.QueryRowContext(ctx, query, args...)
}
