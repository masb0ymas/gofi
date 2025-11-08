package models

type Role struct {
	Base
	Name string `db:"name" json:"name"`
}
