CREATE TABLE IF NOT EXISTS "sessions" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  "user_id" UUID NOT NULL,
  "token" TEXT NOT NULL,
  "expires_at" TIMESTAMP NOT NULL,
  "ip_address" VARCHAR NOT NULL,
  "user_agent" TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_id ON "sessions" ("id");
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON "sessions" ("created_at");
CREATE INDEX IF NOT EXISTS idx_sessions_updated_at ON "sessions" ("updated_at");
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON "sessions" ("user_id");
CREATE INDEX IF NOT EXISTS idx_sessions_token ON "sessions" ("token");
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON "sessions" ("expires_at");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
