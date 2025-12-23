-include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# MIGRATIONS
# ==================================================================================== #

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	./migrate --db-dsn=$(DB_DSN) --app-pwd=$(APP_DEFAULT_PASS) up

## db/migrations/up/seed: apply all up database migrations and run seeders
.PHONY: db/migrations/up/seed
db/migrations/up/seed:
	./migrate --db-dsn=$(DB_DSN) --app-pwd=$(APP_DEFAULT_PASS) --seed=dev up

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down:
	./migrate --db-dsn=$(DB_DSN) --app-pwd=$(APP_DEFAULT_PASS) down

## db/migrations/refresh: apply all refresh database migrations
.PHONY: db/migrations/refresh
db/migrations/refresh:
	./migrate --db-dsn=$(DB_DSN) --app-pwd=$(APP_DEFAULT_PASS) refresh

## db/migrations/refresh/seed: apply all refresh database migrations and run seeders
.PHONY: db/migrations/refresh/seed
db/migrations/refresh/seed:
	./migrate --db-dsn=$(DB_DSN) --app-pwd=$(APP_DEFAULT_PASS) --seed=dev refresh
