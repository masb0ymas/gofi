# GoFi — Golang Fiber + sqlx starter
  
Modern Go starter using Fiber v2 and sqlx with a pragmatic developer workflow.

## Features
- **Go 1.24**
- **Fiber v2** HTTP framework
- **sqlx** with **PostgreSQL** (`lib/pq`)
- **Environment management** via `godotenv`
- **JWT utilities** (`github.com/golang-jwt/jwt/v5`)
- **Error tracking** with Sentry (optional)
- **Google Cloud Storage** helper (optional)
- **Live reload** with Air (`.air.toml`, `./bin/air`)
- **Database migrations** with golang-migrate
- **Dockerfile** for production builds
- **Makefile** with common tasks

## Project structure
```
.
├─ .air.toml
├─ .env.example
├─ Dockerfile
├─ Makefile
├─ public/
│  ├─ docs/
│  └─ email-template/
├─ script/
│  ├─ setup.sh
│  ├─ release.sh
│  └─ bump-version.sh
├─ src/
│  ├─ app/
│  │  └─ database/
│  │     └─ migrations/        # migration .sql files
│  ├─ config/
│  ├─ lib/
│  └─ main.go                   # application entrypoint
├─ bin/                         # local tools (air, migrate, etc.)
├─ go.mod
└─ go.sum
```

## Prerequisites
- Go 1.24+
- PostgreSQL database
- golang-migrate
  - Global: `brew install golang-migrate` (macOS), or see https://github.com/golang-migrate/migrate
  - Or adjust `CMD_MIGRATE` in `Makefile` to point to your local binary

## Quick start
1. Clone the repository
2. Copy environment file
   ```bash
   cp .env.example .env
   ```
3. Edit `.env` and set at least:
   - `APP_NAME`, `APP_ENV`, `APP_PORT`
   - `DB_CONNECTION`, `DB_HOST`, `DB_PORT`, `DB_USERNAME`, `DB_PASSWORD`, `DB_DATABASE`
4. Create your database in PostgreSQL
5. Run database migrations
   ```bash
   make db/migration/up
   ```
6. Start in development mode (hot-reload)
   ```bash
   make dev
   ```
   App runs on `http://localhost:${APP_PORT}`

## Common tasks (Makefile)
- **help**: show available commands
  ```bash
  make help
  ```
- **deps/update**: upgrade deps and tidy
  ```bash
  make deps/update
  ```
- **db/migration/create**: create a new migration
  ```bash
  make db/migration/create
  # then enter a name, files appear in src/app/database/migrations
  ```
- **db/migration/up / down**: apply or rollback migrations
  ```bash
  make db/migration/up
  make db/migration/down
  ```
- **build**: produce a binary into `./cmd/${APP_NAME}`
  ```bash
  make build
  ```
- **start**: build and run the binary
  ```bash
  make start
  ```
- **setup**, **release**, **release/bump**: utility scripts under `script/`
- **tag/staging**, **tag/release**: create version tags for deployments

## Configuration
- Local development reads `.env` (via `godotenv`).
- Migrations use `DATABASE_URL` derived from the `DB_*` variables in `Makefile`.
- Air config is in `.air.toml` and uses `./bin/air` by default.

## Docker
Build a minimal production image and run it:
```bash
docker build -t gofi:local .
docker run --env-file .env -p 8000:8000 gofi:local
```
Notes:
- The container exposes port `8000` by default.
- Ensure your env vars are suitable for the container environment (e.g., DB host is reachable from the container).

## Testing & linting
- This template doesn’t ship opinionated test/lint tools by default. Add your preferred stack (e.g., `go test ./...`, `golangci-lint`) as needed.

## License
This project is licensed under the terms of the MIT License. See `LICENSE.md`.
