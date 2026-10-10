package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"gofi/internal/app"
	"gofi/internal/config"
	"gofi/internal/repositories"
	"gofi/internal/seeders"
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

	repos := repositories.New(db, cfg.App.Debug)
	mirror := mirrorUser(&repos.User)

	// Creating the instance runs Authula's own migrations.
	auth := newAuth(cfg, mirror, logger)
	defer auth.CloseSystems()
	defer auth.ClosePlugins()

	// Seeding follows Authula because the super user is created through it.
	if cfg.App.SeedOnBoot {
		superUser := seeders.SuperUserSeeder{
			DB:       db,
			Auth:     auth.CoreServices(),
			Mirror:   mirror,
			Email:    cfg.App.SuperUserEmail,
			Password: cfg.App.SuperUserPassword,
		}
		if err := seedAll(db, superUser, logger); err != nil {
			return err
		}
	}

	return serve(&app.Application{
		Config:       cfg,
		Logger:       logger,
		Repositories: repos,
		Auth:         auth,
	})
}
