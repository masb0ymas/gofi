package repositories

import (
	"database/sql"

	"gofi/internal/config"
)

type Repositories struct {
	Role              RoleRepository
	User              UserRepository
	UserVerifyAccount UserVerifyAccountRepository
	Session           SessionRepository
	RefreshToken      RefreshTokenRepository
	UserOAuth         UserOAuthRepository
}

func New(db *sql.DB, config *config.ConfigApp) Repositories {
	return Repositories{
		Role:              RoleRepository{BaseRepository: BaseRepository{DB: db, TableName: "roles", Config: config}},
		User:              UserRepository{BaseRepository: BaseRepository{DB: db, TableName: "users", Config: config}},
		UserVerifyAccount: UserVerifyAccountRepository{BaseRepository: BaseRepository{DB: db, TableName: "user_verify_accounts", Config: config}},
		Session:           SessionRepository{BaseRepository: BaseRepository{DB: db, TableName: "sessions", Config: config}},
		RefreshToken:      RefreshTokenRepository{BaseRepository: BaseRepository{DB: db, TableName: "refresh_tokens", Config: config}},
		UserOAuth:         UserOAuthRepository{BaseRepository: BaseRepository{DB: db, TableName: "user_oauths", Config: config}},
	}
}
