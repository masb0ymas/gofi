DROP INDEX IF EXISTS idx_uploads_id;
DROP INDEX IF EXISTS idx_uploads_created_at;
DROP INDEX IF EXISTS idx_uploads_updated_at;
DROP INDEX IF EXISTS idx_uploads_deleted_at;
DROP INDEX IF EXISTS idx_uploads_key_file;
DROP INDEX IF EXISTS idx_uploads_file_name;
DROP INDEX IF EXISTS idx_uploads_mimetype;
DROP INDEX IF EXISTS idx_uploads_signed_url;
DROP INDEX IF EXISTS idx_uploads_expires_at;

DROP TABLE IF EXISTS public."uploads";
