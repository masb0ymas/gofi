package repositories

import "database/sql"

type Repositories struct {
	Role    RoleRepository
	User    UserRepository
	Session SessionRepository
}

func New(db *sql.DB) Repositories {
	return Repositories{
		Role:    RoleRepository{BaseRepository: BaseRepository{DB: db, TableName: "roles"}},
		User:    UserRepository{BaseRepository: BaseRepository{DB: db, TableName: "users"}},
		Session: SessionRepository{DB: db},
	}
}
