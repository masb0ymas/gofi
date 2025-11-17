package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gofi/internal/models"

	"braces.dev/errtrace"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RoleRepository struct {
	BaseRepository
}

func (r RoleRepository) Count() (int64, error) {
	return r.BaseRepository.countExec(r.DB)
}

func (r RoleRepository) List(opts *QueryOptions) ([]*models.Role, PaginationMetadata, error) {
	return r.listExec(r.DB, opts)
}

func (r RoleRepository) listExec(exc Executor, opts *QueryOptions) ([]*models.Role, PaginationMetadata, error) {
	selectFields := `"id", "name", "created_at", "updated_at"`
	baseQuery := fmt.Sprintf(`
		SELECT %s 
		FROM "roles"
		WHERE "deleted_at" IS NULL
	`, selectFields)

	var args []any
	argIndex := 1

	var queryBuilder strings.Builder
	queryBuilder.WriteString(baseQuery)

	orderBy := `"created_at"`
	order := "DESC"

	if opts.OrderBy != "" {
		orderBy = opts.OrderBy
	}

	if opts.Order != "" {
		upperOrder := strings.ToUpper(opts.Order)
		if upperOrder != "ASC" && upperOrder != "DESC" {
			return nil, PaginationMetadata{}, errtrace.New("invalid order")
		}
		order = upperOrder
	}

	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", orderBy, order))

	if opts.Limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d", argIndex))
		args = append(args, opts.Limit)
		argIndex++
	}

	if opts.Offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET $%d", argIndex))
		args = append(args, opts.Offset)
		argIndex++
	}

	query := queryBuilder.String()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exc.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Wrap(err)
	}
	defer rows.Close()

	var roles []*models.Role
	for rows.Next() {
		role := &models.Role{}
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, PaginationMetadata{}, errtrace.Errorf("error scanning row: %w", err)
		}
		roles = append(roles, role)
	}

	count, err := r.Count()
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Wrap(err)
	}

	return roles, PaginationMetadata{Total: count}, nil
}

func (r RoleRepository) Get(id uuid.UUID) (*models.Role, error) {
	return r.getExec(r.DB, id)
}

func (r RoleRepository) getExec(exc Executor, id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT "id", "name", "created_at", "updated_at"
		FROM "roles"
		WHERE "id" = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	role := &models.Role{}
	err := exc.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return role, nil
}

func (r RoleRepository) Insert(roles ...*models.Role) error {
	return r.insertExec(r.DB, roles...)
}

func (r RoleRepository) insertExec(exc Executor, roles ...*models.Role) error {
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
	`, strings.Join(columns[:], ", "), strings.Join(valueStrings, ", "))

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

func (r RoleRepository) Update(id uuid.UUID, role *models.Role) error {
	return r.updateExec(r.DB, id, role)
}

func (r RoleRepository) updateExec(exc Executor, id uuid.UUID, role *models.Role) error {
	query := `
		UPDATE "roles"
		SET "name" = $1, "updated_at" = now()
		WHERE "id" = $2;
	`

	args := []any{
		role.Name,
		id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := exc.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrEditConflict
	}

	return nil
}

func (r RoleRepository) Delete(id uuid.UUID) error {
	return r.BaseRepository.deleteExec(r.DB, id)
}

func (r RoleRepository) SoftDelete(id uuid.UUID) error {
	return r.BaseRepository.softDeleteExec(r.DB, id)
}

func (r RoleRepository) Restore(id uuid.UUID) error {
	return r.BaseRepository.restoreExec(r.DB, id)
}
