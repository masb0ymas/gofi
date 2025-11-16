package models

import (
	"time"

	"gofi/internal/lib/argon2"

	"github.com/google/uuid"
)

type User struct {
	Base
	FirstName string     `db:"first_name" json:"first_name"`
	LastName  *string    `db:"last_name" json:"last_name"`
	Email     string     `db:"email" json:"email"`
	Password  *string    `db:"password" json:"password,omitempty"`
	Phone     *string    `db:"phone" json:"phone,omitempty"`
	ActiveAt  *time.Time `db:"active_at" json:"active_at,omitempty"`
	BlockedAt *time.Time `db:"blocked_at" json:"blocked_at,omitempty"`
	RoleID    uuid.UUID  `db:"role_id" json:"role_id"`
	UploadID  *uuid.UUID `db:"upload_id" json:"upload_id,omitempty"`
	// Relation
	Role   *Role   `json:"role,omitempty"`
	Upload *Upload `json:"upload,omitempty"`
}

func (entity *User) BeforeCreate() (err error) {
	argon2 := argon2.New()

	if entity.Password != nil {
		password, err := argon2.Generate(*entity.Password)
		if err != nil {
			return err
		}

		entity.Password = &password
	}

	return
}
