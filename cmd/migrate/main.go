package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func main() {
	var cfg config
	parseFlag(&cfg)

	db, err := connectDB(cfg.dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
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

	fmt.Println("Completed")
}

func migrateUp(m *migrate.Migrate) {
	fmt.Println("Running up migrations...")
	if err := m.Up(); err != nil {
		log.Fatalf("failed to run up migrations: %v", err)
	}

	fmt.Println("Up migrations completed successfully")
}

func migrateDown(m *migrate.Migrate) {
	fmt.Println("Running down migrations...")
	if err := m.Down(); err != nil {
		log.Fatalf("failed to run down migrations: %v", err)
	}

	fmt.Println("Down migrations completed successfully")
}

// func execSeeders(db *sqlx.DB, seeders ...seeders.Seeder) {
// 	for _, seeder := range seeders {
// 		if err := seeder(db); err != nil {
// 			log.Fatalf("failed to run seeders: %v", err)
// 		}
// 	}
// }
