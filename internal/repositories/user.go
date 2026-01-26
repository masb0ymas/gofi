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
	"github.com/maxrichie5/go-sqlfmt/sqlfmt"
)

type UserRepository struct {
	BaseRepository
}

func (r UserRepository) Count() (int64, error) {
	return r.BaseRepository.countExec(r.DB)
}

func (r UserRepository) List(opts *QueryOptions) ([]*models.User, PaginationMetadata, error) {
	return r.listExec(r.DB, opts)
}

func (r UserRepository) listExec(exc Executor, opts *QueryOptions) ([]*models.User, PaginationMetadata, error) {
	selectFields := `"u"."id", "u"."created_at", "u"."updated_at", "u"."deleted_at", "u"."first_name", "u"."last_name", "u"."email", "u"."phone", "u"."active_at", "u"."blocked_at", "u"."role_id", "u"."upload_id"`
	selectRoleFields := `"r"."id", "r"."name", "r"."created_at", "r"."updated_at"`
	baseQuery := fmt.Sprintf(`
		SELECT %s, %s
		FROM "users" "u"
		LEFT JOIN "roles" "r" ON "u"."role_id" = "r"."id"
		WHERE "u"."deleted_at" IS NULL
	`, selectFields, selectRoleFields)

	var args []any
	argIndex := 1

	var queryBuilder strings.Builder
	queryBuilder.WriteString(baseQuery)

	orderBy := `"u"."created_at"`
	order := "DESC"

	// Whitelist of allowed columns for ORDER BY to prevent SQL injection
	allowedOrderByColumns := map[string]bool{
		`"u"."id"`:         true,
		`"u"."created_at"`: true,
		`"u"."updated_at"`: true,
		`"u"."first_name"`: true,
		`"u"."last_name"`:  true,
		`"u"."email"`:      true,
	}

	if opts.OrderBy != "" {
		if !allowedOrderByColumns[opts.OrderBy] {
			return nil, PaginationMetadata{}, errtrace.New("invalid order by column")
		}
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

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := exc.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Wrap(err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		role := &models.Role{}

		err = rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Phone,
			&user.ActiveAt,
			&user.BlockedAt,
			&user.RoleID,
			&user.UploadID,
			&role.ID,
			&role.Name,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, PaginationMetadata{}, errtrace.Errorf("error scanning row: %w", err)
		}

		user.Role = role
		users = append(users, user)
	}

	count, err := r.Count()
	if err != nil {
		return nil, PaginationMetadata{}, errtrace.Wrap(err)
	}

	return users, PaginationMetadata{Total: count}, nil
}

func (r UserRepository) Get(id uuid.UUID) (*models.User, error) {
	return r.getExec(r.DB, id)
}

func (r UserRepository) getExec(exc Executor, id uuid.UUID) (*models.User, error) {
	selectFields := `"u"."id", "u"."first_name", "u"."last_name", "u"."email", "u"."phone", "u"."active_at", "u"."blocked_at", "u"."role_id", "u"."upload_id", "u"."created_at", "u"."updated_at"`
	selectRoleFields := `"r"."id", "r"."name", "r"."created_at", "r"."updated_at"`
	query := fmt.Sprintf(`
		SELECT %s, %s
		FROM "users" "u"
		LEFT JOIN "roles" "r" ON "u"."role_id" = "r"."id"
		WHERE "u"."id" = $1 AND "u"."deleted_at" IS NULL;
	`, selectFields, selectRoleFields)

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.User{}
	role := &models.Role{}
	err := exc.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.ActiveAt,
		&user.BlockedAt,
		&user.RoleID,
		&user.UploadID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&role.ID,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	user.Role = role
	return user, nil
}

func (r UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	return r.getByIDExec(r.DB, id)
}

