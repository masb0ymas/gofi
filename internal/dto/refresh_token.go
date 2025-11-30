package dto

import (
	"gofi/internal/lib/validator"
	"time"
)

type RefreshTokenPagination struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func (dto RefreshTokenPagination) Validate(v *validator.MapValidator) {
	v.Field("offset").Required().Num()
	v.Field("limit").Required().Num()
}

type RefreshTokenCreate struct {
	UserID    string    `json:"user_id" form:"user_id"`
	Token     string    `json:"token" form:"token"`
	ExpiresAt time.Time `json:"expires_at" form:"expires_at"`
}

func (dto RefreshTokenCreate) Validate(v *validator.MapValidator) {
	v.Field("user_id").Required().String()
	v.Field("token").Required().String()
	v.Field("expires_at").Required().Date()
}

type RefreshTokenUpdate struct {
	Token     *string `json:"token" form:"token"`
	ExpiresAt *string `json:"expires_at" form:"expires_at"`
}

func (dto RefreshTokenUpdate) Validate(v *validator.MapValidator) {
	v.Field("token").String()
	v.Field("expires_at").String()
}
