package lib

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"gofi/internal/dtos"
	"gofi/internal/lib/apperr"
	"gofi/internal/lib/validator"

	"github.com/gofiber/fiber/v3"
)

func validationErr(t *testing.T, err error) *ErrValidationFailed {
	t.Helper()
	var ve *ErrValidationFailed
	if !errors.As(err, &ve) {
		t.Fatalf("expected ErrValidationFailed, got %v", err)
	}
	return ve
}

func runBody(t *testing.T, body string, obj Validatable) error {
	t.Helper()
	app := fiber.New()
	var out error
	app.Post("/t", func(c fiber.Ctx) error {
		out = ValidateRequestBody(c, obj)
		if out != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"err": out.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("POST", "/t", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()
	return out
}

func runQuery(t *testing.T, target string, obj Validatable) error {
	t.Helper()
	app := fiber.New()
	var out error
	app.Get("/t", func(c fiber.Ctx) error {
		out = ValidateRequestQuery(c, obj)
		if out != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"err": out.Error()})
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", target, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()
	return out
}

func TestValidateRequestQuery_OptionalDefaultsPass(t *testing.T) {
	q := &dtos.ListQuery{}
	if err := runQuery(t, "/t", q); err != nil {
		t.Fatalf("empty query should pass, got %v", err)
	}
}

func TestValidateRequestQuery_InvalidOffsetFails(t *testing.T) {
	err := runQuery(t, "/t?offset=sideways", &dtos.ListQuery{})
	validationErr(t, err)
}

func TestValidateRequestQuery_ValidParamsPass(t *testing.T) {
	q := &dtos.ListQuery{}
	if err := runQuery(t, "/t?offset=20&limit=10", q); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Offset != 20 || q.Limit != 10 {
		t.Fatalf("query not bound: %+v", q)
	}
}

func TestValidateRequestQuery_EmptyParamTreatedAsAbsent(t *testing.T) {
	q := &dtos.ListQuery{}
	if err := runQuery(t, "/t?offset=", q); err != nil {
		t.Fatalf("empty offset should pass like an absent param, got %v", err)
	}
}

func TestValidateRequestBodyRejectsMalformedInput(t *testing.T) {
	for name, body := range map[string]string{
		"broken json":   `{"name":`,
		"not an object": `[1,2,3]`,
		"wrong type":    `{"name":"x","offset":1.5}`,
	} {
		err := runBody(t, body, &intBody{})
		var ae *apperr.Error
		if !errors.As(err, &ae) || ae.Status() != fiber.StatusBadRequest {
			t.Errorf("%s: expected a 400 apperr, got %v", name, err)
		}
	}
}

type intBody struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
}

func (intBody) Validate(*validator.MapValidator) {}
