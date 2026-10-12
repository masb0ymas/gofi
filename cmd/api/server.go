package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryfiber "github.com/gofiber/contrib/v3/sentry"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/static"

	"gofi/internal/app"
	"gofi/internal/middlewares"
)

// serve runs the HTTP server until SIGINT/SIGTERM, then drains in-flight
// requests before returning.
func serve(app *app.Application) error {
	server, err := newServer(app)
	if err != nil {
		return err
	}
	// A no-op unless newServer initialised Sentry.
	defer sentry.Flush(5 * time.Second)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.Logger.Info("starting server", "port", app.Config.App.Port, "env", app.Config.App.Env, "auth", authBasePath)
	return server.Listen(fmt.Sprintf(":%d", app.Config.App.Port), fiber.ListenConfig{GracefulContext: ctx})
}

// newServer builds the Fiber app with its middleware and routes.
func newServer(app *app.Application) (*fiber.App, error) {
	proxies := splitCSV(app.Config.App.TrustedProxies)

	server := fiber.New(fiber.Config{
		AppName:          app.Config.App.Name,
		ErrorHandler:     middlewares.ErrorHandler(app.Logger),
		TrustProxy:       len(proxies) > 0,
		TrustProxyConfig: fiber.TrustProxyConfig{Proxies: proxies},
		ProxyHeader:      fiber.HeaderXForwardedFor,
		// Bound how long a slow client can hold a connection.
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Middleware
	server.Use(recover.New())
	server.Use(logger.New())
	server.Use(helmet.New())
	server.Use(requestid.New())
	server.Use(compress.New())

	// Sentry reports panics, and only when a DSN is configured.
	if dsn := app.Config.Sentry.Dsn; dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{Dsn: dsn, Environment: app.Config.App.Env}); err != nil {
			return nil, fmt.Errorf("init sentry: %w", err)
		}

		server.Use(sentryfiber.New(sentryfiber.Config{
			Repanic:         true,
			WaitForDelivery: true,
			Timeout:         5 * time.Second,
		}))
	}

	origins := splitCSV(app.Config.App.CORSAllowedOrigins)

	// CSRF: the session cookie is SameSite=Lax, which a sibling subdomain can
	// get around. This covers Authula's routes too, as they are mounted here.
	server.Use(middlewares.CrossOrigin(origins))

	// Authula answers CORS for its own routes from the same origin list, so
	// the Fiber middleware skips them to avoid duplicate headers.
	if len(origins) > 0 {
		server.Use(cors.New(cors.Config{
			Next:             func(c fiber.Ctx) bool { return strings.HasPrefix(c.Path(), authBasePath) },
			AllowOrigins:     origins,
			AllowCredentials: true,
			// Without this a browser client cannot read the marker that says a
			// write was replayed instead of applied (see middlewares.Idempotency).
			ExposeHeaders: []string{middlewares.HeaderIdempotentReplayed},
		}))
	}

	server.Use(static.New("./public"))

	routes(server, app)

	return server, nil
}
