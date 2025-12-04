package main

import (
	"fmt"

	"gofi/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func newGoogleOAuth(cfg config.ConfigGoogle) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/v1/auth/google/callback", cfg.RedirectURL),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
}
