package main

import (
	fiberadapter "github.com/Authula/authula/adapters/fiber"
	"github.com/gofiber/fiber/v3"

	"gofi/internal/app"
	"gofi/internal/config"
	"gofi/internal/handlers"
	"gofi/internal/middlewares"
)

// Per-IP request budgets per minute: tight on the auth routes, loose on the
// rest of the API.
const (
	authRateLimit = 30
	apiRateLimit  = 120
)

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

	exemptLoopback := app.Config.App.Env == config.EnvDevelopment
	exemptIPs := splitCSV(app.Config.App.RateLimitExemptIPs)

	// ponytail: the limiters count in memory, so each replica has its own
	// budget; give them shared storage once there is more than one instance.
	v1 := r.Group("/v1", middlewares.RateLimit(apiRateLimit, exemptLoopback, exemptIPs))

	// Reads are public; writes need a signed-in super user. Idempotency comes
	// last, after the caller is known, and replays a repeated write instead of
	// applying it twice (see middlewares.Idempotency). One guard serves every
	// write route: its keys are already scoped per route, and a shared cache
	// means a single in-memory store instead of one per route.
	idempotent := middlewares.Idempotency()

	categories := v1.Group("/categories")
	categories.Get("/", h.Category.List)
	categories.Get("/:id", h.Category.Get)
	categories.Post("/", m.Authorization, m.SuperUser, idempotent, h.Category.Create)
	categories.Put("/:id", m.Authorization, m.SuperUser, idempotent, h.Category.Update)
	categories.Delete("/:id", m.Authorization, m.SuperUser, idempotent, h.Category.Delete)
	categories.Put("/:id/restore", m.Authorization, m.SuperUser, idempotent, h.Category.Restore)
	categories.Delete("/:id/force", m.Authorization, m.SuperUser, idempotent, h.Category.DeletePermanently)
}
