package models

import (
	"time"

	"github.com/google/uuid"
)

type Base struct {
	ID        uuid.UUID `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// SoftDelete marks tables whose rows are hidden by setting deleted_at instead
// of being removed.
type SoftDelete struct {
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
