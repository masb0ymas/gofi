DROP INDEX IF EXISTS idx_refresh_tokens_id;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_created_at;
DROP INDEX IF EXISTS idx_refresh_tokens_revoked_at;

DROP TABLE IF EXISTS public."refresh_tokens";
