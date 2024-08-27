package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
)

var validate = validator.New()

// validate struct
func validateStruct(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// validate
func Validate(data interface{}) (code int32, message string, errors []string) {
	if errs := validateStruct(data); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return http.StatusBadRequest, fiber.ErrBadRequest.Message, errMsgs
	}

	return 400, "", nil
}

// parse form data body
func ParseFormData(c *fiber.Ctx, body interface{}) (code int32, message string, errors []string) {
	if err := c.BodyParser(body); err != nil {
		errMsgs := make([]string, 0)
		errMsgs = append(errMsgs, fiber.ErrUnprocessableEntity.Message)

		return http.StatusUnprocessableEntity, fiber.ErrUnprocessableEntity.Message, errMsgs
	}

	return 400, "", nil
}

// parse form data body and validate form
func ParseFormDataAndValidate(c *fiber.Ctx, body interface{}) (code int32, message string, errors []string) {
	if code, message, errors := ParseFormData(c, body); errors != nil {
		return code, message, errors
	}

	return Validate(body)
}
