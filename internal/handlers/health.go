package handlers

import (
	"github.com/gofiber/fiber/v3"

	"gofi/internal/app"
	"gofi/internal/dtos"
)

type healthHandler struct {
	app *app.Application
}

// Check pings the database, so the probe fails when the API cannot serve.
// It reports nothing but the status: the endpoint is public.
func (h *healthHandler) Check(c fiber.Ctx) error {
	if err := h.app.DB.PingContext(c.Context()); err != nil {
		h.app.Logger.Error("health check: database unreachable", "error", err)
		return dtos.Response(c, fiber.StatusServiceUnavailable, fiber.Map{"status": "unavailable"}, nil, "")
	}
	return dtos.OK(c, fiber.Map{"status": "ok"})
}
