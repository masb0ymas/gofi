package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
