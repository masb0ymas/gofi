CREATE TABLE IF NOT EXISTS "categories" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "deleted_at" TIMESTAMPTZ,
  "name" VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_categories_created_at ON "categories" ("created_at");
CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON "categories" ("deleted_at");
CREATE INDEX IF NOT EXISTS idx_categories_name ON "categories" ("name");