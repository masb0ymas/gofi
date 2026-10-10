package dtos

import "gofi/internal/lib/validator"

// ListQuery is the shared pagination binding for list endpoints.
// Offset starts at 0; limit is clamped to 100. Sorted is a
// "column:asc|desc" pair (e.g. "name:asc"); repositories resolve the column
// against a whitelist.
type ListQuery struct {
	Offset int    `query:"offset"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	Sorted string `query:"sorted"`
}

// Validate declares rules over the raw query map: every parameter is a
// string and optional, but when present it must be well formed.
func (dto ListQuery) Validate(v *validator.MapValidator) {
	dto.base(v)
}

// base declares the shared pagination/search rules. Query DTOs that embed
// ListQuery call it so their own fields extend the same validator.
func (dto ListQuery) base(v *validator.MapValidator) {
	v.Field("offset").Regex(`^\d+$`)
	v.Field("limit").Regex(`^\d+$`)
	v.Field("search").String()
	v.Field("sorted").Regex(`^[a-z_]+:(asc|desc)$`)
}

const (
	defaultListLimit = 10
	maxListLimit     = 100
)

// Clamp normalizes pagination values: a missing or zero limit falls back to
// 10 and the limit never exceeds 100.
func (dto *ListQuery) Clamp() {
	if dto.Limit <= 0 {
		dto.Limit = defaultListLimit
	}
	if dto.Limit > maxListLimit {
		dto.Limit = maxListLimit
	}
	if dto.Offset < 0 {
		dto.Offset = 0
	}
}
