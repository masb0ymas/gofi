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
	flag.StringVar(&cfg.App.JWTSecret, "jwt-secret", "", "JWT Secret")

	// Database
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "Database DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "Database max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "Database max idle connections")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "Database max idle time")

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

	if cfg.App.JWTSecret == "" {
		log.Fatal("flag jwt-secret must be provided")
	}

	if cfg.DB.DSN == "" {
		log.Fatal("flag db-dsn must be provided")
	}
}
