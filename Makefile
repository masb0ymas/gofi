include .env

# build dir
BUILD_DIR=./cmd

# migration path
MIGRATION_PATH=./database/migrations

# database url
DATABASE_URL="$(DB_CONNECTION)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"

.PHONY: before-setup
before-setup:
	chmod +x setup.sh

.PHONY: setup
setup: before-setup
	bash setup.sh

.PHONY: update-deps
update-deps:
	go get -u && go mod tidy

.PHONY: dev
dev:
	./bin/air server --port $(APP_PORT)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

.PHONY: start
start: build
	$(BUILD_DIR)/$(APP_NAME)

# Using Golang Migrate
.PHONY: migration-create
migration-create:
	read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq "$$(date +%Y%m%d%H%M%S)_$$name"

# Using Golang Migrate
.PHONE: migration-up
migration-up:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose up

# Using Golang Migrate
.PHONE: migration-down
migration-down:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose down
