package app

import (
	"log/slog"

	authula "github.com/Authula/authula"

	"gofi/internal/config"
	"gofi/internal/repositories"
)

var Version = "dev"

type Application struct {
	Config       config.Config
	Logger       *slog.Logger
	Repositories repositories.Repositories
	Auth         *authula.Auth
}
