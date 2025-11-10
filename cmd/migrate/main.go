package main

import (
	"database/sql"
	"fmt"
	"gofi/internal/seeders"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var cfg config
	parseFlag(&cfg)

	db, err := connectDB(cfg.dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://./migrations/", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	switch cfg.mode {
	case modeUp:
		migrateUp(migrate)
	case modeDown:
		migrateDown(migrate)
	case modeRefresh:
		migrateDown(migrate)
		migrateUp(migrate)
	}

	if cfg.seed != "" {
		s := []seeders.Seeder{
			seeders.RoleSeeder{DB: db},
		}

		execSeeders(db, s...)
	}

	fmt.Println("Completed")
}

func migrateUp(m *migrate.Migrate) {
	fmt.Println("Running up migrations...")
	if err := m.Up(); err != nil {
		log.Fatalf("failed to run up migrations: %v", err)
	}
}

func migrateDown(m *migrate.Migrate) {
	fmt.Println("Running down migrations...")
	if err := m.Down(); err != nil {
		log.Fatalf("failed to run down migrations: %v", err)
	}
}

func execSeeders(db *sql.DB, seeders ...seeders.Seeder) {
	for _, seeder := range seeders {
		fmt.Printf("Running %s seeder...\n", seeder.Name())
		seeder.Seed()
	}
}
