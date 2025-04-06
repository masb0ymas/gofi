package model

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	BaseModel
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	IpAddress string    `db:"ip_address" json:"ip_address"`
	UserAgent string    `db:"user_agent" json:"user_agent"`
	Latitude  *string   `db:"latitude" json:"latitude"`
	Longitude *string   `db:"longitude" json:"longitude"`
	// Relation
	User *User `json:"user,omitempty"`
}
