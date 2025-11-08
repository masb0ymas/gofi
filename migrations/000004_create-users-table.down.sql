DROP INDEX IF EXISTS idx_users_id;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_first_name;
DROP INDEX IF EXISTS idx_users_last_name;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_active_at;
DROP INDEX IF EXISTS idx_users_blocked_at;
DROP INDEX IF EXISTS idx_users_role_id;

DROP TABLE IF EXISTS public."users";
