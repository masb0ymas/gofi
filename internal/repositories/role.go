package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gofi/internal/models"

	"braces.dev/errtrace"
	"github.com/lib/pq"
)

type RoleRepository struct {
	DB *sql.DB
}

func (r RoleRepository) Insert(roles ...*models.Role) error {
	return r.InsertExec(r.DB, roles...)
}

func (r RoleRepository) InsertExec(exc Executor, roles ...*models.Role) error {
	if len(roles) == 0 {
		return nil
	}

	columns := []string{"id", "name"}

	valueStrings := make([]string, 0, len(roles))
	valueArgs := make([]any, 0, len(roles)*len(columns))

	for i, role := range roles {
		values := []any{role.ID, role.Name}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "roles" (%s)
		VALUES %s
		RETURNING "id", "created_at", "updated_at";
	`, strings.Join(columns[:], ","), strings.Join(valueStrings, ","))

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

	for _, role := range roles {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}
