package models

// User mirrors the Authula user of the same id (see migration 000002) and
// adds the application role.
type User struct {
	Base
	Name  string  `db:"name" json:"name"`
	Email string  `db:"email" json:"email"`
	Image *string `db:"image" json:"image"`
	Role  string  `db:"role" json:"role"`
}
