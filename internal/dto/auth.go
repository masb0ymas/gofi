package dto

import "gofi/internal/lib/validator"

type AuthSignUp struct {
	FirstName string  `json:"first_name" form:"first_name"`
	LastName  *string `json:"last_name" form:"last_name"`
	Email     string  `json:"email" form:"email"`
	Phone     *string `json:"phone" form:"phone"`
	Password  string  `json:"password" form:"password"`
}

func (dto AuthSignUp) Validate(v *validator.MapValidator) {
	v.Field("first_name").Required().String()
	v.Field("last_name").String()
	v.Field("email").Required().String()
	v.Field("phone").String()
	v.Field("password").Required().String()
}

type AuthSignIn struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (dto AuthSignIn) Validate(v *validator.MapValidator) {
	v.Field("email").Required().String()
	v.Field("password").Required().String()
}

type AuthVerifyRegistration struct {
	Token string `json:"token" form:"token"`
}

func (dto AuthVerifyRegistration) Validate(v *validator.MapValidator) {
	v.Field("token").Required().String()
}

type AuthRefreshToken struct {
	Token string `json:"token" form:"token"`
}

func (dto AuthRefreshToken) Validate(v *validator.MapValidator) {
	v.Field("token").Required().String()
}
