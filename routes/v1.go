package routes

import (
	"gofi/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func v1Route(db *sqlx.DB, app *fiber.App) {
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")
		return c.Next()
	})

	handler.RoleHandler(db, v1)
	handler.SessionHandler(db, v1)
}
