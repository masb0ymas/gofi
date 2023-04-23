package migrations

import "os"

func BaseMigration() string {
	DB_DATABASE := os.Getenv("DB_DATABASE")
	DB_TIMEZONE := os.Getenv("DB_TIMEZONE")

	schema := `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		SELECT * FROM pg_timezone_names;

		ALTER DATABASE ` + DB_DATABASE + ` SET timezone TO '` + DB_TIMEZONE + `';

		CREATE TABLE IF NOT EXISTS "migrations" (
			id uuid default uuid_generate_v4(),
			created_at timestamp default now() NOT NULL,
			updated_at timestamp default now() NOT NULL,
			name varchar NOT NULL,
			PRIMARY KEY (id)
		);
	`

	return schema
}
