-- Authula keeps its own tables (users, accounts, sessions, verifications) in
-- this schema; it creates them when the API boots.
CREATE SCHEMA IF NOT EXISTS "authula";

-- Mirror of "authula"."users", kept in sync by the API so application tables
-- can reference users without reaching into the auth schema.
CREATE TABLE IF NOT EXISTS "users" (
  "id" UUID PRIMARY KEY NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "name" VARCHAR NOT NULL,
  "email" VARCHAR NOT NULL,
  "image" TEXT
);

CREATE INDEX IF NOT EXISTS idx_users_created_at ON "users" ("created_at");
CREATE INDEX IF NOT EXISTS idx_users_email ON "users" ("email");
