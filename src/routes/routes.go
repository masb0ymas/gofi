package routes

import (
	"fmt"
	"gofi/src/pkg/helpers"
	"gofi/src/pkg/modules"
	"net/http"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

/*
Initialize Routes
*/
func InitializeRoutes(app *fiber.App) {
	// index route
	app.Get("/", func(c *fiber.Ctx) error {

		httpResponse := modules.HttpResponse(modules.Response{
			Message: "gofi - golang fiber v" + fiber.Version,
			Data: fiber.Map{
				"maintaner": "masb0ymas, <n.fajri@outlook.com>",
				"source":    "https://github.com/masb0ymas/gofi",
			}})

		return c.Status(http.StatusOK).JSON(httpResponse)
	})

	// health route
	app.Get("/health", func(c *fiber.Ctx) error {

		httpResponse := modules.HttpResponse(modules.Response{
			Data: fiber.Map{
				"cpu":     runtime.NumCPU(),
				"date":    helpers.TimeIn("ID").Format(time.RFC850),
				"golang":  runtime.Version(),
				"gofiber": fiber.Version,
				"status":  "Ok",
			}})

		return c.Status(http.StatusOK).JSON(httpResponse)
	})

	// monitor route
	app.Get("/monitor", monitor.New())

	// forbidden route version
	app.Get("/v1", func(c *fiber.Ctx) error {
		return c.Status(http.StatusForbidden).JSON(fiber.NewError(http.StatusForbidden))
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).JSON(fiber.NewError(http.StatusNotFound, "Sorry, HTTP resource you are looking for was not found."))
	})

	logMessage := helpers.PrintLog("Routes", "initialized successfully")
	fmt.Println(logMessage)
}
