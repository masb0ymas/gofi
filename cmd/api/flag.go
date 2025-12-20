package main

import (
	"flag"
	"log"
	"time"

	"gofi/internal/config"
)

func parseFlag(cfg *config.Config) {
	var machineID uint

	// App
	flag.UintVar(&machineID, "machine-id", 0, "Machine ID")
	flag.StringVar(&cfg.App.Env, "env", "development", "Environment")
	flag.BoolVar(&cfg.App.Debug, "debug", false, "Debug mode")
	flag.IntVar(&cfg.App.Port, "port", 8080, "Port")
	flag.StringVar(&cfg.App.Name, "app-name", "gofi", "App Name")
	flag.StringVar(&cfg.App.Secret, "app-secret", "", "App Secret")
	flag.StringVar(&cfg.App.JWTSecret, "jwt-secret", "", "JWT Secret")
	flag.StringVar(&cfg.App.ClientURL, "client-url", "", "Client URL")
	flag.StringVar(&cfg.App.ServerURL, "server-url", "", "Server URL")

	// Database
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "Database DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "Database max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "Database max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "Database max idle time")

	// Redis
	flag.StringVar(&cfg.Redis.Addr, "redis-addr", "localhost:6379", "Redis address")
	flag.StringVar(&cfg.Redis.Password, "redis-password", "", "Redis password")
	flag.IntVar(&cfg.Redis.DB, "redis-db", 0, "Redis database")

	// Resend
	flag.StringVar(&cfg.Resend.ApiKey, "resend-api-key", "", "Resend API key")
	flag.StringVar(&cfg.Resend.FromEmail, "resend-from-email", "", "Resend from email")
	flag.StringVar(&cfg.Resend.DebugToEmail, "resend-debug-to-email", "", "Resend debug to email")

	// Google
	flag.StringVar(&cfg.Google.ClientID, "google-client-id", "", "Google client ID")
	flag.StringVar(&cfg.Google.ClientSecret, "google-client-secret", "", "Google client secret")
	flag.StringVar(&cfg.Google.RedirectURL, "google-redirect-url", "", "Google redirect URL")

	// S3
	flag.StringVar(&cfg.S3.ClientID, "s3-client-id", "", "S3 client ID")
	flag.StringVar(&cfg.S3.ClientSecret, "s3-client-secret", "", "S3 client secret")
	flag.StringVar(&cfg.S3.Region, "s3-region", "", "S3 region")
	flag.StringVar(&cfg.S3.Endpoint, "s3-endpoint", "", "S3 endpoint")
	flag.StringVar(&cfg.S3.Token, "s3-token", "", "S3 token")

	flag.Parse()

	uint16Max := uint(1<<16 - 1)
	if machineID > uint16Max {
		log.Fatal("flag machine-id can only handle uint16")
		return
	}

	cfg.App.MachineID = uint16(machineID)

	validateFlag(cfg)
}

func validateFlag(cfg *config.Config) {
	if cfg.App.Env == "" {
		log.Println("flag environment is marked as local")
	}

	if cfg.App.MachineID == 0 {
		log.Fatal("flag machine-id must be provided and cannot be 0")
	}

	if cfg.App.Secret == "" {
		log.Fatal("flag app-secret must be provided")
	}

	if cfg.App.JWTSecret == "" {
		log.Fatal("flag jwt-secret must be provided")
	}

	if cfg.App.ClientURL == "" {
		log.Fatal("flag client-url must be provided")
	}

	if cfg.App.ServerURL == "" {
		log.Fatal("flag server-url must be provided")
	}

	if cfg.DB.DSN == "" {
		log.Fatal("flag db-dsn must be provided")
	}

	if cfg.Redis.Addr == "" {
		log.Fatal("flag redis-addr must be provided")
	}

	if cfg.Resend.ApiKey == "" {
		log.Fatal("flag resend-api-key must be provided")
	}

	if cfg.Resend.FromEmail == "" {
		log.Fatal("flag resend-from-email must be provided")
	}

	if cfg.Resend.DebugToEmail == "" {
		log.Fatal("flag resend-debug-to-email must be provided")
	}
}
