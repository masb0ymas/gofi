package repositories

import "database/sql"

type Repositories struct {
	Role    RoleRepository
	User    UserRepository
	Session SessionRepository
}

func New(DB *sql.DB) Repositories {
	return Repositories{
		Role:    RoleRepository{DB: DB},
		User:    UserRepository{DB: DB},
		Session: SessionRepository{DB: DB},
	}
}
