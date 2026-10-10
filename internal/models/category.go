package models

type Category struct {
	Base
	SoftDelete
	Name string `db:"name" json:"name"`
}
