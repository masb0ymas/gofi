CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT * FROM pg_timezone_names;
ALTER DATABASE "dev_dbgofi" SET timezone TO "Asia/Jakarta";

CREATE TABLE "upload" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now(),
  "deleted_at" timestamp,
  "keyfile" varchar NOT NULL,
  "filename" text NOT NULL,
  "mimetype" varchar NOT NULL,
  "size" int NOT NULL,
  "signed_url" text NOT NULL,
  "expired_at" timestamp NOT NULL
);

CREATE INDEX idx_upload_id ON "upload" (id);
CREATE INDEX idx_upload_created_at ON "upload" (created_at);
CREATE INDEX idx_upload_updated_at ON "upload" (updated_at);
CREATE INDEX idx_upload_deleted_at ON "upload" (deleted_at);
CREATE INDEX idx_upload_keyfile ON "upload" (keyfile);
CREATE INDEX idx_upload_filename ON "upload" (filename);
CREATE INDEX idx_upload_expired_at ON "upload" (expired_at);

CREATE TABLE "role" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now(),
  "deleted_at" timestamp,
  "name" varchar NOT NULL
);

CREATE INDEX idx_role_id ON "role" (id);
CREATE INDEX idx_role_created_at ON "role" (created_at);
CREATE INDEX idx_role_updated_at ON "role" (updated_at);
CREATE INDEX idx_role_deleted_at ON "role" (deleted_at);
CREATE INDEX idx_role_name ON "role" (name);

INSERT INTO "role" ("id","created_at","updated_at","deleted_at","name") VALUES
	 ('d6547c6b-16a8-4e84-9792-3eb6e7c35d35',now(),now(),NULL,'Super Admin'),
	 ('560e4ac4-09cc-4a63-91dc-2bfce03bf9e6',now(),now(),NULL,'Admin'),
	 ('fb81445d-0190-499e-b3af-9c8a4522b8e1',now(),now(),NULL,'User');

CREATE TABLE "user" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now(),
  "deleted_at" timestamp,
  "fullname" varchar NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" text NOT NULL,
  "phone" varchar(20) NULL,
  "token_verify" text NULL,
  "is_active" bool NOT NULL,
  "is_blocked" bool NOT NULL,
  "role_id" uuid NOT NULL,
  "upload_id" uuid NULL
);

CREATE INDEX idx_user_id ON "user" (id);
CREATE INDEX idx_user_created_at ON "user" (created_at);
CREATE INDEX idx_user_updated_at ON "user" (updated_at);
CREATE INDEX idx_user_deleted_at ON "user" (deleted_at);
CREATE INDEX idx_user_fullname ON "user" (fullname);
CREATE INDEX idx_user_email ON "user" (email);
CREATE INDEX idx_user_token_verify ON "user" (token_verify);
CREATE INDEX idx_user_is_active ON "user" (is_active);
CREATE INDEX idx_user_is_blocked ON "user" (is_blocked);
CREATE INDEX idx_user_role_id ON "user" (role_id);

ALTER TABLE "user" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");
ALTER TABLE "user" ADD FOREIGN KEY ("upload_id") REFERENCES "upload" ("id");

INSERT INTO "user" ("id","created_at","updated_at","deleted_at","fullname","email","password","phone","token_verify","is_active","is_blocked","role_id","upload_id") VALUES
	 (uuid_generate_v4(),now(),now(),NULL,'Super Admin','super.admin@example.com','$argon2id$v=19$m=65536,t=3,p=2$hXwlaW+1NCwqKWDySLUk4g$ftx5ZLF5QjKLi50RW6qxPKZVDAPOvs6DxCY0L+GZz6A',NULL,NULL,true,false,'d6547c6b-16a8-4e84-9792-3eb6e7c35d35',NULL),
	 (uuid_generate_v4(),now(),now(),NULL,'Admin','admin@example.com','$argon2id$v=19$m=65536,t=3,p=2$ssShjR+1zMucGwSWI1p7rw$vTHTwnKQejOrxC4SlirCsJ7NfA1IC9pHonRAzBqKOUA',NULL,NULL,true,false,'560e4ac4-09cc-4a63-91dc-2bfce03bf9e6',NULL),
	 (uuid_generate_v4(),now(),now(),NULL,'User','user@example.com','$argon2id$v=19$m=65536,t=3,p=2$wnMuSBm5Fbw6mo5p4f3I6A$FzqhdZTYyklKziq506MM7cA2Cm7n4ud7GoSXMw6VVnc',NULL,NULL,true,false,'fb81445d-0190-499e-b3af-9c8a4522b8e1',NULL);

CREATE TABLE "session" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now(),
  "user_id" uuid NOT NULL,
  "token" text NOT NULL,
  "expired_at" timestamp NOT NULL
);

CREATE INDEX idx_session_id ON "session" (id);
CREATE INDEX idx_session_created_at ON "session" (created_at);
CREATE INDEX idx_session_updated_at ON "session" (updated_at);
CREATE INDEX idx_session_user_id ON "session" (user_id);
CREATE INDEX idx_session_token ON "session" (token);
CREATE INDEX idx_session_expired_at ON "session" (expired_at);

ALTER TABLE "session" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");
