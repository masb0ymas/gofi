package middleware

import (
	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
)

func SentryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if hub := sentryfiber.GetHubFromContext(c); hub != nil {
			hub.Scope().SetTag("middleware", "sentry")
		}

		return c.Next()
	}
}
