package entity

import (
	"time"

	"github.com/google/uuid"
)

type Upload struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
	KeyFile   string     `json:"key_file" db:"key_file"`
	Filename  string     `json:"filename" db:"filename"`
	Mimetype  string     `json:"mimetype" db:"mimetype"`
	Size      int32      `json:"size" db:"size"`
	SignedUrl string     `json:"signed_url" db:"signed_url"`
	ExpiredAt time.Time  `json:"expired_at" db:"expired_at"`
}

type UploadReq struct {
	KeyFile   string    `json:"key_file"`
	Filename  string    `json:"filename"`
	Mimetype  string    `json:"mimetype"`
	Size      int32     `json:"size"`
	SignedUrl string    `json:"signed_url"`
	ExpiredAt time.Time `json:"expired_at"`
}

type UploadRes struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	KeyFile   string     `json:"key_file"`
	Filename  string     `json:"filename"`
	Mimetype  string     `json:"mimetype"`
	Size      int32      `json:"size"`
	SignedUrl string     `json:"signed_url"`
	ExpiredAt time.Time  `json:"expired_at"`
}
