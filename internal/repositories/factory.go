package repositories

import "database/sql"

type Repositories struct {
	Category CategoryRepository
	User     UserRepository
}

// New builds every repository over db.
func New(db *sql.DB) Repositories {
	base := BaseRepository{DB: db}

	return Repositories{
		Category: CategoryRepository{base},
		User:     UserRepository{base},
	}
}
