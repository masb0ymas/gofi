package middlewares

import (
	"fmt"
	"net/http"

	"gofi/internal/lib"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m Middlewares) PermissionAccess(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid := lib.ContextGetUID(c)

		user, err := m.app.Repositories.User.GetByID(uid)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": fmt.Sprintf("Unauthorized, permission access failed: %s", err.Error()),
			})
		}

		if user.ID != uuid.Nil && !lib.Contains(roles, user.RoleID.String()) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized, permission access failed: you are not allowed!",
			})
		}

		return c.Next()
	}
}
