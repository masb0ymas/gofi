package middlewares

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"gofi/internal/lib"
	"gofi/internal/lib/apperr"
)

// ErrorHandler renders handler errors: validation failures as 422 with the
// field list, apperr kinds with their status, and anything unknown as a
// generic 500 that is logged but never leaked.
func ErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		var ve *lib.ErrValidationFailed
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(lib.WrapValidationError(ve.MessageRecord))
		}

		var fe *fiber.Error
		if errors.As(err, &fe) {
			return c.Status(fe.Code).JSON(fiber.Map{"message": fe.Message})
		}

		status, message := apperr.From(err)
		if status >= fiber.StatusInternalServerError {
			logger.Error("request failed", "method", c.Method(), "path", c.Path(), "error", err)
		}
		return c.Status(status).JSON(fiber.Map{"message": message})
	}
}
