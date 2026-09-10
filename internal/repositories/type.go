package repositories

import (
	"context"
	"database/sql"
	"strings"

	"braces.dev/errtrace"
)

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type QueryOptions struct {
	Offset int64
	Limit  int64

	OrderBy string
	Order   string // asc | desc
}

type PaginationMetadata struct {
	Total int64 `json:"total"`
}

// buildOrderBy resolves the ORDER BY column and direction from user-supplied
// options. OrderBy must appear in allowedColumns (quoted, table-qualified as
// written in the query); Order is case-insensitively matched to ASC/DESC.
func buildOrderBy(opts *QueryOptions, allowedColumns map[string]bool, defaultOrderBy string) (string, string, error) {
	orderBy := defaultOrderBy
	order := "DESC"

	if opts.OrderBy != "" {
		if !allowedColumns[opts.OrderBy] {
			return "", "", errtrace.New("invalid order by column")
		}
		orderBy = opts.OrderBy
	}

	if opts.Order != "" {
		upperOrder := strings.ToUpper(opts.Order)
		if upperOrder != "ASC" && upperOrder != "DESC" {
			return "", "", errtrace.New("invalid order")
		}
		order = upperOrder
	}

	return orderBy, order, nil
}
