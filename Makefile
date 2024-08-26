include .env

# build dir
BUILD_DIR=./dist

# migration path
MIGRATION_PATH=./database/migrations

# database url
DATABASE_URL="$(DB_CONNECTION)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"

.PHONY: update-deps
update-deps:
	go get -u && go mod tidy

.PHONY: dev
dev:
	./bin/air server --port $(APP_PORT)

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: build
build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

.PHONY: start
start: build
	$(BUILD_DIR)/$(APP_NAME)

.PHONY: migration-create
migration-create: 
	migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_initial_table

.PHONE: migration-up
migration-up:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose up

.PHONE: migration-down
migration-down:
	migrate -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose down
