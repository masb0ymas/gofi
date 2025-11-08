CREATE TABLE IF NOT EXISTS "user_verify_accounts" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(), -- Polymorphic ID (User ID)
  "token" TEXT NOT NULL,
  "expires_at" TIMESTAMP NOT NULL
);

-- Polymorphic table
ALTER TABLE "user_verify_accounts" ADD FOREIGN KEY ("id") REFERENCES "users" ("id") ON DELETE CASCADE;
