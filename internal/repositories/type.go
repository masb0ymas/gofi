package repositories

import (
	"database/sql"
	"fmt"

	"gofi/internal/lib/apperr"
)

// rowScanner is the scan surface shared by *sql.Row and *sql.Rows, so a single
// scan function serves both single-row and multi-row reads.
type rowScanner interface {
	Scan(dest ...any) error
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
