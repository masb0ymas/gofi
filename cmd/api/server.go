package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"gofi/internal/app"
	"gofi/internal/docs"
	"gofi/internal/lib/constant"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func serve(app *app.Application) error {
	// Fiber Configuration
	server := fiber.New(fiber.Config{
		BodyLimit:               2 * 1024 * 1024, // 2MB
		IdleTimeout:             time.Minute,
		ReadTimeout:             20 * time.Second,
		WriteTimeout:            3 * time.Minute,
		EnableTrustedProxyCheck: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	server.Use(recover.New())
	server.Use(logger.New())
	server.Use(helmet.New())
	server.Use(requestid.New())
	server.Use(compress.New())

	// CORS
	server.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(constant.AllowedOrigins(app), ","),
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		MaxAge:       3600,
	}))

	// Rate Limit
	server.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	server.Static("/", "./public")

	// generate /docs with scalar
	if app.Config.App.Env != "production" {
		docs.SetupDocsRoutes(server, app)
	}

	// Initial Routes
	routes(server, app)

	// Create channel to listen for interrupt signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		app.Logger.Info("server started on port", "port", app.Config.App.Port)
		listerPort := fmt.Sprintf(":%d", app.Config.App.Port)

		if err := server.Listen(listerPort); err != nil {
			app.Logger.Error("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-c
	app.Logger.Info("Received interrupt signal, shutting down...")

	// Stop server
	if err := server.Shutdown(); err != nil {
		app.Logger.Error("failed to stop server", "error", err)
	}

	return nil
}
