package repository

import (
	"gofi/database/model"

	"github.com/jmoiron/sqlx"
)

type UploadRepository struct {
	*Repository[model.Upload]
}

func NewUploadRepository(db *sqlx.DB) *UploadRepository {
	return &UploadRepository{
		Repository: NewRepository[model.Upload](db, "upload"),
	}
}

func (r *UploadRepository) Create(values *model.Upload) error {
	query := `
		INSERT INTO "upload" (
			"id", "created_at", "updated_at", "key_file", "filename", "mimetype", "size", "signed_url", "expires_at"
		) VALUES (
			:id, :created_at, :updated_at, :key_file, :filename, :mimetype, :size, :signed_url, :expires_at
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"created_at": values.CreatedAt,
		"updated_at": values.UpdatedAt,
		"key_file":   values.KeyFile,
		"filename":   values.Filename,
		"mimetype":   values.Mimetype,
		"size":       values.Size,
		"signed_url": values.SignedURL,
		"expires_at": values.ExpiresAt,
	})
	return err
}

func (r *UploadRepository) Update(values *model.Upload) error {
	query := `
		UPDATE "upload" SET
			"updated_at" = :updated_at,
			"key_file" = :key_file,
			"filename" = :filename,
			"mimetype" = :mimetype,
			"size" = :size,
			"signed_url" = :signed_url,
			"expires_at" = :expires_at,
		WHERE "id" = :id AND "deleted_at" IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"updated_at": values.UpdatedAt,
		"key_file":   values.KeyFile,
		"filename":   values.Filename,
		"mimetype":   values.Mimetype,
		"size":       values.Size,
		"signed_url": values.SignedURL,
		"expires_at": values.ExpiresAt,
	})
	return err
}
