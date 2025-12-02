CREATE TABLE IF NOT EXISTS "users_oauth" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "user_id" UUID NOT NULL,
  "provider" VARCHAR NOT NULL,
  "access_token" TEXT NOT NULL,
  "refresh_token" TEXT,
  "expires_at" TIMESTAMP NOT NULL,
);

CREATE INDEX IF NOT EXISTS idx_users_oauth_id ON "users_oauth" ("id");
CREATE INDEX IF NOT EXISTS idx_users_oauth_user_id ON "users_oauth" ("user_id");
CREATE INDEX IF NOT EXISTS idx_users_oauth_provider ON "users_oauth" ("provider");
CREATE INDEX IF NOT EXISTS idx_users_oauth_access_token ON "users_oauth" ("access_token");
CREATE INDEX IF NOT EXISTS idx_users_oauth_expires_at ON "users_oauth" ("expires_at");

ALTER TABLE "users_oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
