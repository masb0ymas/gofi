package repository

import (
	"gofi/src/app/database/model"

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

func (r *RoleRepository) Create(values *model.Role) error {
	query := `
		INSERT INTO "role" (
			"id", "created_at", "updated_at", "name"
		) VALUES (
			:id, :created_at, :updated_at, :name
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"created_at": values.CreatedAt,
		"updated_at": values.UpdatedAt,
		"name":       values.Name,
	})
	return err
}

func (r *RoleRepository) Update(values *model.Role) error {
	query := `
		UPDATE "role" SET
			"updated_at" = :updated_at,
			"name" = :name
		WHERE "id" = :id AND "deleted_at" IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"updated_at": values.UpdatedAt,
		"name":       values.Name,
	})
	return err
}
