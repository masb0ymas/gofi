package repository

import (
	"database/sql"
	"gofi/src/app/database/model"
	"gofi/src/lib"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	*Repository[model.User]
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		Repository: NewRepository[model.User](db, "user"),
	}
}

func (r *UserRepository) FindAllWithPagination(req *lib.Pagination) ([]model.User, int64, error) {
	var records []model.User
	var total int64

	query := `
		SELECT 
			"u"."id", "u"."created_at", "u"."updated_at", "u"."deleted_at", "u"."fullname", "u"."email", "u"."phone", "u"."is_active", "u"."is_blocked", "u"."role_id", "u"."upload_id",
			"r".*,
			"up".*
		FROM "user" "u"
		LEFT JOIN "role" "r" ON "u"."role_id" = "r"."id"
		LEFT JOIN "upload" "up" ON "u"."upload_id" = "up"."id"
		WHERE "u"."deleted_at" IS NULL
	`
	query = lib.QueryBuilder(query, req.Filtered, req.Sorted, req.Page, req.PageSize)
	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		var role model.Role

		// Use nullable types for Upload fields
		var uploadID sql.NullString
		var uploadCreatedAt sql.NullTime
		var uploadUpdatedAt sql.NullTime
		var uploadDeletedAt sql.NullTime
		var uploadKeyFile sql.NullString
		var uploadFilename sql.NullString
		var uploadMimetype sql.NullString
		var uploadSize sql.NullInt64
		var uploadSignedURL sql.NullString
		var uploadExpiresAt sql.NullTime

		err = rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Fullname, &user.Email, &user.Phone, &user.IsActive, &user.IsBlocked, &user.RoleID, &user.UploadID,
			&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt, &role.Name,
			&uploadID, &uploadCreatedAt, &uploadUpdatedAt, &uploadDeletedAt, &uploadKeyFile, &uploadFilename, &uploadMimetype, &uploadSize, &uploadSignedURL, &uploadExpiresAt,
		)
		if err != nil {
			return nil, 0, err
		}

		user.Role = &role
		user.Upload = nil
		// Explicitly set these fields to nil to exclude them from JSON
		user.Password = nil
		user.TokenVerify = nil

		// Only set Upload if it exists (not null)
		if uploadID.Valid {
			upload := model.Upload{
				BaseModel: model.BaseModel{
					ID:        uuid.MustParse(uploadID.String),
					CreatedAt: uploadCreatedAt.Time,
					UpdatedAt: uploadUpdatedAt.Time,
				},
				KeyFile:   uploadKeyFile.String,
				Filename:  uploadFilename.String,
				Mimetype:  uploadMimetype.String,
				Size:      uploadSize.Int64,
				SignedURL: uploadSignedURL.String,
				ExpiresAt: uploadExpiresAt.Time,
			}

			if uploadDeletedAt.Valid {
				upload.DeletedAt = &uploadDeletedAt.Time
			}

			user.Upload = &upload
		}

		records = append(records, user)
	}

	query_total := `
		SELECT COUNT(*) FROM "user" WHERE "deleted_at" IS NULL
	`
	query_total = lib.QueryBuilderForCount(query_total, req.Filtered)
	err = r.db.Get(&total, query_total)
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	query := `
		SELECT * FROM "user" WHERE "email" = $1 AND "deleted_at" IS NULL
	`
	err := r.db.Get(&user, query, email)
	return &user, err
}

func (r *UserRepository) Create(values *model.User) error {
	hash, err := lib.Hash(*values.Password)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO "user" (
			"id", "created_at", "updated_at", "fullname", "email", "password", "is_active", "role_id"
		) VALUES (
			:id, :created_at, :updated_at, :fullname, :email, :password, :is_active, :role_id
		)
	`
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"created_at": values.CreatedAt,
		"updated_at": values.UpdatedAt,
		"fullname":   values.Fullname,
		"email":      values.Email,
		"password":   hash,
		"is_active":  values.IsActive,
		"role_id":    values.RoleID,
	})
	return err
}

func (r *UserRepository) Update(values *model.User) error {
	query := `
		UPDATE "user" SET
			"updated_at" = :updated_at,
			"fullname" = :fullname,
			"email" = :email,
		WHERE "id" = :id AND "deleted_at" IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"updated_at": values.UpdatedAt,
		"fullname":   values.Fullname,
		"email":      values.Email,
	})
	return err
}

func (r *UserRepository) FindUserByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	var role model.Role

	// Use nullable types for Upload fields
	var uploadID sql.NullString
	var uploadCreatedAt sql.NullTime
	var uploadUpdatedAt sql.NullTime
	var uploadDeletedAt sql.NullTime
	var uploadKeyFile sql.NullString
	var uploadFilename sql.NullString
	var uploadMimetype sql.NullString
	var uploadSize sql.NullInt64
	var uploadSignedURL sql.NullString
	var uploadExpiresAt sql.NullTime

	query := `
		SELECT 
			"u"."id", "u"."created_at", "u"."updated_at", "u"."deleted_at", "u"."fullname", "u"."email", "u"."phone", "u"."is_active", "u"."is_blocked", "u"."role_id", "u"."upload_id",
			"r".*,
			"up".*
		FROM "user" "u"
		LEFT JOIN "role" "r" ON "u"."role_id" = "r"."id"
		LEFT JOIN "upload" "up" ON "u"."upload_id" = "up"."id"
		WHERE "u"."id" = $1 AND "u"."deleted_at" IS NULL
	`
	rows, err := r.db.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.Fullname, &user.Email, &user.Phone,
			&user.IsActive, &user.IsBlocked, &user.RoleID, &user.UploadID,
			&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt, &role.Name,
			&uploadID, &uploadCreatedAt, &uploadUpdatedAt, &uploadDeletedAt, &uploadKeyFile, &uploadFilename, &uploadMimetype, &uploadSize, &uploadSignedURL, &uploadExpiresAt,
		)

		if err != nil {
			return nil, err
		}

		user.Role = &role
		user.Upload = nil
		// Explicitly set these fields to nil to exclude them from JSON
		user.Password = nil
		user.TokenVerify = nil

		// Only set Upload if it exists (not null)
		if uploadID.Valid {
			upload := model.Upload{
				BaseModel: model.BaseModel{
					ID:        uuid.MustParse(uploadID.String),
					CreatedAt: uploadCreatedAt.Time,
					UpdatedAt: uploadUpdatedAt.Time,
				},
				KeyFile:   uploadKeyFile.String,
				Filename:  uploadFilename.String,
				Mimetype:  uploadMimetype.String,
				Size:      uploadSize.Int64,
				SignedURL: uploadSignedURL.String,
				ExpiresAt: uploadExpiresAt.Time,
			}

			if uploadDeletedAt.Valid {
				upload.DeletedAt = &uploadDeletedAt.Time
			}

			user.Upload = &upload
		}
	}

	return &user, nil
}
