package repositories

import (
	"context"
	"database/sql"
)

type Executor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repositories struct {
	Role RoleRepository
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
