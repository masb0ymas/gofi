package repository

import (
	"errors"
	"fmt"
	"goarif-api/config"
	"goarif-api/database/model"
	"goarif-api/lib"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	*Repository[model.User]
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		Repository: NewRepository[model.User](db, "user"),
	}
}

func (r *AuthRepository) SignUp(user *model.User) error {
	var existingUser model.User

	query := `
		SELECT * FROM "user" WHERE email = $1
	`
	_ = r.db.Get(&existingUser, query, user.Email)

	if existingUser.ID != uuid.Nil {
		return errors.New("user already exists")
	}

	token, _, err := lib.GenerateToken(&lib.Payload{
		UID:       user.ID,
		SecretKey: config.Env("JWT_SECRET_KEY", "secret"),
		ExpiresAt: config.Env("JWT_EXPIRES_IN", "30"), // days
	})
	if err != nil {
		return err
	}

	password, err := lib.Hash(*user.Password)
	if err != nil {
		return err
	}

	fmt.Println(token, password)

	query = `
		INSERT INTO "user" (
			id, created_at, updated_at, fullname, email, password, token_verify, role_id
		) VALUES (
			:id, :created_at, :updated_at, :fullname, :email, :password, :token_verify, :role_id
		)
	`
	_, err = r.db.NamedExec(query, map[string]interface{}{
		"id":           uuid.New(),
		"created_at":   time.Now(),
		"updated_at":   time.Now(),
		"fullname":     user.Fullname,
		"email":        user.Email,
		"password":     password,
		"token_verify": token,
		"role_id":      user.RoleID,
	})

	return err
}

func (r *AuthRepository) SignIn(user *model.User, ip_address string, user_agent string) (*model.SignInResponse, error) {
	var existingUser model.User

	query := `
		SELECT * FROM "user" WHERE email = $1
	`
	err := r.db.Get(&existingUser, query, user.Email)
	if err != nil {
		return nil, err
	}

	if existingUser.ID == uuid.Nil {
		return nil, errors.New("user not found")
	}

	if !existingUser.IsActive {
		return nil, errors.New("user is not active")
	}

	matchPassword, err := lib.VerifyHash(*user.Password, *existingUser.Password)
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
		INSERT INTO "session" (
			id, created_at, updated_at, user_id, token, expires_at, ip_address, user_agent
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

	return &model.SignInResponse{UID: existingUser.ID, Token: token, ExpiresAt: time.Unix(expiresAtUnix, 0)}, nil
}

func (r *AuthRepository) VerifySession(user_id uuid.UUID, token string) (*model.User, error) {
	var existingSession model.Session
	var existingUser model.User

	query := `
		SELECT * FROM "session" WHERE user_id = $1 AND token = $2
	`
	err := r.db.Get(&existingSession, query, user_id, token)
	if err != nil {
		return nil, err
	}

	if existingSession.ID == uuid.Nil {
		return nil, errors.New("session not found")
	}

	if existingSession.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session is expired")
	}

	query = `
		SELECT * FROM "user" WHERE id = $1
	`
	err = r.db.Get(&existingUser, query, existingSession.UserID)
	if err != nil {
		return nil, err
	}

	if existingUser.ID == uuid.Nil {
		return nil, errors.New("user not found")
	}

	return &existingUser, nil
}

func (r *AuthRepository) SignOut(user_id uuid.UUID, token string) error {
	var existingUser model.User

	query := "SELECT * FROM user WHERE id = $1"
	err := r.db.Get(&existingUser, query, user_id)
	if err != nil {
		return err
	}

	if existingUser.ID == uuid.Nil {
		return errors.New("user not found")
	}

	query = "DELETE FROM session WHERE user_id = $1 AND token = $2"
	_, err = r.db.Exec(query, user_id, token)
	return err
}
