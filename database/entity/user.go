package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
	Fullname    string     `json:"fullname" db:"fullname"`
	Email       string     `json:"email" db:"email"`
	Password    string     `json:"password" db:"password"`
	Phone       string     `json:"phone" db:"phone"`
	TokenVerify string     `json:"token_verify" db:"token_verify"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	IsBlocked   bool       `json:"is_blocked" db:"is_blocked"`
	RoleID      uuid.UUID  `json:"role_id" db:"role_id"`
	UploadID    *uuid.UUID `json:"upload_id" db:"upload_id"`
}

type UserReq struct {
	Fullname    string     `json:"fullname"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	Phone       string     `json:"phone"`
	TokenVerify string     `json:"token_verify"`
	IsActive    bool       `json:"is_active"`
	IsBlocked   bool       `json:"is_blocked"`
	RoleID      uuid.UUID  `json:"role_id"`
	UploadID    *uuid.UUID `json:"upload_id"`
}

type UserRes struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	Fullname    string     `json:"fullname"`
	Email       string     `json:"email"`
	Password    string     `json:"password"`
	Phone       string     `json:"phone"`
	TokenVerify string     `json:"token_verify"`
	IsActive    bool       `json:"is_active"`
	IsBlocked   bool       `json:"is_blocked"`
	RoleID      uuid.UUID  `json:"role_id"`
	UploadID    *uuid.UUID `json:"upload_id"`
}
