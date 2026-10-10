DROP TABLE IF EXISTS "users";
-- Dropped with the mirror so the two never disagree; Authula recreates its
-- tables on the next boot.
DROP SCHEMA IF EXISTS "authula" CASCADE;
