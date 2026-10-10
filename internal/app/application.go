package app

import (
	"database/sql"
	"log/slog"

	authula "github.com/Authula/authula"

	"gofi/internal/config"
	"gofi/internal/repositories"
)

type Application struct {
	Config       config.Config
	Logger       *slog.Logger
	DB           *sql.DB
	Repositories repositories.Repositories
	Auth         *authula.Auth
}
