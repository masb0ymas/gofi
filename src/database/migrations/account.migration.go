package migrations

func RoleMigration() string {
	schema := `
		CREATE TABLE IF NOT EXISTS "role" (
			id uuid default uuid_generate_v4(),
			created_at timestamp default now() NOT NULL,
			updated_at timestamp default now() NOT NULL,
			deleted_at timestamp,
			name varchar NOT NULL,
			PRIMARY KEY (id)
		);
	`

	return schema
}

func UserMigration() string {
	schema := `
		CREATE TABLE IF NOT EXISTS "user" (
			id uuid default uuid_generate_v4(),
			created_at timestamp default now() NOT NULL,
			updated_at timestamp default now() NOT NULL,
			deleted_at timestamp,
			fullname varchar NOT NULL,
			email varchar NOT NULL,
			password varchar,
			phone varchar(20),
			token_verify text,
			address text,
			is_active bool default false NOT NULL,
			is_blocked bool default false NOT NULL,
			role_id uuid NOT NULL,
			PRIMARY KEY (id),
			UNIQUE (email),
			CONSTRAINT fk_role 
				FOREIGN KEY (role_id) 
					REFERENCES "role"(id)
		);
	`

	return schema
}

func SessionMigration() string {
	schema := `
		CREATE TABLE IF NOT EXISTS "session" (
			id uuid default uuid_generate_v4(),
			created_at timestamp default now() NOT NULL,
			updated_at timestamp default now() NOT NULL,
			deleted_at timestamp,
			user_id uuid NOT NULL,
			token text NOT NULL,
			ipAddress varchar,
			device varchar,
			platform varchar,
			latitude varchar,
			longitude varchar,
			PRIMARY KEY (id),
			CONSTRAINT fk_user 
				FOREIGN KEY (user_id) 
					REFERENCES "user"(id)
		);
	`

	return schema
}
