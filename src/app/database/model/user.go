package model

import (
	"gofi/src/lib"
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Fullname    string     `db:"fullname" json:"fullname"`
	Email       string     `db:"email" json:"email"`
	Password    *string    `db:"password" json:"password,omitempty"`
	Phone       *string    `db:"phone" json:"phone,omitempty"`
	TokenVerify *string    `db:"token_verify" json:"token_verify,omitempty"`
	IsActive    bool       `db:"is_active" json:"is_active"`
	IsBlocked   bool       `db:"is_blocked" json:"is_blocked"`
	RoleID      uuid.UUID  `db:"role_id" json:"role_id"`
	UploadID    *uuid.UUID `db:"upload_id" json:"upload_id"`
	// Relation
	Role   *Role   `json:"role,omitempty"`
	Upload *Upload `json:"upload,omitempty"`
}

type SignInResponse struct {
	UID         uuid.UUID `json:"uid"`
	Fullname    string    `json:"fullname"`
	Email       string    `json:"email"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsAdmin     bool      `json:"is_admin"`
}

func (entity *User) BeforeCreate() (err error) {
	if entity.Password != nil {
		password, err := lib.Hash(*entity.Password)
		if err != nil {
			return err
		}

		entity.Password = &password
	}

	return
}
