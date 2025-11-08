package handlers

import (
	"gofi/internal/app"

	"github.com/gofiber/fiber/v2"
)

type healthHandler struct {
	app *app.Application
}

func (h *healthHandler) Check(c *fiber.Ctx) error {
	v := fiber.Map{
		"machineID": h.app.Config.App.MachineID,
		"status":    "ok",
		"systemInfo": map[string]interface{}{
			"debug": h.app.Config.App.Debug,
		},
	}

	return c.Status(fiber.StatusOK).JSON(v)
}
