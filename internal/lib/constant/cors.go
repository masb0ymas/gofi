package constant

import "gofi/internal/app"

func AllowedOrigins(app *app.Application) []string {
	var allowedOrigins []string

	// local development
	if app.Config.App.Env != "production" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	}

	// production
	allowedOrigins = append(allowedOrigins, "https://example.com")

	return allowedOrigins
}
