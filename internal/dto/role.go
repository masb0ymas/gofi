package dto

import "gofi/internal/lib/validator"

type RolePagination struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func (dto RolePagination) Validate(v *validator.MapValidator) {
	v.Field("offset").Required().Num()
	v.Field("limit").Required().Num()
}

type RoleCreate struct {
	Name string `json:"name" form:"name"`
}

func (dto RoleCreate) Validate(v *validator.MapValidator) {
	v.Field("name").Required().String()
}

type RoleUpdate struct {
	Name string `json:"name" form:"name"`
}

func (dto RoleUpdate) Validate(v *validator.MapValidator) {
	v.Field("name").Required().String()
}
