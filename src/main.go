package main

import (
	"gofi/src/app/database"
	"gofi/src/app/routes"
	"gofi/src/config"
	"gofi/src/lib"
	"log"
	"os"
	"time"

	sentryfiber "github.com/getsentry/sentry-go/fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/masb0ymas/go-utils/pkg"
)

func main() {
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = database.DB.Close() }()

	// Initial Provider
	// lib.InitGCS()

	if err := lib.InitSentry(); err != nil {
		log.Fatalf("failed to initialize Sentry: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:               10 * 1024 * 1024, // 10MB
		StreamRequestBody:       true,
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(config.Cors()))
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Use(requestid.New())
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too Many Requests",
			})
		},
	}))

	app.Static("/", "./public")

	// Sentry
	sentryHandler := sentryfiber.New(sentryfiber.Options{
		Repanic:         true,
		WaitForDelivery: true,
	})
	app.Use(sentryHandler)

	// Routes
	routes.Root(database.DB, app)

	// Get port from environment variable or use default
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	// Start server
	msg := pkg.Println("Fiber", "Server starting on port "+port+"")
	log.Println(msg)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
