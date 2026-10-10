package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

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
	// The repositories log their queries through the default logger.
	slog.SetDefault(logger)

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

	repos := repositories.New(db)
	mirror := mirrorUser(&repos.User)

	// Creating the instance runs Authula's own migrations.
	auth := newAuth(cfg, mirror, logger)
	defer auth.CloseSystems()
	defer auth.ClosePlugins()

	// Repairs mirror rows a failed hook left missing or stale.
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	err = repos.User.Reconcile(ctx)
	cancel()
	if err != nil {
		return fmt.Errorf("reconcile users mirror: %w", err)
	}

	// Seeding follows Authula because the super user is created through it.
	if cfg.App.SeedOnBoot {
		superUser := seeders.SuperUserSeeder{
			DB:       db,
			Auth:     auth.CoreServices(),
			Promote:  promoteSuperUser(&repos.User, mirror),
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
		DB:           db,
		Repositories: repos,
		Auth:         auth,
	})
}
