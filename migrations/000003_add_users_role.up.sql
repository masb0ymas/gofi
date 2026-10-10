-- Application role of a mirrored user. Authula knows nothing about it, so the
-- mirror upsert never touches this column.
ALTER TABLE "users" ADD COLUMN IF NOT EXISTS "role" VARCHAR NOT NULL DEFAULT 'user';
