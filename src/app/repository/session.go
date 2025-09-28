package repository

import (
	"gofi/src/app/database/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	*Repository[model.Session]
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		Repository: NewRepository[model.Session](db, "session"),
	}
}

func (r *SessionRepository) FindByUserId(id uuid.UUID) (*model.Session, error) {
	var session model.Session
	query := `
		SELECT * FROM "session" WHERE "user_id" = $1 AND "deleted_at" IS NULL
	`
	err := r.db.Get(&session, query, id)
	return &session, err
}

func (r *SessionRepository) Create(values *model.Session) error {
	query := `
		INSERT INTO "session" (
			"id", "created_at", "updated_at", "user_id", "token", "expires_at", "ip_address", "user_agent", "latitude", "longitude"
		) VALUES (
			:id, :created_at, :updated_at, :user_id, :token, :expires_at, :ip_address, :user_agent, :latitude, :longitude
		)
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"created_at": values.CreatedAt,
		"updated_at": values.UpdatedAt,
		"user_id":    values.UserID,
		"token":      values.Token,
		"expires_at": values.ExpiresAt,
		"ip_address": values.IpAddress,
		"user_agent": values.UserAgent,
		"latitude":   values.Latitude,
		"longitude":  values.Longitude,
	})
	return err
}

func (r *SessionRepository) Update(values *model.Session) error {
	query := `
		UPDATE "session" SET
			"updated_at" = :updated_at,
			"user_id" = :user_id,
			"token" = :token,
			"expires_at" = :expires_at,
			"ip_address" = :ip_address,
			"user_agent" = :user_agent,
			"latitude" = :latitude,
			"longitude" = :longitude
		WHERE "id" = :id AND "deleted_at" IS NULL
	`
	_, err := r.db.NamedExec(query, map[string]interface{}{
		"id":         values.ID,
		"updated_at": values.UpdatedAt,
		"user_id":    values.UserID,
		"token":      values.Token,
		"expires_at": values.ExpiresAt,
		"ip_address": values.IpAddress,
		"user_agent": values.UserAgent,
		"latitude":   values.Latitude,
		"longitude":  values.Longitude,
	})
	return err
}
