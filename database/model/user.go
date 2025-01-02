package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Fullname        string            `db:"fullname" json:"fullname"`
	Email           string            `db:"email" json:"email"`
	Password        *string           `db:"password" json:"password"`
	Phone           *string           `db:"phone" json:"phone"`
	TokenVerify     *string           `db:"token_verify" json:"token_verify"`
	IsActive        bool              `db:"is_active" json:"is_active"`
	IsBlocked       bool              `db:"is_blocked" json:"is_blocked"`
	RoleID          uuid.UUID         `db:"role_id" json:"role_id"`
	UploadID        *uuid.UUID        `db:"upload_id" json:"upload_id"`
	Role            Role              `json:"role,omitempty"`
	Upload          Upload            `json:"upload,omitempty"`
	AccountProvider []AccountProvider `json:"account_provider,omitempty"`
}

type SignInResponse struct {
	UID       uuid.UUID `json:"uid"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
