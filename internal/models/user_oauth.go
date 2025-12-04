package models

import (
	"time"

	"github.com/google/uuid"
)

type UserOAuth struct {
	ID                 uuid.UUID `db:"id" json:"id"`
	UserID             uuid.UUID `db:"user_id" json:"user_id"`
	IdentityProviderID string    `db:"identity_provider_id" json:"identity_provider_id"`
	Provider           string    `db:"provider" json:"provider"`
	AccessToken        string    `db:"access_token" json:"access_token"`
	RefreshToken       *string   `db:"refresh_token" json:"refresh_token,omitempty"`
	ExpiresAt          time.Time `db:"expires_at" json:"expires_at"`
}
