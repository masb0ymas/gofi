package repository

import (
	"goarif-api/database/model"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	*Repository[model.Role]
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{
		Repository: NewRepository[model.Role](db, "role"),
	}
}

func (r *RoleRepository) Create(role *model.Role) error {
	query := `
		INSERT INTO role (
			id, created_at, updated_at, name
		) VALUES (
			:id, :created_at, :updated_at, :name
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         role.ID,
		"created_at": role.CreatedAt,
		"updated_at": role.UpdatedAt,
		"name":       role.Name,
	})
	return err
}

func (r *RoleRepository) Update(role *model.Role) error {
	query := `
		UPDATE role SET
			updated_at = :updated_at,
			name = :name,
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         role.ID,
		"updated_at": role.UpdatedAt,
		"name":       role.Name,
	})
	return err
}
