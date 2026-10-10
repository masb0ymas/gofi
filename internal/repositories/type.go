package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

// Executor is satisfied by both *sql.DB and *sql.Tx so repository methods can
// run inside or outside a transaction.
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// rowScanner is the scan surface shared by *sql.Row and *sql.Rows, so a single
// scan function serves both single-row and multi-row reads.
type rowScanner interface {
	Scan(dest ...any) error
}

// withTx runs fn inside a transaction, committing when fn returns nil and
// rolling back otherwise. fn receives the transaction as an Executor, so the
// same SQL runs standalone or atomically.
func (r BaseRepository) withTx(ctx context.Context, fn func(tx Executor) error) error {
	return lib.WithTransaction(ctx, r.DB, func(tx *sql.Tx) error { return fn(tx) })
}

// requireAffected turns a write that touched no row into apperr.ErrNotFound,
// so callers can tell "updated nothing" apart from a successful write.
func requireAffected(res sql.Result, what string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%w: %s", apperr.ErrNotFound, what)
	}
	return nil
}

// ListParams carries the list controls bound from the request query:
// pagination plus a search term and a "<column>:asc|desc" sort value.
type ListParams struct {
	Offset int
	Limit  int
	Search string
	Sorted string
}

// execAffected runs a statement and returns how many rows it changed. Cascades
// use it to report what a delete removed.
func (r BaseRepository) execAffected(ctx context.Context, ex Executor, query string, args ...any) (int64, error) {
	res, err := r.execContext(ctx, ex, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
