package main

import (
	"gofi/internal/app"
	"gofi/internal/handlers"
	"gofi/internal/lib/constant"
	"gofi/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func routes(r *fiber.App, app *app.Application) {
	h := handlers.New(app)
	m := middlewares.New(app)

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	r.Get("/healthcheck", h.Health.Check)

	// Public routes
	r.Get("/v1/roles", h.Role.Index)
	r.Get("/v1/roles/:roleID", h.Role.Show)

	// Admin protected routes
	adminAuthorized := r.Group("/")
	adminAuthorized.Use(m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}))

	adminAuthorized.Post("/v1/roles", h.Role.Create)
	adminAuthorized.Put("/v1/roles/:roleID", h.Role.Update)
	adminAuthorized.Delete("/v1/roles/:roleID", h.Role.Delete)
	adminAuthorized.Delete("/v1/roles/:roleID/soft-delete", h.Role.SoftDelete)
	adminAuthorized.Patch("/v1/roles/:roleID/restore", h.Role.Restore)

	adminAuthorized.Get("/v1/users", h.User.Index)
	adminAuthorized.Get("/v1/users/:userID", h.User.Show)
	adminAuthorized.Post("/v1/users", h.User.Create)
	adminAuthorized.Put("/v1/users/:userID", h.User.Update)
	adminAuthorized.Delete("/v1/users/:userID", h.User.Delete)
	adminAuthorized.Delete("/v1/users/:userID/soft-delete", h.User.SoftDelete)
	adminAuthorized.Patch("/v1/users/:userID/restore", h.User.Restore)

	// Not found handler
	r.Get("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Sorry, HTTP resource you are looking for was not found.",
		})
	})
}
