package lib

import (
	"encoding/json"
	"fmt"

	"gofi/internal/lib/validator"

	"github.com/gofiber/fiber/v2"
)

type ErrValidationFailed struct {
	MessageRecord validator.MessageRecord
}

func (e ErrValidationFailed) Error() string {
	return fmt.Sprintf("%v", e.MessageRecord)
}

type Validatable interface {
	Validate(v *validator.MapValidator)
}

func ValidateStruct(obj Validatable) error {
	var data map[string]interface{}

	jsonData, err := json.Marshal(obj)
	if err != nil {
		data = make(map[string]interface{})
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		data = make(map[string]interface{})
	}

	v := validator.NewMapValidator()
	obj.Validate(v)

	mr, passed := v.Validate(data)
	if !passed {
		return &ErrValidationFailed{MessageRecord: mr}
	}

	return nil

}

func ValidateRequestQuery(c *fiber.Ctx, obj Validatable) error {
	err := c.QueryParser(obj)
	if err != nil {
		return err
	}

	return ValidateStruct(obj)
}

func ValidateRequestBody(c *fiber.Ctx, obj Validatable) error {
	err := c.BodyParser(obj)
	if err != nil {
		return err
	}

	return ValidateStruct(obj)
}

func WrapValidationError(mr validator.MessageRecord) map[string]interface{} {
	return map[string]interface{}{
		"message": "validation failed",
		"errors":  mr,
	}
}
