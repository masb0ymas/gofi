package model

import (
	"time"
)

type Upload struct {
	BaseModel
	KeyFile   string    `db:"keyfile" json:"keyfile"`
	Filename  string    `db:"filename" json:"filename"`
	Mimetype  string    `db:"mimetype" json:"mimetype"`
	Size      int64     `db:"size" json:"size"`
	SignedURL string    `db:"signed_url" json:"signed_url"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
