DROP INDEX IF EXISTS idx_user_oauths_id;
DROP INDEX IF EXISTS idx_user_oauths_user_id;
DROP INDEX IF EXISTS idx_user_oauths_provider;
DROP INDEX IF EXISTS idx_user_oauths_access_token;
DROP INDEX IF EXISTS idx_user_oauths_expires_at;

DROP TABLE IF EXISTS public."user_oauths";
