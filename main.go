package main

import (
	"gofi/src/pkg/config"
	"gofi/src/routes"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

func main() {
	app := fiber.New()

	port := config.Env("APP_PORT", "8000")
	envRateLimit := config.Env("RATE_LIMIT", "10")
	rateLimit, _ := strconv.Atoi(envRateLimit)

	// default middleware
	app.Use(cors.New(config.Cors()))
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(logger.New())
	app.Use(limiter.New(limiter.Config{Max: rateLimit}))
	app.Use(requestid.New())
	app.Use(recover.New())

	// static file
	app.Static("/", "./public")

	// Connect to the Database
	config.ConnectDB()

	// initial app route
	routes.InitializeRoutes(app)

	// listening app
	log.Fatal(app.Listen(":" + port))

}
