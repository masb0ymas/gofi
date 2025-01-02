-- Create extension for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT * FROM pg_timezone_names;
ALTER DATABASE dev_dbgoarif SET timezone TO 'Asia/Jakarta';

-- Create upload table
CREATE TABLE "upload" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    key_file TEXT NOT NULL,
    filename VARCHAR(255) NOT NULL,
    mimetype VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    signed_url TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

-- Create role table
CREATE TABLE "role" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL
);

-- Create user table
CREATE TABLE "user" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    fullname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255),
    phone VARCHAR(50),
    token_verify TEXT,
    is_active BOOLEAN DEFAULT FALSE NOT NULL,
    is_blocked BOOLEAN DEFAULT FALSE NOT NULL,
    role_id UUID NOT NULL,
    upload_id UUID,
    FOREIGN KEY (role_id) REFERENCES "role" (id),
    FOREIGN KEY (upload_id) REFERENCES "upload" (id)
);

-- Create session table
CREATE TABLE "session" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address VARCHAR(255) NOT NULL,
    user_agent TEXT NOT NULL,
    latitude VARCHAR(255),
    longitude VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- Create account_provider table
CREATE TABLE "account_provider" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP,
    user_id UUID NOT NULL,
    provider VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    id_token TEXT,
    FOREIGN KEY (user_id) REFERENCES "user" (id)
);

-- Create indexes for better query performance
CREATE INDEX idx_upload_id ON "upload"(id);
CREATE INDEX idx_upload_created_at ON "upload"(created_at);
CREATE INDEX idx_upload_updated_at ON "upload"(updated_at);
CREATE INDEX idx_upload_deleted_at ON "upload"(deleted_at);
CREATE INDEX idx_upload_key_file ON "upload"(key_file);
CREATE INDEX idx_upload_filename ON "upload"(filename);
CREATE INDEX idx_upload_expires_at ON "upload"(expires_at);

CREATE INDEX idx_role_id ON "role"(id);
CREATE INDEX idx_role_created_at ON "role"(created_at);
CREATE INDEX idx_role_updated_at ON "role"(updated_at);
CREATE INDEX idx_role_deleted_at ON "role"(deleted_at);
CREATE INDEX idx_role_name ON "role"(name);

CREATE INDEX idx_user_id ON "user"(id);
CREATE INDEX idx_user_created_at ON "user"(created_at);
CREATE INDEX idx_user_updated_at ON "user"(updated_at);
CREATE INDEX idx_user_deleted_at ON "user"(deleted_at);
CREATE INDEX idx_user_fullname ON "user"(fullname);
CREATE INDEX idx_user_email ON "user"(email);
CREATE INDEX idx_user_token_verify ON "user"(token_verify);
CREATE INDEX idx_user_is_active ON "user"(is_active);
CREATE INDEX idx_user_is_blocked ON "user"(is_blocked);
CREATE INDEX idx_user_role_id ON "user"(role_id);

CREATE INDEX idx_session_id ON "session"(id);
CREATE INDEX idx_session_created_at ON "session"(created_at);
CREATE INDEX idx_session_updated_at ON "session"(updated_at);
CREATE INDEX idx_session_deleted_at ON "session"(deleted_at);
CREATE INDEX idx_session_user_id ON "session"(user_id);
CREATE INDEX idx_session_token ON "session"(token);
CREATE INDEX idx_session_expires_at ON "session"(expires_at);

CREATE INDEX idx_account_provider_id ON "account_provider"(id);
CREATE INDEX idx_account_provider_created_at ON "account_provider"(created_at);
CREATE INDEX idx_account_provider_updated_at ON "account_provider"(updated_at);
CREATE INDEX idx_account_provider_deleted_at ON "account_provider"(deleted_at);
CREATE INDEX idx_account_provider_user_id ON "account_provider"(user_id);
CREATE INDEX idx_account_provider_provider ON "account_provider"(provider);

-- Insert default role
INSERT INTO "role" (id, created_at, updated_at, name) VALUES 
  ('8d4aa7ad-f5aa-45da-84f0-146fb34b59ef', now(), now(), 'Super Admin'),
  ('b11fe5f0-63bd-4735-90dd-57171de19904', now(), now(), 'Admin'),
  ('cd85009f-c550-4529-bc04-adacfe1cd9a0', now(), now(), 'User');

-- Insert default user
INSERT INTO "user" (id, created_at, updated_at, fullname, email, password, is_active, role_id) VALUES 
  ('798b7d7f-c3d2-4113-a3ef-de003ec798b6', now(), now(), 'Super Admin', 'super.admin@mail.com', '$argon2id$v=19$m=65536,t=3,p=2$hXwlaW+1NCwqKWDySLUk4g$ftx5ZLF5QjKLi50RW6qxPKZVDAPOvs6DxCY0L+GZz6A', TRUE, '8d4aa7ad-f5aa-45da-84f0-146fb34b59ef'),
  ('0ac68f78-9060-45fd-be35-9d28578716ec', now(), now(), 'Admin', 'admin@mail.com', '$argon2id$v=19$m=65536,t=3,p=2$ssShjR+1zMucGwSWI1p7rw$vTHTwnKQejOrxC4SlirCsJ7NfA1IC9pHonRAzBqKOUA', TRUE, 'b11fe5f0-63bd-4735-90dd-57171de19904'),
  ('f89504e2-d98a-4f68-afe3-6e2446a5c3e8', now(), now(), 'User', 'user@mail.com', '$argon2id$v=19$m=65536,t=3,p=2$wnMuSBm5Fbw6mo5p4f3I6A$FzqhdZTYyklKziq506MM7cA2Cm7n4ud7GoSXMw6VVnc', TRUE, 'cd85009f-c550-4529-bc04-adacfe1cd9a0');
