package entities

import (
	"gofi/src/pkg/helpers"
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	Id        uuid.UUID        `json:"id" db:"id"`
	CreatedAt time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time        `json:"updatedAt" db:"updated_at"`
	DeletedAt helpers.NullTime `json:"deletedAt" db:"deleted_at"`
}
