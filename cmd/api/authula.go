package main

import (
	"log/slog"
	"time"

	authula "github.com/Authula/authula"
	authulaconfig "github.com/Authula/authula/config"
	authulamodels "github.com/Authula/authula/models"
	emailplugin "github.com/Authula/authula/plugins/email"
	emailpasswordplugin "github.com/Authula/authula/plugins/email-password"
	emailpasswordplugintypes "github.com/Authula/authula/plugins/email-password/types"
	emailplugintypes "github.com/Authula/authula/plugins/email/types"
	magiclinkplugin "github.com/Authula/authula/plugins/magic-link"
	magiclinkplugintypes "github.com/Authula/authula/plugins/magic-link/types"
	oauth2plugin "github.com/Authula/authula/plugins/oauth2"
	oauth2plugintypes "github.com/Authula/authula/plugins/oauth2/types"
	sessionplugin "github.com/Authula/authula/plugins/session"

	"gofi/internal/config"
)

// authBasePath is where the Authula handler is mounted on the Fiber app;
// every Authula route (and route mapping below) lives under it.
const authBasePath = "/v1/auth"

var (
	sessionAuth         = sessionplugin.HookIDSessionAuth.String()
	sessionAuthOptional = sessionplugin.HookIDSessionAuthOptional.String()
)

// authRouteMappings attaches session capabilities to Authula's routes. Paths
// are relative to authBasePath. Sign-in style routes take an optional session
// so an already signed-in user is recognised; account routes require one.
var authRouteMappings = []authulamodels.RouteMapping{
	{
		Paths:   []string{"GET:/me", "POST:/sign-out"},
		Plugins: []string{sessionAuth},
	},
	{
		Paths: []string{
			"POST:/email-password/sign-up",
			"POST:/email-password/sign-in",
			"GET:/email-password/verify-email",
			"POST:/email-password/request-password-reset",
			"POST:/email-password/change-password",
			"/magic-link/*",
			"/oauth2/*",
		},
		Plugins: []string{sessionAuthOptional},
	},
	{
		Paths: []string{
			"POST:/email-password/send-email-verification",
			"POST:/email-password/request-email-change",
		},
		Plugins: []string{sessionAuth},
	},
}

// newAuth builds the embedded Authula instance: email & password, Google
// OAuth2 and email magic links on top of cookie sessions, with emails sent
// through Resend. Authula migrates its own tables (users, accounts, sessions,
// verifications) into the app database when it starts.
func newAuth(cfg config.Config, logger *slog.Logger) *authula.Auth {
	origins := splitCSV(cfg.App.CORSAllowedOrigins)

	logLevel := "info"
	if cfg.App.Debug {
		logLevel = "debug"
	}

	authConfig := authulaconfig.NewConfig(
		authulaconfig.WithAppName(cfg.App.Name),
		authulaconfig.WithBaseURL(cfg.Auth.BaseURL),
		authulaconfig.WithBasePath(authBasePath),
		authulaconfig.WithSecret(cfg.App.Secret),
		authulaconfig.WithDatabase(authulamodels.DatabaseConfig{
			Provider:     "postgres",
			URL:          cfg.DB.DSN,
			MaxOpenConns: cfg.DB.MaxOpenConns,
			MaxIdleConns: cfg.DB.MaxIdleConns,
		}),
		authulaconfig.WithLogger(authulamodels.LoggerConfig{Level: logLevel, Logger: logger}),
		authulaconfig.WithSession(authulamodels.SessionConfig{
			CookieName:         config.SessionCookieName,
			ExpiresIn:          7 * 24 * time.Hour,
			UpdateAge:          24 * time.Hour,
			Secure:             cfg.App.Env != config.EnvDevelopment,
			HttpOnly:           true,
			SameSite:           "lax",
			AutoCleanup:        true,
			MaxSessionsPerUser: 5,
		}),
		authulaconfig.WithVerification(authulamodels.VerificationConfig{AutoCleanup: true}),
		// Empty origins leave Authula's CORS at its "*" default, which it never
		// answers with credentials: same-origin only, like the rest of the API.
		authulaconfig.WithSecurity(authulamodels.SecurityConfig{
			TrustedOrigins: origins,
			TrustedProxies: splitCSV(cfg.App.TrustedProxies),
			CORS: authulamodels.CORSConfig{
				AllowCredentials: true,
				AllowedOrigins:   origins,
			},
		}),
		authulaconfig.WithRouteMappings(authRouteMappings),
	)

	googleEnabled := cfg.Auth.GoogleClientID != ""
	if !googleEnabled {
		logger.Warn("google sign-in disabled: google-client-id is not set")
	}

	return authula.New(&authula.AuthConfig{
		Config: authConfig,
		Plugins: []authulamodels.Plugin{
			// Registered first: the email-based plugins look up its mailer.
			emailplugin.New(emailplugintypes.EmailPluginConfig{
				Enabled:     true,
				Provider:    emailplugintypes.ProviderResend,
				FromAddress: cfg.Auth.EmailFrom,
				Resend:      &emailplugintypes.ResendConfig{ApiKey: cfg.Auth.ResendAPIKey},
			}),
			sessionplugin.New(sessionplugin.SessionPluginConfig{Enabled: true}),
			emailpasswordplugin.New(emailpasswordplugintypes.EmailPasswordPluginConfig{
				Enabled:                  true,
				MinPasswordLength:        8,
				MaxPasswordLength:        128,
				RequireEmailVerification: true,
				AutoSignIn:               true,
				SendEmailOnSignUp:        true,
				PasswordResetExpiresIn:   time.Hour,
			}),
			magiclinkplugin.New(magiclinkplugintypes.MagicLinkPluginConfig{
				Enabled:   true,
				ExpiresIn: 15 * time.Minute,
			}),
			oauth2plugin.New(oauth2plugintypes.OAuth2PluginConfig{
				Enabled: googleEnabled,
				Providers: map[string]oauth2plugintypes.ProviderConfig{
					"google": {
						Enabled:      googleEnabled,
						ClientID:     cfg.Auth.GoogleClientID,
						ClientSecret: cfg.Auth.GoogleClientSecret,
						RedirectURL:  authConfig.BaseURL + authBasePath + "/oauth2/callback/google",
					},
				},
			}),
		},
	})
}
