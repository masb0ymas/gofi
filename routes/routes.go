package routes

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/masb0ymas/go-utils/pkg"
)

func Routes(db *sqlx.DB, app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message":    "Go Fi with Sqlx",
			"maintainer": "masb0ymas <n.fajri@mail.com>",
			"source":     "https://github.com/masb0ymas/gofi",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"cpu":     runtime.NumCPU(),
			"date":    pkg.TimeIn("ID").Format(time.RFC850),
			"golang":  runtime.Version(),
			"gofiber": fiber.Version,
			"status":  "Ok",
		})
	})

	app.Get("/v1", func(c *fiber.Ctx) error {
		return c.Status(http.StatusForbidden).JSON(fiber.NewError(http.StatusForbidden))
	})

	// initial v1 route
	v1Route(db, app)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(fiber.NewError(http.StatusForbidden, "Sorry, HTTP resource you are looking for was not found."))
	})
}
