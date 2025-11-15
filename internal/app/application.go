package app

import (
	"log/slog"

	"gofi/internal/config"
	"gofi/internal/repositories"
)

type Application struct {
	Config       config.Config
	Logger       *slog.Logger
	Repositories repositories.Repositories
}
