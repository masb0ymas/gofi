package models

import "time"

type Upload struct {
	Base
	KeyFile   string    `db:"key_file" json:"key_file"`
	FileName  string    `db:"file_name" json:"file_name"`
	MimeType  string    `db:"mimetype" json:"mimetype"`
	Size      int64     `db:"size" json:"size"`
	SignedURL string    `db:"signed_url" json:"signed_url"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}
