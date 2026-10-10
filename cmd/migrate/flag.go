package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	modeUp      = "up"
	modeDown    = "down"
	modeRefresh = "refresh"
	modeForce   = "force"
)

type config struct {
	mode        string
	dbDSN       string
	seed        bool
	allowRemote bool
	version     int
}

func parseFlag(cfg *config) {
	// Defaults to the environment so the DSN, which carries the database
	// password, stays out of the process arguments.
	flag.StringVar(&cfg.dbDSN, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN (default $DB_DSN)")
	flag.BoolVar(&cfg.seed, "seed", false, "run the seeders after migrating")
	flag.BoolVar(&cfg.allowRemote, "allow-remote", false, "allow down and refresh on a database that is not on localhost")

	flag.Parse()
	validateFlag(cfg)
}

// isLocalDSN reports whether a URL-style DSN points at this machine. Anything
// it cannot parse counts as remote.
func isLocalDSN(dsn string) bool {
	u, err := url.Parse(dsn)
	if err != nil {
		return false
	}
	return slices.Contains([]string{"localhost", "127.0.0.1", "::1"}, u.Hostname())
}

func validateFlag(cfg *config) {
	cfg.mode = flag.Arg(0)
	availableModes := []string{modeUp, modeDown, modeRefresh, modeForce}

	if !slices.Contains(availableModes, cfg.mode) {
		log.Fatalf("command must be provided: %s", strings.Join(availableModes, ", "))
	}

	if cfg.dbDSN == "" {
		log.Fatal("flag --db-dsn or env DB_DSN must be provided")
	}

	if cfg.seed && (cfg.mode == modeDown || cfg.mode == modeForce) {
		log.Fatalf("flag --seed must not be provided when mode is %s", cfg.mode)
	}

	// down rolls back every migration, dropping the authula schema with all
	// users and sessions, so it takes an explicit flag off localhost.
	if (cfg.mode == modeDown || cfg.mode == modeRefresh) && !cfg.allowRemote && !isLocalDSN(cfg.dbDSN) {
		log.Fatalf("%s deletes all data and the database is not on localhost; pass --allow-remote to proceed", cfg.mode)
	}

	if cfg.mode == modeForce {
		version, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Fatal("force needs the schema version to set, e.g. `force 2` (-1 for none)")
		}
		cfg.version = version
	}
}
