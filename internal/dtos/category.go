package dtos

import "gofi/internal/lib/validator"

// CategoryBody is the create/update payload for a category.
type CategoryBody struct {
	Name string `json:"name"`
}

func (dto CategoryBody) Validate(v *validator.MapValidator) {
	v.Field("name").Required().String()
}
