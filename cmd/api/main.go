package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"gofi/internal/app"
	"gofi/internal/config"
	"gofi/internal/repositories"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg config.Config
	parseFlag(&cfg)

	level := slog.LevelInfo
	if cfg.App.Debug {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	if cfg.App.MigrateOnBoot {
		if err := migrateUp(&cfg.DB, logger); err != nil {
			return err
		}
	}

	db, err := connectDB(&cfg.DB)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	if cfg.App.SeedOnBoot {
		if err := seedAll(db, cfg.App.SeedPassword, logger); err != nil {
			return err
		}
	}

	// Creating the instance runs Authula's own migrations.
	auth := newAuth(cfg, logger)
	defer auth.CloseSystems()
	defer auth.ClosePlugins()

	return serve(&app.Application{
		Config:       cfg,
		Logger:       logger,
		Repositories: repositories.New(db, cfg.App.Debug),
		Auth:         auth,
	})
}
