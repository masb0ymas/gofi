package app

import (
	"log/slog"

	"gofi/internal/config"
	"gofi/internal/repositories"
	"gofi/internal/services"
)

type Application struct {
	Config       config.Config
	Logger       *slog.Logger
	Repositories repositories.Repositories
	Services     services.Services
}
