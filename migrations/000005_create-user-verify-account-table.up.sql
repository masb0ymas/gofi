CREATE TABLE IF NOT EXISTS "user_verify_accounts" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(), -- Polymorphic ID (User ID)
  "token" TEXT NOT NULL,
  "expires_at" TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_verify_accounts_id ON "user_verify_accounts" ("id");
CREATE INDEX IF NOT EXISTS idx_user_verify_accounts_token ON "user_verify_accounts" ("token");
CREATE INDEX IF NOT EXISTS idx_user_verify_accounts_expires_at ON "user_verify_accounts" ("expires_at");

-- Polymorphic table
ALTER TABLE "user_verify_accounts" ADD FOREIGN KEY ("id") REFERENCES "users" ("id") ON DELETE CASCADE;
