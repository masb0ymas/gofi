package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"gofi/internal/lib"
	"gofi/internal/lib/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m Middlewares) Authorization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jsonWebToken := jwt.New(&m.app.Config.App)

		extractToken, err := jsonWebToken.ExtractToken(c)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("Unauthorized, %s", err.Error()),
			})
		}

		session, err := m.app.Repositories.Session.GetByToken(extractToken)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("Unauthorized, %s", err.Error()),
			})
		}

		if session.ID != uuid.Nil {
			claims, err := jsonWebToken.Verify(extractToken)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": fmt.Sprintf("Unauthorized, %s", err.Error()),
				})
			}

			if claims.UID != session.UserID.String() {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized, invalid session",
				})
			}

			if claims.Exp < time.Now().Unix() {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": fmt.Sprintf("Unauthorized, %s", "token has expired"),
				})
			}

			uid, err := uuid.Parse(claims.UID)
			if err != nil {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": "Unauthorized, invalid session",
				})
			}

			lib.ContextSetUID(c, uid)
		}

		return c.Next()
	}
}
