package middlewares

import (
	"slices"
	"time"

	"gofi/internal/lib/apperr"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func RateLimit(max int, exemptLoopback bool, exemptIPs []string) fiber.Handler {
	return limiter.New(limiter.Config{
		Next: func(c fiber.Ctx) bool {
			return slices.Contains(exemptIPs, c.IP()) || (exemptLoopback && c.IP() == "127.0.0.1")
		},
		Max:        max,
		Expiration: time.Minute,
		LimitReached: func(c fiber.Ctx) error {
			return apperr.ErrTooManyRequests
		},
	})
}
