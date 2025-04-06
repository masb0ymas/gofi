package lib

import (
	"errors"
	"gofi/config"

	"github.com/getsentry/sentry-go"
)

func InitSentry() error {
	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.Env("SENTRY_DSN", ""),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		return errors.New("sentry initialization failed: " + err.Error())
	}

	return nil
}
