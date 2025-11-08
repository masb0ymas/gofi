package handlers

import (
	"gofi/internal/app"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type roleHandler struct {
	app *app.Application
}

func (h *roleHandler) Index(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Index",
	})
}
