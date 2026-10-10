package repositories

import "database/sql"

type Repositories struct {
	Category CategoryRepository
}

// New builds every repository over db; debug prints each query (DEBUG=true).
func New(db *sql.DB, debug bool) Repositories {
	base := BaseRepository{DB: db, Debug: debug}

	return Repositories{
		Category: CategoryRepository{base},
	}
}
