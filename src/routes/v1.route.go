package routes

import (
	"gofi/src/apps/controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteV1(app *fiber.App) {
	// group v1
	v1 := app.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("Version", "v1")

		return c.Next()
	})

	// Endpoint Role
	roleHandler := v1.Group("/role")
	roleHandler.Get("/", controllers.FindAllRole)
}
