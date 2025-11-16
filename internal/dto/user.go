package dto

import (
	"gofi/internal/lib/validator"

	"github.com/google/uuid"
)

type UserPagination struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func (dto UserPagination) Validate(v *validator.MapValidator) {
	v.Field("offset").Required().Num()
	v.Field("limit").Required().Num()
}

type UserCreate struct {
	FirstName string     `json:"first_name" form:"first_name"`
	LastName  *string    `json:"last_name" form:"last_name"`
	Email     string     `json:"email" form:"email"`
	Phone     *string    `json:"phone" form:"phone"`
	Password  *string    `json:"password" form:"password"`
	RoleID    uuid.UUID  `json:"role_id" form:"role_id"`
	UploadID  *uuid.UUID `json:"upload_id" form:"upload_id"`
}

func (dto UserCreate) Validate(v *validator.MapValidator) {
	v.Field("first_name").Required().String()
	v.Field("last_name").String()
	v.Field("email").Required().String()
	v.Field("phone").String()
	v.Field("password").String()
	v.Field("role_id").Required().UUID()
	v.Field("upload_id").UUID()
}

type UserUpdate struct {
	FirstName string     `json:"first_name" form:"first_name"`
	LastName  *string    `json:"last_name" form:"last_name"`
	Email     string     `json:"email" form:"email"`
	Phone     *string    `json:"phone" form:"phone"`
	Password  *string    `json:"password" form:"password"`
	RoleID    uuid.UUID  `json:"role_id" form:"role_id"`
	UploadID  *uuid.UUID `json:"upload_id" form:"upload_id"`
}

func (dto UserUpdate) Validate(v *validator.MapValidator) {
	v.Field("first_name").String()
	v.Field("last_name").String()
	v.Field("email").String()
	v.Field("phone").String()
	v.Field("password").String()
	v.Field("role_id").UUID()
	v.Field("upload_id").UUID()
}
