CREATE TABLE IF NOT EXISTS "user_oauths" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "user_id" UUID NOT NULL,
  "provider" VARCHAR NOT NULL,
  "access_token" TEXT NOT NULL,
  "refresh_token" TEXT,
  "expires_at" TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_user_oauths_id ON "user_oauths" ("id");
CREATE INDEX IF NOT EXISTS idx_user_oauths_user_id ON "user_oauths" ("user_id");
CREATE INDEX IF NOT EXISTS idx_user_oauths_provider ON "user_oauths" ("provider");
CREATE INDEX IF NOT EXISTS idx_user_oauths_access_token ON "user_oauths" ("access_token");
CREATE INDEX IF NOT EXISTS idx_user_oauths_expires_at ON "user_oauths" ("expires_at");

ALTER TABLE "user_oauths" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
