package entities

type RoleEntity struct {
	BaseEntity
	Name string `json:"name" db:"name"`
}
