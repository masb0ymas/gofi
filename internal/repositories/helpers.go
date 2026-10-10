package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gofi/internal/lib/apperr"

	"github.com/jackc/pgx/v5/pgconn"
)

// translateNotFound maps sql.ErrNoRows onto the typed app error so services
// and handlers can handle misses uniformly.
func translateNotFound(table string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.New(apperr.KindNotFound, "record not found in %s", table)
	}
	return err
}

// Postgres SQLSTATE codes the repositories translate into 409s.
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
)

// translateUnique maps a unique-constraint violation onto a typed 409
// carrying a human message (e.g. "organization name already exists"); any
// other error passes through unchanged.
func translateUnique(err error, message string) error {
	return translateCode(err, pgUniqueViolation, message)
}

// translateForeignKey maps a foreign-key violation onto a typed 409, so a
// permanent delete blocked by referencing rows reads as a conflict rather than
// a server error.
func translateForeignKey(err error, message string) error {
	return translateCode(err, pgForeignKeyViolation, message)
}

func translateCode(err error, code, message string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == code {
		return apperr.New(apperr.KindConflict, "%s", message)
	}
	return err
}

// escapeLikePattern escapes ILIKE wildcards (%, _) so a search term matches
// literally; the backslash is Postgres' default LIKE escape character.
func escapeLikePattern(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

// resolveOrder resolves a "column:asc|desc" sort value against a whitelist.
// ORDER BY can't take bind parameters, so untrusted input only ever selects a
// fixed column; anything unknown falls back to the repository's default.
func resolveOrder(sorted string, columns map[string]string, fallback string) string {
	field, dir, _ := strings.Cut(sorted, ":")
	col, ok := columns[field]
	if !ok {
		return fallback
	}
	if dir == "asc" {
		return col + " ASC"
	}
	return col + " DESC"
}

// whereBuilder accumulates AND-ed list conditions with positional args.
// Conditions reference their own arg with a single %d (or %[1]d when the same
// value appears more than once), which is replaced by its $n index.
type whereBuilder struct {
	conds []string
	args  []any
}

func newWhere(base string) *whereBuilder {
	return &whereBuilder{conds: []string{base}}
}

func (w *whereBuilder) add(cond string, arg any) {
	w.args = append(w.args, arg)
	w.conds = append(w.conds, fmt.Sprintf(cond, len(w.args)))
}

// search adds a case-insensitive substring match over columns.
func (w *whereBuilder) search(term string, columns ...string) {
	if term == "" {
		return
	}
	ors := make([]string, len(columns))
	for i, col := range columns {
		ors[i] = col + " ILIKE $%[1]d"
	}
	w.add("("+strings.Join(ors, " OR ")+")", "%"+escapeLikePattern(term)+"%")
}

// page renders the WHERE body and appends LIMIT/OFFSET args, returning the
// clause, the full arg list and the placeholders for LIMIT and OFFSET.
func (w *whereBuilder) page(p ListParams) (where string, args []any, limit, offset int) {
	args = append(w.args, p.Limit, p.Offset)
	return strings.Join(w.conds, " AND "), args, len(args) - 1, len(args)
}

// countRows returns how many rows match the builder's conditions, for the
// list metadata. table is a caller-supplied literal, never request input.
func (r BaseRepository) countRows(ctx context.Context, ex Executor, table string, w *whereBuilder) (int64, error) {
	var total int64
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, table, strings.Join(w.conds, " AND "))
	err := r.queryRowContext(ctx, ex, query, w.args...).Scan(&total)
	return total, err
}

// collectRows scans every row with scan, always returning a non-nil slice so
// empty lists serialize as [].
func collectRows[T any](rows *sql.Rows, scan func(rowScanner) (T, error)) ([]T, error) {
	defer rows.Close()

	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
