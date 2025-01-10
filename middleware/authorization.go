package middleware

import (
	"goarif-api/config"
	"goarif-api/lib"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		secretKey := config.Env("JWT_SECRET_KEY", "secret")
		claims, err := lib.VerifyToken(c, secretKey)
		if err != nil {
			return lib.SendUnauthorizedResponse(c, err.Error())
		}

		uid, err := uuid.Parse(claims.UID)
		if err != nil {
			return lib.SendUnauthorizedResponse(c, "invalid UID in token")
		}

		if claims.Exp < time.Now().Unix() {
			return lib.SendUnauthorizedResponse(c, "token is invalid")
		}

		c.Set("uid", uid.String())

		return c.Next()
	}
}
