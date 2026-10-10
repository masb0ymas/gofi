package dtos

import "gofi/internal/lib/validator"

// maxNameLen bounds a category name; the column itself is unbounded.
const maxNameLen = 100

// CategoryBody is the create/update payload for a category.
type CategoryBody struct {
	Name string `json:"name"`
}

func (dto CategoryBody) Validate(v *validator.MapValidator) {
	v.Field("name").Required().String().MaxLen(maxNameLen)
}