func (r UserRepository) getByIDExec(exc Executor, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT "u"."id", "u"."email", "u"."active_at", "u"."blocked_at", "u"."role_id"
		FROM "users" AS "u"
		WHERE "u"."id" = $1 AND
					"u"."active_at" IS NOT NULL AND
					"u"."blocked_at" IS NULL AND
					"u"."deleted_at" IS NULL;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.User{}
	err := exc.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.ActiveAt,
		&user.BlockedAt,
		&user.RoleID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return user, nil
}

func (r UserRepository) GetByEmail(email string) (*models.User, error) {
	return r.getByEmailExec(r.DB, email)
}

func (r UserRepository) getByEmailExec(exc Executor, email string) (*models.User, error) {
	query := `
		SELECT "u"."id", "u"."created_at", "u"."updated_at", "u"."deleted_at", "u"."first_name", "u"."last_name", "u"."email", "u"."phone", "u"."password", "u"."active_at", "u"."blocked_at", "u"."role_id", "u"."upload_id"
		FROM "users" AS "u"
		WHERE "u"."email" = $1 AND
				"u"."active_at" IS NOT NULL AND
				"u"."blocked_at" IS NULL AND
				"u"."deleted_at" IS NULL;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &models.User{}
	err := exc.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.ActiveAt,
		&user.BlockedAt,
		&user.RoleID,
		&user.UploadID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errtrace.Wrap(ErrRecordNotFound)
		default:
			return nil, errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return user, nil
}

func (r UserRepository) Insert(users ...*models.User) error {
	for _, user := range users {
		if user.Password != nil {
			if err := user.BeforeCreate(); err != nil {
				return errtrace.Wrap(err)
			}
		}
	}
	return r.InsertExec(r.DB, users...)
}

func (r UserRepository) InsertExec(exc Executor, users ...*models.User) error {
	if len(users) == 0 {
		return nil
	}

	columns := []string{
		"id",
		"first_name",
		"last_name",
		"email",
		"phone",
		"password",
		"active_at",
		"blocked_at",
		"role_id",
		"upload_id",
	}

	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]any, 0, len(users)*len(columns))

	for i, user := range users {
		values := []any{
			user.ID,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Phone,
			user.Password,
			user.ActiveAt,
			user.BlockedAt,
			user.RoleID,
			user.UploadID,
		}

		placeholders := make([]string, 0, len(values))
		for j := range columns {
			placeholders = append(placeholders, "$"+strconv.Itoa(i*len(columns)+j+1))
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
		valueArgs = append(valueArgs, values...)
	}

	query := fmt.Sprintf(`
		INSERT INTO "users" (%s)
		VALUES %s
		RETURNING "id", "created_at", "updated_at";
	`, strings.Join(columns[:], ", "), strings.Join(valueStrings, ", "))

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

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

	for _, user := range users {
		if !rows.Next() {
			return errtrace.New("error scanning row: no next row")
		}

		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return errtrace.Errorf("error scanning row: %w", err)
		}
	}

	return nil
}

func (r UserRepository) Update(id uuid.UUID, user *models.User) error {
	return r.UpdateExec(r.DB, id, user)
}

func (r UserRepository) UpdateExec(exc Executor, id uuid.UUID, user *models.User) error {
	query := `
		UPDATE "users"
		SET "first_name" = $1,
				"last_name" = $2,
				"email" = $3,
				"phone" = $4,
				"active_at" = $5,
				"blocked_at" = $6,
				"role_id" = $7,
				"upload_id" = $8,
				"updated_at" = now()
		WHERE "id" = $9;
	`

	if r.Config != nil && r.Config.Debug {
		fmt.Println()
		sqlfmt.PrettyPrint(query)
	}

	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		user.ActiveAt,
		user.BlockedAt,
		user.RoleID,
		user.UploadID,
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

func (r UserRepository) Delete(id uuid.UUID) error {
	return r.BaseRepository.deleteExec(r.DB, id)
}

func (r UserRepository) SoftDelete(id uuid.UUID) error {
	return r.BaseRepository.softDeleteExec(r.DB, id)
}

func (r UserRepository) Restore(id uuid.UUID) error {
	return r.BaseRepository.restoreExec(r.DB, id)
}
