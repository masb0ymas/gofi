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

	r.Get("/health-check", h.Health.Check)

	authRoutes := r.Group("/v1/auth")
	authRoutes.Post("/sign-up", h.Auth.SignUp)
	authRoutes.Post("/sign-in", h.Auth.SignIn)
	authRoutes.Post("/verify-registration", h.Auth.VerifyRegistration)
	authRoutes.Get("/verify-session", m.Authorization(), h.Auth.VerifySession)
	authRoutes.Post("/refresh-token", m.Authorization(), h.Auth.RefreshToken)
	authRoutes.Post("/sign-out", m.Authorization(), h.Auth.SignOut)
	authRoutes.Post("/google", h.Auth.GoogleAuthURL)
	authRoutes.Get("/google/callback", h.Auth.GoogleAuthCallback)

	adminOnly := []string{constant.RoleAdmin}

	sessionRoutes := r.Group("/v1/sessions")
	sessionRoutes.Use(m.Authorization(), m.PermissionAccess(adminOnly))
	sessionRoutes.Get("", h.Session.Index)

	roleRoutes := r.Group("/v1/roles")
	roleRoutes.Use(m.Authorization())
	roleRoutes.Get("", h.Role.Index)
	roleRoutes.Get("/:roleID", h.Role.Show)
	roleRoutes.Post("", m.PermissionAccess(adminOnly), h.Role.Create)
	roleRoutes.Put("/:roleID", m.PermissionAccess(adminOnly), h.Role.Update)
	roleRoutes.Delete("/:roleID", m.PermissionAccess(adminOnly), h.Role.Delete)
	roleRoutes.Delete("/:roleID/soft-delete", m.PermissionAccess(adminOnly), h.Role.SoftDelete)
	roleRoutes.Patch("/:roleID/restore", m.PermissionAccess(adminOnly), h.Role.Restore)

	userRoutes := r.Group("/v1/users")
	userRoutes.Use(m.Authorization())
	userRoutes.Get("", h.User.Index)
	userRoutes.Get("/:userID", h.User.Show)
	userRoutes.Post("", m.PermissionAccess(adminOnly), h.User.Create)
	userRoutes.Put("/:userID", m.PermissionAccess(adminOnly), h.User.Update)
	userRoutes.Delete("/:userID", m.PermissionAccess(adminOnly), h.User.Delete)
	userRoutes.Delete("/:userID/soft-delete", m.PermissionAccess(adminOnly), h.User.SoftDelete)
	userRoutes.Patch("/:userID/restore", m.PermissionAccess(adminOnly), h.User.Restore)

	// Not found handler
	r.Use("*", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Sorry, HTTP resource you are looking for was not found.",
		})
	})
}
