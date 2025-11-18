package main

import (
	"log/slog"
	"os"

	"gofi/internal/app"
	"gofi/internal/config"
	"gofi/internal/repositories"
	"gofi/internal/services"
)

func main() {
	var cfg config.Config
	parseFlag(&cfg)

	loggerLevel := slog.LevelInfo

	if cfg.App.Debug {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	db, err := connectDB(&cfg.DB)
	if err != nil {
		logger.Error("failed to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Dependencies Injection
	app := &app.Application{
		Config:       cfg,
		Logger:       logger,
		Repositories: repositories.New(db),
		Services: services.Services{
			Email: services.EmailService{Config: cfg.Resend},
		},
	}

	if err := serve(app); err != nil {
		logger.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}
