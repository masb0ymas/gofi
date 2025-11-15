CREATE TABLE IF NOT EXISTS "roles" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
  "created_at" TIMESTAMP DEFAULT now(),
  "updated_at" TIMESTAMP DEFAULT now(),
  "deleted_at" TIMESTAMP,
  "name" VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_roles_id ON "roles" ("id");
CREATE INDEX IF NOT EXISTS idx_roles_created_at ON "roles" ("created_at");
CREATE INDEX IF NOT EXISTS idx_roles_updated_at ON "roles" ("updated_at");
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON "roles" ("deleted_at");
CREATE INDEX IF NOT EXISTS idx_roles_name ON "roles" ("name");
