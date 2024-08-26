DROP INDEX IF EXISTS idx_session_id, idx_session_created_at, idx_session_updated_at, idx_session_user_id, idx_session_token, idx_session_expired_at;
DROP INDEX IF EXISTS idx_user_id, idx_user_created_at, idx_user_updated_at, idx_user_deleted_at, idx_user_fullname, idx_user_email, idx_user_token_verify, idx_user_is_active, idx_user_is_blocked, idx_user_role_id;
DROP INDEX IF EXISTS idx_role_id, idx_role_created_at, idx_role_updated_at, idx_role_deleted_at, idx_role_name;
DROP INDEX IF EXISTS idx_upload_id, idx_upload_created_at, idx_upload_updated_at, idx_upload_deleted_at, idx_upload_key_file, idx_upload_filename, idx_upload_expired_at;

DROP TABLE IF EXISTS public."session";
DROP TABLE IF EXISTS public."user";
DROP TABLE IF EXISTS public."role";
DROP TABLE IF EXISTS public."upload";

DROP EXTENSION IF EXISTS "uuid-ossp";
