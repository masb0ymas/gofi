CREATE TABLE IF NOT EXISTS "refresh_tokens" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "user_id" UUID NOT NULL,
  "token" TEXT NOT NULL,
  "expires_at" TIMESTAMP NOT NULL,
  "created_at" TIMESTAMP DEFAULT now(),
  "revoked_at" TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sessions_id ON "refresh_tokens" ("id");
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON "refresh_tokens" ("user_id");
CREATE INDEX IF NOT EXISTS idx_sessions_token ON "refresh_tokens" ("token");
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON "refresh_tokens" ("expires_at");
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON "refresh_tokens" ("created_at");
CREATE INDEX IF NOT EXISTS idx_sessions_revoked_at ON "refresh_tokens" ("revoked_at");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
