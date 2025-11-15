CREATE TABLE IF NOT EXISTS "uploads" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  "deleted_at" TIMESTAMP,
  "key_file" TEXT NOT NULL, -- format: /bucket/file_name
  "file_name" TEXT NOT NULL,
  "mimetype" VARCHAR NOT NULL, -- RFC 6838 standard (e.g., "application/pdf")
  "size" INTEGER NOT NULL, -- in bytes
  "signed_url" TEXT NOT NULL,
  "expires_at" TIMESTAMP NOT NULL -- ISO 8601 format
);

CREATE INDEX IF NOT EXISTS idx_uploads_id ON "uploads" ("id");
CREATE INDEX IF NOT EXISTS idx_uploads_created_at ON "uploads" ("created_at");
CREATE INDEX IF NOT EXISTS idx_uploads_updated_at ON "uploads" ("updated_at");
CREATE INDEX IF NOT EXISTS idx_uploads_deleted_at ON "uploads" ("deleted_at");
CREATE INDEX IF NOT EXISTS idx_uploads_key_file ON "uploads" ("key_file");
CREATE INDEX IF NOT EXISTS idx_uploads_file_name ON "uploads" ("file_name");
CREATE INDEX IF NOT EXISTS idx_uploads_mimetype ON "uploads" ("mimetype");
CREATE INDEX IF NOT EXISTS idx_uploads_signed_url ON "uploads" ("signed_url");
CREATE INDEX IF NOT EXISTS idx_uploads_expires_at ON "uploads" ("expires_at");
