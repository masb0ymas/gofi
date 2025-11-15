CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

SELECT * FROM pg_timezone_names;
ALTER DATABASE "dev_gofi" SET timezone TO "Asia/Jakarta";
