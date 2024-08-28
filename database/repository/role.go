package repository

import (
	"context"
	"fmt"
	"gofi/database/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (repo *RoleRepository) CreateRole(ctx context.Context, r *entity.Role) (*entity.Role, error) {
	var (
		lastInsertID uuid.UUID
		createdAt    time.Time
		updatedAt    time.Time
	)

	const query_insert = `
		INSERT INTO "role" (name) 
		VALUES ($1)
		RETURNING id, created_at, updated_at
	`

	err := repo.db.QueryRowContext(ctx, query_insert, r.Name).
		Scan(&lastInsertID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting role: %w", err)
	}

	r.ID = lastInsertID
	r.CreatedAt = createdAt
	r.UpdatedAt = updatedAt

	return r, nil
}

func (repo *RoleRepository) GetRole(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var r entity.Role

	const query_find_one = `
		SELECT * FROM "role" 
		WHERE id=$1
	`

	err := repo.db.GetContext(ctx, &r, query_find_one, id)
	if err != nil {
		return nil, fmt.Errorf("error getting role: %v", err)
	}

	return &r, nil
}

func (repo *RoleRepository) ListRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role

	const query_find_all = `
		SELECT * FROM "role"
	`

	err := repo.db.SelectContext(ctx, &roles, query_find_all)
	if err != nil {
		return nil, fmt.Errorf("error listing roles: %v", err)
	}

	return roles, nil
}

func (repo *RoleRepository) UpdateRole(ctx context.Context, r *entity.Role) (*entity.Role, error) {
	const query_update = `
		UPDATE "role" SET name=:name, updated_at=:updated_at 
		WHERE id=:id
	`

	_, err := repo.db.NamedExecContext(ctx, query_update, r)
	if err != nil {
		return nil, fmt.Errorf("error updating role: %v", err)
	}

	return r, nil
}

func (repo *RoleRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	const query_delete = `
		DELETE FROM "role" 
		WHERE id=$1
	`

	_, err := repo.db.ExecContext(ctx, query_delete, id)
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}

	return nil
}
