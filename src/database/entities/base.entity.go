package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type BaseEntity struct {
	Id        uuid.UUID   `json:"id" db:"id"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt pq.NullTime `json:"deletedAt" db:"deleted_at"`
}
