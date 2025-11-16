package models

import "github.com/google/uuid"

type Session struct {
	Base
	UserID uuid.UUID `db:"user_id" json:"user_id"`
	Token  string    `db:"token" json:"token"`
}
