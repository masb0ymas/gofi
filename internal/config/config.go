package config

import (
	"cmp"
	"os"
	"strconv"
	"time"
)

// EnvDevelopment is the App.Env value that enables development-only behaviour:
// the loopback rate-limit exemption and session cookies without Secure.
const EnvDevelopment = "development"

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
	MachineID          string
	Env                string
	Debug              bool
	Port               int
	Name               string
	Secret             string
	MigrateOnBoot      bool
	SeedOnBoot         bool
	SeedPassword       string
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

// Load reads the settings the binaries use today from the environment.
func Load() Config {
	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		port = 8000
	}

	return Config{
		App: ConfigApp{
			Env:           cmp.Or(os.Getenv("APP_ENV"), EnvDevelopment),
			Port:          port,
			Name:          cmp.Or(os.Getenv("APP_NAME"), "gofi"),
			MigrateOnBoot: os.Getenv("MIGRATE_ON_BOOT") == "true",
		},
		DB: ConfigDB{
			DSN: cmp.Or(os.Getenv("DB_DSN"), "postgres://postgres:postgres@localhost:5432/gofi?sslmode=disable"),
		},
	}
}
