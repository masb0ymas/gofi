package entities

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type UserEntity struct {
	BaseEntity
	Fullname    string      `json:"fullname" db:"fullname"`
	Email       string      `json:"email" db:"email"`
	Password    string      `json:"password" db:"password"`
	Phone       null.String `json:"phone" db:"phone"`
	TokenVerify null.String `json:"tokenVerify" db:"token_verify"`
	Address     null.String `json:"address" db:"address"`
	IsActive    bool        `json:"isActive" db:"is_active"`
	IsBlocked   bool        `json:"isBlocked" db:"is_blocked"`
	RoleId      uuid.UUID   `json:"RoleId" db:"role_id"`
}
