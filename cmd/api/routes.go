package main

import (
	fiberadapter "github.com/Authula/authula/adapters/fiber"
	"github.com/gofiber/fiber/v3"

	"gofi/internal/app"
	"gofi/internal/config"
	"gofi/internal/handlers"
	"gofi/internal/middlewares"
)

// authRateLimit is the per-IP request budget per minute on the auth routes.
const authRateLimit = 30

func routes(r *fiber.App, app *app.Application) {
	h := handlers.New(app)
	m := middlewares.New(app)

	r.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})
	r.Get("/health", h.Health.Check)

	// Authula: email & password, Google OAuth2 and magic-link sign-in,
	// plus GET /me and POST /sign-out, all under authBasePath. The limiter
	// guards password guessing and sign-in/reset email flooding.
	r.Use(authBasePath,
		middlewares.RateLimit(authRateLimit, app.Config.App.Env == config.EnvDevelopment, splitCSV(app.Config.App.RateLimitExemptIPs)),
		fiberadapter.New(fiberadapter.Config{Handler: app.Auth.Handler()}),
	)

	v1 := r.Group("/v1")

	// Reads are public; writes need a signed-in user.
	category := v1.Group("/category")
	category.Get("/", h.Category.List)
	category.Get("/:id", h.Category.Get)
	category.Post("/", m.Authorization, h.Category.Create)
	category.Put("/:id", m.Authorization, h.Category.Update)
	category.Delete("/:id", m.Authorization, h.Category.Delete)
	category.Put("/:id/restore", m.Authorization, h.Category.Restore)
	category.Delete("/:id/force", m.Authorization, h.Category.DeletePermanently)
}
