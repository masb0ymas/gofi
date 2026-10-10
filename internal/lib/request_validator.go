package lib

import (
	"encoding/json"
	"fmt"

	"gofi/internal/lib/apperr"
	"gofi/internal/lib/validator"

	"github.com/gofiber/fiber/v3"
)

type ErrValidationFailed struct {
	MessageRecord validator.MessageRecord
}

func (e ErrValidationFailed) Error() string {
	return fmt.Sprintf("%v", e.MessageRecord)
}

var errInvalidBody = apperr.New(apperr.KindBadRequest, "invalid JSON body")

type Validatable interface {
	Validate(v *validator.MapValidator)
}

// ValidateRequestQuery validates the raw query map before binding it into
// obj, so optional parameters that are absent stay absent.
func ValidateRequestQuery(c fiber.Ctx, obj Validatable) error {
	data := make(map[string]interface{}, len(c.Queries()))
	for key, val := range c.Queries() {
		data[key] = val
	}
	dropEmpty(data)

	if err := validateDict(obj, data); err != nil {
		return err
	}

	if err := c.Bind().Query(obj); err != nil {
		return apperr.New(apperr.KindBadRequest, "invalid query string")
	}
	return nil
}

// ValidateRequestBody validates the raw JSON object before binding it into
// obj, so optional fields that are absent stay absent (an empty struct field
// would otherwise fail rules like Email/WithinS). An empty body is treated
// as an empty object.
func ValidateRequestBody(c fiber.Ctx, obj Validatable) error {
	body := c.Body()

	data := make(map[string]interface{})
	if len(body) > 0 {
		if err := json.Unmarshal(body, &data); err != nil {
			return errInvalidBody
		}
	}
	dropEmpty(data)

	if err := validateDict(obj, data); err != nil {
		return err
	}

	if len(body) == 0 {
		return nil
	}
	// The rules passed but a value still does not fit its field (e.g. 1.5
	// into an integer): that is the client's mistake, not a server error.
	if err := c.Bind().Body(obj); err != nil {
		return errInvalidBody
	}
	return nil
}

// dropEmpty removes empty-string values so an explicit "" is validated like
// an absent key: optional format/enum rules skip it and Required still
// rejects it. The struct binding still sees the original value.
func dropEmpty(data map[string]interface{}) {
	for key, val := range data {
		if s, ok := val.(string); ok && s == "" {
			delete(data, key)
		}
	}
}

func validateDict(obj Validatable, data map[string]interface{}) error {
	v := validator.NewMapValidator()
	obj.Validate(v)

	mr, passed := v.Validate(data)
	if !passed {
		return &ErrValidationFailed{MessageRecord: mr}
	}

	return nil
}

// FieldError mirrors the validation error item shape the web UI expects:
// an array of {code, field, message, param} in the response body.
type FieldError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
	Param   string `json:"param"`
}

// WrapValidationError renders a MessageRecord into the validation error
// envelope: {"message": "validation failed", "errors": [FieldError...]}.
func WrapValidationError(mr validator.MessageRecord) map[string]interface{} {
	return map[string]interface{}{
		"message": "validation failed",
		"errors":  FieldErrorsFromRecord(mr),
	}
}

// FieldErrorsFromRecord flattens the field → messages record into the array
// shape the web UI consumes.
func FieldErrorsFromRecord(mr validator.MessageRecord) []FieldError {
	out := make([]FieldError, 0, len(mr))
	for field, messages := range mr {
		for _, message := range messages {
			out = append(out, FieldError{
				Code:    "validation",
				Field:   field,
				Message: message,
				Param:   field,
			})
		}
	}
	return out
}
