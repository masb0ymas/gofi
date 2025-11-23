package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Base
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token,omitempty"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IPAddress string    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
}
