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

const insertRole = `
	INSERT INTO "role" (name) 
	VALUES ($1)
	RETURNING id, created_at, updated_at
`

func (repo *RoleRepository) CreateRole(ctx context.Context, r *entity.Role) (*entity.Role, error) {
	var (
		lastInsertID uuid.UUID
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := repo.db.QueryRowContext(ctx, insertRole, r.Name).
		Scan(&lastInsertID, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("error inserting role: %w", err)
	}

	r.ID = lastInsertID
	r.CreatedAt = createdAt
	r.UpdatedAt = updatedAt

	return r, nil
}

const findRoleById = `
	SELECT * FROM "role" 
	WHERE id=$1
`

func (repo *RoleRepository) GetRole(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var r entity.Role

	err := repo.db.GetContext(ctx, &r, findRoleById, id)
	if err != nil {
		return nil, fmt.Errorf("error getting role: %v", err)
	}

	return &r, nil
}

const findAllRole = `
	SELECT * FROM "role"
`

func (repo *RoleRepository) ListRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role

	err := repo.db.SelectContext(ctx, &roles, findAllRole)
	if err != nil {
		return nil, fmt.Errorf("error listing roles: %v", err)
	}

	return roles, nil
}

const updateRole = `
	UPDATE "role" SET name=:name, updated_at=:updated_at 
	WHERE id=:id
`

func (repo *RoleRepository) UpdateRole(ctx context.Context, r *entity.Role) (*entity.Role, error) {
	_, err := repo.db.NamedExecContext(ctx, updateRole, r)
	if err != nil {
		return nil, fmt.Errorf("error updating role: %v", err)
	}

	return r, nil
}

const deleteRole = `
	DELETE FROM "role" 
	WHERE id=?
`

func (repo *RoleRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, deleteRole, id)
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}

	return nil
}
