package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Include is the parsed set of relation names a request asked to embed via
// ?include=role,organization,…. A dotted name (submission.assignment) embeds a
// relation on an embedded relation.
type Include map[string]struct{}

// ParseInclude splits a comma-separated include query value.
func ParseInclude(raw string) Include {
	set := Include{}
	for _, part := range strings.Split(raw, ",") {
		if part = strings.TrimSpace(part); part != "" {
			set[part] = struct{}{}
		}
	}
	return set
}

// Has reports whether name was requested.
func (i Include) Has(name string) bool {
	_, ok := i[name]
	return ok
}

// Relator loads relations in batches, so a list page embeds nested objects with
// a handful of IN queries rather than one per row.
type Relator struct {
	BaseRepository
}

// loadByID loads the rows of table whose id is in ids, keyed by id. table and
// columns are caller-supplied literals, never request input.
func loadByID[T any](
	ctx context.Context, r *Relator, table, columns string,
	scan func(rowScanner) (T, error), id func(T) uuid.UUID, ids []uuid.UUID,
) (map[uuid.UUID]T, error) {
	out := map[uuid.UUID]T{}
	if len(ids) == 0 {
		return out, nil
	}
	rows, err := r.queryContext(ctx, r.DB, fmt.Sprintf(`SELECT %s FROM %s WHERE id = ANY($1)`, columns, table), ids)
	if err != nil {
		return nil, err
	}
	list, err := collectRows(rows, scan)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		out[id(v)] = v
	}
	return out, nil
}
