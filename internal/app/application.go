package app

import (
	"log/slog"

	"gofi/internal/config"
)

type Application struct {
	Config config.Config
	Logger *slog.Logger
}
