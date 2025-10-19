package repository

import (
	"database/sql"
	"errors"
	"gofi/src/app/database/model"
	"gofi/src/config"
	"gofi/src/lib"
	"gofi/src/lib/constant"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/masb0ymas/go-utils/pkg"
)

type AuthRepository struct {
	*Repository[model.User]
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		Repository: NewRepository[model.User](db, "user"),
	}
}

func (r *AuthRepository) SignUp(values *model.User) (*string, error) {
	var existingUser model.User

	token, _, err := lib.GenerateToken(&lib.Payload{
		UID:       uuid.New(),
		SecretKey: config.Env("JWT_SECRET_KEY", "secret"),
		ExpiresAt: config.Env("JWT_EXPIRES_IN", "30"), // days
	})
	if err != nil {
		return nil, err
	}

	query := `
		SELECT * FROM "user" WHERE "email" = $1
	`
	err = r.db.Get(&existingUser, query, values.Email)
	// skip error no rows ( user not found )
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	if existingUser.ID != uuid.Nil {
		return nil, errors.New("user already exists")
	}

	if err := values.BeforeCreate(); err != nil {
		return nil, err
	}

	query = `
		INSERT INTO "user" (
			"id", "created_at", "updated_at", "fullname", "email", "password", "token_verify", "role_id"
		) VALUES (
			:id, :created_at, :updated_at, :fullname, :email, :password, :token_verify, :role_id
		)
	`
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"id":           values.ID,
		"created_at":   values.CreatedAt,
		"updated_at":   values.UpdatedAt,
		"fullname":     values.Fullname,
		"email":        values.Email,
		"password":     values.Password,
		"token_verify": token,
		"role_id":      values.RoleID,
	})

	return &token, err
}

func (r *AuthRepository) SignIn(values *model.User, ip_address string, user_agent string) (*model.SignInResponse, error) {
	var (
		existingUser    model.User
		existingSession model.Session
	)

	query := `
		SELECT * FROM "user" WHERE "email" = $1
	`
	err := r.db.Get(&existingUser, query, values.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user is not registered")
		}

		return nil, err
	}

	if existingUser.ID == uuid.Nil {
		return nil, errors.New("user not found")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user is not active")
	}

	matchPassword, err := lib.VerifyHash(*values.Password, *existingUser.Password)
	if err != nil {
		return nil, err
	}

	if !matchPassword {
		return nil, errors.New("password is incorrect")
	}

	token, expiresAtUnix, err := lib.GenerateToken(&lib.Payload{
		UID:       existingUser.ID,
		SecretKey: config.Env("JWT_SECRET_KEY", "secret"),
		ExpiresAt: config.Env("JWT_EXPIRES_IN", "30"), // days
	})
	if err != nil {
		return nil, err
	}

	query = `
		SELECT * FROM "session" WHERE "user_id" = $1
	`
	err = r.db.Get(&existingSession, query, existingUser.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	query = `
		INSERT INTO "session" (
			"id", "created_at", "updated_at", "user_id", "token", "expires_at", "ip_address", "user_agent"
		) VALUES (
			:id, :created_at, :updated_at, :user_id, :token, :expires_at, :ip_address, :user_agent
		)
	`
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"id":         uuid.New(),
		"created_at": time.Now(),
		"updated_at": time.Now(),
		"user_id":    existingUser.ID,
		"token":      token,
		"expires_at": time.Unix(expiresAtUnix, 0),
		"ip_address": ip_address,
		"user_agent": user_agent,
	})
	if err != nil {
		return nil, err
	}

	adminRole := []string{constant.ID_SUPER_ADMIN, constant.ID_ADMIN}

	return &model.SignInResponse{
		Fullname:    existingUser.Fullname,
		Email:       existingUser.Email,
		UID:         existingUser.ID,
		AccessToken: token,
		ExpiresAt:   time.Unix(expiresAtUnix, 0),
		IsAdmin:     pkg.Contains(adminRole, existingUser.RoleID.String()),
	}, nil
}

func (r *AuthRepository) VerifySession(user_id uuid.UUID, token string) (*model.User, error) {
	var existingSession model.Session

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
		SELECT * FROM "session" WHERE "user_id" = $1 AND "token" = $2
	`
	err := r.db.Get(&existingSession, query, user_id, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("session not found")
		}

		return nil, err
	}

	if existingSession.ID == uuid.Nil {
		return nil, errors.New("session not found")
	}

	if existingSession.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session is expired")
	}

	query = `
		SELECT 
			"u"."id", "u"."created_at", "u"."updated_at", "u"."deleted_at", "u"."fullname", "u"."email", "u"."phone", "u"."is_active", "u"."is_blocked", "u"."role_id", "u"."upload_id",
			"r".*,
			"up".*
		FROM "user" "u"
		LEFT JOIN "role" "r" ON "u"."role_id" = "r"."id"
		LEFT JOIN "upload" "up" ON "u"."upload_id" = "up"."id"
		WHERE "u"."id" = $1 AND "u"."deleted_at" IS NULL
	`
	rows, err := r.db.Queryx(query, user_id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

func (r *AuthRepository) VerifyToken(uid uuid.UUID, token string) error {
	var existingUser model.User

	query := `
		SELECT * FROM "user" WHERE "id" = $1 AND "token_verify" = $2
	`
	err := r.db.Get(&existingUser, query, uid, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found or verification token is invalid")
		}

		return err
	}

	if existingUser.ID == uuid.Nil {
		return errors.New("user not found")
	}

	if existingUser.ID != uuid.Nil {
		query = `
		UPDATE "user" SET
			"updated_at" = :updated_at, 
			"token_verify" = :token_verify,
			"is_active" = :is_active
		WHERE "id" = :id AND "deleted_at" IS NULL
	`
		_, err = r.db.NamedExec(query, map[string]interface{}{
			"updated_at":   time.Now(),
			"token_verify": token,
			"is_active":    true,
			"id":           uid,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *AuthRepository) SignOut(user_id uuid.UUID, token string) error {
	var existingUser model.User

	query := `
		SELECT * FROM "user" WHERE "id" = $1
	`
	err := r.db.Get(&existingUser, query, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}

		return err
	}

	if existingUser.ID == uuid.Nil {
		return errors.New("user not found")
	}

	query = `
		DELETE FROM "session" WHERE "user_id" = $1 AND "token" = $2
	`
	_, err = r.db.Exec(query, user_id, token)
	return err
}
