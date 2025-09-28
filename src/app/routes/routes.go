package routes

import (
	"fmt"
	"gofi/src/config"
	"gofi/src/lib"
	"net/http"
	"runtime"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
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

	app.Get("/api-docs", func(c *fiber.Ctx) error {
		url := config.Env("APP_SERVER_URL", "http://localhost:8000")

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: fmt.Sprintf("%s/docs/swagger.json", url),
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Docs Go-Fi API",
			},
			DarkMode: true,
		})

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Type("html").SendString(htmlContent)
	})

	// initial v1 route
	v1Route(db, app)

	app.Use("*", func(c *fiber.Ctx) error {
		return lib.SendNotFoundResponse(c, "Sorry, HTTP resource you are looking for was not found.")
	})
}
