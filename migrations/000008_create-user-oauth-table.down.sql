DROP INDEX IF EXISTS idx_users_oauth_id;
DROP INDEX IF EXISTS idx_users_oauth_user_id;
DROP INDEX IF EXISTS idx_users_oauth_provider;
DROP INDEX IF EXISTS idx_users_oauth_access_token;
DROP INDEX IF EXISTS idx_users_oauth_expires_at;

DROP TABLE IF EXISTS public."users_oauth";
