package main

import (
	"gofi/internal/app"
	"gofi/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func routes(r *fiber.App, app *app.Application) {
	h := handlers.New(app)

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	r.Get("/healthcheck", h.Health.Check)

	r.Get("/v1/roles", h.Role.Index)
	r.Get("/v1/roles/:roleID", h.Role.Show)
	r.Post("/v1/roles", h.Role.Create)
	r.Put("/v1/roles/:roleID", h.Role.Update)
	r.Delete("/v1/roles/:roleID", h.Role.Delete)
	r.Delete("/v1/roles/:roleID/soft-delete", h.Role.SoftDelete)
	r.Patch("/v1/roles/:roleID/restore", h.Role.Restore)
}
