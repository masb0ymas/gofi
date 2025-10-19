package routes

import (
	"gofi/src/app/docs"
	"gofi/src/config"
	"gofi/src/lib"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/masb0ymas/go-utils/pkg"
)

func Root(db *sqlx.DB, app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message":    "Gofi Api",
			"maintainer": "masb0ymas <n.fajri@mail.com>",
			"source":     "https://github.com/masb0ymas/gofi",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")

		return c.Status(http.StatusOK).JSON(fiber.Map{
			"cpu":        runtime.NumCPU(),
			"date":       pkg.TimeIn("ID").Format(time.RFC850),
			"golang":     runtime.Version(),
			"gofiber":    fiber.Version,
			"status":     "Ok",
			"client ip":  c.IP(),
			"user agent": userAgent,
		})
	})

	app.Get("/v1", func(c *fiber.Ctx) error {
		return lib.SendForbiddenResponse(c, fiber.NewError(http.StatusForbidden).Error())
	})

	// api reference not include for production "/docs"
	if config.Env("APP_ENV", "development") != "production" {
		docs.SetupDocsRoutes(app)
	}

	// initial v1 route
	v1Route(db, app)

	app.Use("*", func(c *fiber.Ctx) error {
		return lib.SendNotFoundResponse(c, "Sorry, HTTP resource you are looking for was not found.")
	})
}
