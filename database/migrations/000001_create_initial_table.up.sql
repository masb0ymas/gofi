CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT * FROM pg_timezone_names;
ALTER DATABASE "dev_dbgofi" SET timezone TO "Asia/Jakarta";

CREATE TABLE "upload" (
  "id" UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  "created_at" timestamp DEFAULT now(),
  "updated_at" timestamp DEFAULT now(),
  "deleted_at" timestamp,
  "key_file" varchar NOT NULL,
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
CREATE INDEX idx_upload_key_file ON "upload" (key_file);
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
	 ('03ba326e-f9ed-410a-818f-eaa409c13622',now(),now(),NULL,'Super Admin'),
	 ('9dc8b32b-aefe-44d3-bf19-6dc088d13174',now(),now(),NULL,'Admin'),
	 ('d7efa7e9-3c97-4217-a6bd-59e2eba53068',now(),now(),NULL,'User'),
	 ('be8482c9-7410-45eb-8c28-4dfd508a0de6',now(),now(),NULL,'Guest');

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
	 (uuid_generate_v4(),now(),now(),NULL,'Super Admin','super.admin@example.com','$argon2id$v=19$m=65536,t=3,p=2$hXwlaW+1NCwqKWDySLUk4g$ftx5ZLF5QjKLi50RW6qxPKZVDAPOvs6DxCY0L+GZz6A',NULL,NULL,true,false,'03ba326e-f9ed-410a-818f-eaa409c13622',NULL),
	 (uuid_generate_v4(),now(),now(),NULL,'Admin','admin@example.com','$argon2id$v=19$m=65536,t=3,p=2$ssShjR+1zMucGwSWI1p7rw$vTHTwnKQejOrxC4SlirCsJ7NfA1IC9pHonRAzBqKOUA',NULL,NULL,true,false,'9dc8b32b-aefe-44d3-bf19-6dc088d13174',NULL),
	 (uuid_generate_v4(),now(),now(),NULL,'User','user@example.com','$argon2id$v=19$m=65536,t=3,p=2$wnMuSBm5Fbw6mo5p4f3I6A$FzqhdZTYyklKziq506MM7cA2Cm7n4ud7GoSXMw6VVnc',NULL,NULL,true,false,'d7efa7e9-3c97-4217-a6bd-59e2eba53068',NULL);

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
