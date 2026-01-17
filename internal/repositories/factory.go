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
		UserVerifyAccount: UserVerifyAccountRepository{DB: db, Config: config},
		Session:           SessionRepository{DB: db, Config: config},
		RefreshToken:      RefreshTokenRepository{DB: db, Config: config},
		UserOAuth:         UserOAuthRepository{DB: db, Config: config},
	}
}
