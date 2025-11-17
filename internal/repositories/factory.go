package repositories

import "database/sql"

type Repositories struct {
	Role    RoleRepository
	User    UserRepository
	Session SessionRepository
}

func New(db *sql.DB) Repositories {
	return Repositories{
		Role:    RoleRepository{DB: db},
		User:    UserRepository{DB: db},
		Session: SessionRepository{DB: db},
	}
}
