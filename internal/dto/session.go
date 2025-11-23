package dto

import "gofi/internal/lib/validator"

type SessionPagination struct {
	Offset int64 `json:"offset" form:"offset"`
	Limit  int64 `json:"limit" form:"limit"`
}

func (dto SessionPagination) Validate(v *validator.MapValidator) {
	v.Field("offset").Required().Num()
	v.Field("limit").Required().Num()
}
