package main

import (
	"flag"
	"log"
	"slices"
	"strings"
)

const (
	modeUp      = "up"
	modeDown    = "down"
	modeRefresh = "refresh"
)

const (
	seedProduction  = "prod"
	seedDevelopment = "dev"
)

type config struct {
	mode   string
	dbDSN  string
	seed   string
	appPwd string
}

func parseFlag(cfg *config) {
	flag.StringVar(&cfg.dbDSN, "db-dsn", "", "PostgreSQL DSN")
	flag.StringVar(&cfg.seed, "seed", "", "seed")
	flag.StringVar(&cfg.appPwd, "app-pwd", "", "app default password")

	flag.Parse()
	validateFlag(cfg)
}

func validateFlag(cfg *config) {
	cfg.mode = flag.Arg(0)
	availableModes := []string{modeUp, modeDown, modeRefresh}

	if !slices.Contains(availableModes, cfg.mode) {
		log.Fatalf("command must be provided: %s", strings.Join(availableModes, ", "))
	}

	if cfg.dbDSN == "" {
		log.Fatal("flag --db-dsn must be provided")
	}

	if cfg.appPwd == "" {
		log.Fatal("flag --app-pwd must be provided")
	}

	if cfg.mode == modeDown && cfg.seed != "" {
		log.Fatal("flag --seed must not be provided when mode is down")
	}

	if cfg.seed != "" {
		availableSeeds := []string{seedProduction, seedDevelopment}
		if !slices.Contains(availableSeeds, cfg.seed) {
			log.Fatalf("seed must be provided: %s", strings.Join(availableSeeds, ", "))
		}
	}
}
