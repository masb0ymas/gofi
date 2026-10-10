package config

import "time"

// Accepted App.Env values. Only EnvDevelopment enables development-only
// behaviour: the loopback rate-limit exemption and session cookies without
// Secure. The default is EnvProduction, so forgetting ENV fails safe.
const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
)

// SessionCookieName is the cookie Authula issues on sign-in; the API's
// authorization middleware reads the same cookie.
const SessionCookieName = "authula.session_token"

type Config struct {
	App    ConfigApp
	DB     ConfigDB
	Auth   ConfigAuth
	Sentry ConfigSentry
}

// ConfigAuth configures the embedded Authula server.
type ConfigAuth struct {
	BaseURL            string
	EmailFrom          string
	ResendAPIKey       string
	GoogleClientID     string
	GoogleClientSecret string
}

type ConfigApp struct {
	Env                string
	Debug              bool
	Port               int
	Name               string
	Secret             string
	MigrateOnBoot      bool
	SeedOnBoot         bool
	SuperUserEmail     string
	SuperUserPassword  string
	CORSAllowedOrigins string
	RateLimitExemptIPs string
	TrustedProxies     string
}

type ConfigDB struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type ConfigSentry struct {
	Dsn string
}
