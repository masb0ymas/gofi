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

	r.Post("/v1/auth/sign-up", h.Auth.SignUp)
	r.Post("/v1/auth/sign-in", h.Auth.SignIn)
	r.Get("/v1/auth/verify-session", h.Auth.VerifySession)
	r.Post("/v1/auth/sign-out", h.Auth.SignOut)

	r.Get("/v1/roles", h.Role.Index)
	r.Get("/v1/roles/:roleID", h.Role.Show)
	r.Post("/v1/roles", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.Role.Create)
	r.Put("/v1/roles/:roleID", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.Role.Update)
	r.Delete("/v1/roles/:roleID", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.Role.Delete)
	r.Delete("/v1/roles/:roleID/soft-delete", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.Role.SoftDelete)
	r.Patch("/v1/roles/:roleID/restore", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.Role.Restore)

	r.Get("/v1/users", h.User.Index)
	r.Get("/v1/users/:userID", h.User.Show)
	r.Post("/v1/users", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.User.Create)
	r.Put("/v1/users/:userID", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.User.Update)
	r.Delete("/v1/users/:userID", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.User.Delete)
	r.Delete("/v1/users/:userID/soft-delete", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.User.SoftDelete)
	r.Patch("/v1/users/:userID/restore", m.Authorization(), m.PermissionAccess([]string{constant.RoleAdmin}), h.User.Restore)

	// Not found handler
	r.Get("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Sorry, HTTP resource you are looking for was not found.",
		})
	})
}
