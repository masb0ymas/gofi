package dto

import "gofi/internal/lib/validator"

type UserOAuthPagination struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func (dto UserOAuthPagination) Validate(v *validator.MapValidator) {
	v.Field("offset").Required().Num()
	v.Field("limit").Required().Num()
}

type UserOAuthCreate struct {
	Provider     string  `json:"provider" form:"provider"`
	AccessToken  string  `json:"access_token" form:"access_token"`
	RefreshToken *string `json:"refresh_token" form:"refresh_token"`
	ExpiresAt    string  `json:"expires_at" form:"expires_at"`
}

func (dto UserOAuthCreate) Validate(v *validator.MapValidator) {
	v.Field("provider").Required().String()
	v.Field("access_token").Required().String()
	v.Field("refresh_token").String()

	if dto.ExpiresAt != "" {
		v.Field("expires_at").Required().String()
	}
}

type UserOAuthUpdate struct {
	Provider     *string `json:"provider" form:"provider"`
	AccessToken  *string `json:"access_token" form:"access_token"`
	RefreshToken *string `json:"refresh_token" form:"refresh_token"`
	ExpiresAt    *string `json:"expires_at" form:"expires_at"`
}

func (dto UserOAuthUpdate) Validate(v *validator.MapValidator) {
	v.Field("provider").String()
	v.Field("access_token").String()
	v.Field("refresh_token").String()
	v.Field("expires_at").String()
}
