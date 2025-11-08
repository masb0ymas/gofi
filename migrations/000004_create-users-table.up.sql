CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  "deleted_at" TIMESTAMP,
  "first_name" VARCHAR NOT NULL,
  "last_name" VARCHAR,
  "email" VARCHAR NOT NULL UNIQUE,
  "password" TEXT,
  "phone" VARCHAR(20),
  "active_at" TIMESTAMP, -- ISO 8601 format
  "blocked_at" TIMESTAMP, -- ISO 8601 format
  "role_id" UUID NOT NULL,
  "upload_id" UUID
);

CREATE INDEX IF NOT EXISTS idx_users_id ON "users" ("id");
CREATE INDEX IF NOT EXISTS idx_users_created_at ON "users" ("created_at");
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON "users" ("updated_at");
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON "users" ("deleted_at");
CREATE INDEX IF NOT EXISTS idx_users_first_name ON "users" ("first_name");
CREATE INDEX IF NOT EXISTS idx_users_last_name ON "users" ("last_name");
CREATE INDEX IF NOT EXISTS idx_users_email ON "users" ("email");
CREATE INDEX IF NOT EXISTS idx_users_active_at ON "users" ("active_at");
CREATE INDEX IF NOT EXISTS idx_users_blocked_at ON "users" ("blocked_at");
CREATE INDEX IF NOT EXISTS idx_users_role_id ON "users" ("role_id");

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON DELETE CASCADE;
ALTER TABLE "users" ADD FOREIGN KEY ("upload_id") REFERENCES "uploads" ("id") ON DELETE CASCADE;
