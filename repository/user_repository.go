package repository

import (
	"goarif-api/database/model"
	"goarif-api/lib"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*Repository[model.User]
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		Repository: NewRepository[model.User](db, "user"),
	}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM user WHERE email = $1 AND deleted_at IS NULL"
	err := r.db.Get(&user, query, email)
	return &user, err
}

func (r *UserRepository) Create(user *model.User) error {
	hash, err := lib.Hash(*user.Password)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user (
			id, created_at, updated_at, fullname, email, password, is_active, role_id
		) VALUES (
			:id, :created_at, :updated_at, :fullname, :email, :password, :is_active, :role_id
		)
	`
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"id":         user.ID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"fullname":   user.Fullname,
		"email":      user.Email,
		"password":   hash,
		"is_active":  user.IsActive,
		"role_id":    user.RoleID,
	})
	return err
}

func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE user SET
			updated_at = :updated_at,
			fullname = :fullname,
			email = :email,
		WHERE id = :id AND deleted_at IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         user.ID,
		"updated_at": user.UpdatedAt,
		"fullname":   user.Fullname,
		"email":      user.Email,
	})
	return err
}
