-include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# APPLICATION
# ==================================================================================== #

## run: run the application
.PHONY: run
run:
	go run ./cmd/api \
		--machine-id=$(MACHINE_ID) \
		--debug=$(DEBUG) \
		--env=$(ENV) \
		--port=$(PORT) \
		--app-name=$(APP_NAME) \
		--app-secret=$(APP_SECRET) \
		--cors-allowed-origins=$(CORS_ALLOWED_ORIGINS) \
		--db-dsn=$(DB_DSN) \
		--db-max-open-conns=$(DB_MAX_OPEN_CONNS) \
		--db-max-idle-conns=$(DB_MAX_IDLE_CONNS) \
		--db-max-idle-time=$(DB_MAX_IDLE_TIME) \
		--migrate-on-boot=$(MIGRATE_ON_BOOT) \
		--seed-on-boot=$(SEED_ON_BOOT) \
		--auth-base-url=$(AUTH_BASE_URL) \
		--email-from=$(FROM_ADDRESS) \
		--resend-api-key=$(RESEND_API_KEY) \
		--google-client-id=$(GOOGLE_CLIENT_ID) \
		--google-client-secret=$(GOOGLE_CLIENT_SECRET)

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for $(name)...'
	migrate create -seq -ext=.sql -dir=./migrations $(name)

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@go run ./cmd/migrate --db-dsn=$(DB_DSN) up

## db/migrations/up/seed: apply all up database migrations and run seeders
.PHONY: db/migrations/up/seed
db/migrations/up/seed:
	@go run ./cmd/migrate --db-dsn=$(DB_DSN) --seed=dev up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	@go run ./cmd/migrate --db-dsn=$(DB_DSN) down

## db/migrations/refresh: apply all refresh database migrations
.PHONY: db/migrations/refresh
db/migrations/refresh:
	@go run ./cmd/migrate --db-dsn=$(DB_DSN) refresh

## db/migrations/refresh/seed: apply all refresh database migrations and run seeders
.PHONY: db/migrations/refresh/seed
db/migrations/refresh/seed:
	@go run ./cmd/migrate --db-dsn=$(DB_DSN) --seed=dev refresh

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags="-s" -o=./bin/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/api ./cmd/api

## build/migrate: build the cmd/migrate application
.PHONY: build/migrate
build/migrate:
	@echo 'Building cmd/migrate...'
	go build -ldflags="-s" -o=./bin/migrate ./cmd/migrate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/linux_amd64/migrate ./cmd/migrate

# ==================================================================================== #
# RELEASE
# ==================================================================================== #

## release: bump version and create a release tag (auto-detect from commits); extras: make release ARGS="--type=minor --no-push"
.PHONY: release
release:
	@bash scripts/release.sh $(ARGS)

## release/dry-run: preview the next version bump without creating a tag
.PHONY: release/dry-run
release/dry-run:
	@bash scripts/release.sh --dry-run $(ARGS)