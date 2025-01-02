-- Drop tables
DROP TABLE IF EXISTS "account_provider";
DROP TABLE IF EXISTS "session";
DROP TABLE IF EXISTS "user";
DROP TABLE IF EXISTS "role";
DROP TABLE IF EXISTS "upload";

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Drop indexes
DROP INDEX IF EXISTS idx_account_provider_id;
DROP INDEX IF EXISTS idx_account_provider_created_at;
DROP INDEX IF EXISTS idx_account_provider_updated_at;
DROP INDEX IF EXISTS idx_account_provider_deleted_at;
DROP INDEX IF EXISTS idx_account_provider_user_id;
DROP INDEX IF EXISTS idx_account_provider_provider;

DROP INDEX IF EXISTS idx_session_id;
DROP INDEX IF EXISTS idx_session_created_at;
DROP INDEX IF EXISTS idx_session_updated_at;
DROP INDEX IF EXISTS idx_session_deleted_at;
DROP INDEX IF EXISTS idx_session_user_id;
DROP INDEX IF EXISTS idx_session_token;
DROP INDEX IF EXISTS idx_session_expires_at;

DROP INDEX IF EXISTS idx_role_id;
DROP INDEX IF EXISTS idx_role_created_at;
DROP INDEX IF EXISTS idx_role_updated_at;
DROP INDEX IF EXISTS idx_role_deleted_at;
DROP INDEX IF EXISTS idx_role_name;

DROP INDEX IF EXISTS idx_user_id;
DROP INDEX IF EXISTS idx_user_created_at;
DROP INDEX IF EXISTS idx_user_updated_at;
DROP INDEX IF EXISTS idx_user_deleted_at;
DROP INDEX IF EXISTS idx_user_fullname;
DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_user_token_verify;
DROP INDEX IF EXISTS idx_user_is_active;
DROP INDEX IF EXISTS idx_user_is_blocked;
DROP INDEX IF EXISTS idx_user_role_id;

DROP INDEX IF EXISTS idx_role_id;
DROP INDEX IF EXISTS idx_role_created_at;
DROP INDEX IF EXISTS idx_role_updated_at;
DROP INDEX IF EXISTS idx_role_deleted_at;
DROP INDEX IF EXISTS idx_role_name;

DROP INDEX IF EXISTS idx_upload_id;
DROP INDEX IF EXISTS idx_upload_created_at;
DROP INDEX IF EXISTS idx_upload_updated_at;
DROP INDEX IF EXISTS idx_upload_deleted_at;
DROP INDEX IF EXISTS idx_upload_key_file;
DROP INDEX IF EXISTS idx_upload_filename;
DROP INDEX IF EXISTS idx_upload_expires_at;
