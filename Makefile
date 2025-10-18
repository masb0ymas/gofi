include .env

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# command air
CMD_AIR=./bin/air # for linux / macos
# CMD_AIR=./bin/air.exe # for windows

# command migrate
# CMD_MIGRATE=./bin/migrate # local golang-migrate
CMD_MIGRATE=migrate # global golang-migrate

# build dir
BUILD_DIR=./cmd

# migration path
MIGRATION_PATH=./src/app/database/migrations

# database url
DATABASE_URL="$(DB_CONNECTION)://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## setup: setup the application
.PHONY: setup
setup:
	chmod +x ./script/setup.sh
	bash ./script/setup.sh

## release: release the application
.PHONY: release
release:
	chmod +x ./script/release.sh
	bash ./script/release.sh

## release/bump: bump the version
.PHONY: release/bump
release/bump:
	chmod +x ./script/bump-version.sh
	bash ./script/bump-version.sh

## deps/update: update the dependencies
.PHONY: deps/update
deps/update:
	go get -u && go mod tidy

## lint: lint the code & check typing
.PHONY: lint
lint:
	golangci-lint run

## dev: run the application in development mode
.PHONY: dev
dev:
	$(CMD_AIR) server --port $(APP_PORT)

## clean: clean the build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

## build: build the application
.PHONY: build
build: clean
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

## start: start the application
.PHONY: start
start: build
	$(BUILD_DIR)/$(APP_NAME)

## db/migration/version: print the migration version
.PHONY: db/migration/version
db/migration/version:
	$(CMD_MIGRATE) --version

## db/migration/create: create a new migration
.PHONY: db/migration/create
db/migration/create:
	read -p "Enter migration name: " name; \
	$(CMD_MIGRATE) create -ext sql -dir $(MIGRATION_PATH) -seq "$$(date +%Y%m%d%H%M%S)_$$name"

## db/migration/up: run the migration up
.PHONY: db/migration/up
db/migration/up:
	$(CMD_MIGRATE) -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose up

## db/migration/down: run the migration down
.PHONY: db/migration/down
db/migration/down:
	$(CMD_MIGRATE) -path $(MIGRATION_PATH) -database $(DATABASE_URL) -verbose down

.PHONY: db/migration/clear-dirty-flag
db/migration/clear-dirty-flag:
	@read -p "Version to clear flag: " DIRTY_VERSION; \
	$(CMD_MIGRATE) -path $(MIGRATION_PATH) -database $(DATABASE_URL) force $$DIRTY_VERSION

.PHONY: db/migration/migration-stage
db/migration/migration-stage:
	$(CMD_MIGRATE) -path $(MIGRATION_PATH) -database $(DATABASE_URL) version

.PHONY: db/migration/force-version
db/migration/force-version:
	@read -p "Version to force to: " NEW_VERSION; \
	$(CMD_MIGRATE) -path $(MIGRATION_PATH) -database $(DATABASE_URL) force $$NEW_VERSION

## tag/staging: deploy staging by creating a new tag with incremented patch version (run with 'push' argument to automatically push the tag)
.PHONY: tag/staging
tag/staging:
	@echo "Creating new staging tag with incremented patch version..."
	@latest_tag=$(git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-staging | tail -n1); \
	if [ -z "$latest_tag" ]; then \
		echo "No existing staging tag found. Creating first staging tag v0.0.1-staging"; \
		new_tag="v0.0.1-staging"; \
	else \
		echo "Found latest staging tag: $latest_tag"; \
		base_version=$(echo $latest_tag | sed 's/-staging$//'); \
		major=$(echo $base_version | cut -d. -f1); \
		minor=$(echo $base_version | cut -d. -f2); \
		patch=$(echo $base_version | cut -d. -f3); \
		new_patch=$((patch + 1)); \
		new_tag="v$major.$minor.$new_patch-staging"; \
	fi; \
	echo "Creating new tag: $new_tag"; \
	git tag $new_tag; \
	if [ "$(push)" = "push" ]; then \
		echo "Pushing tag $new_tag to origin..."; \
		git push origin $new_tag; \
		echo "Tag $new_tag has been created and pushed to origin"; \
	else \
		echo "New staging tag $new_tag created successfully. Run 'git push origin $new_tag' to push the tag to remote."; \
	fi

## tag/release: deploy production by creating a new tag with incremented patch version (run with 'push' argument to automatically push the tag)
.PHONY: tag/release
tag/release:
	@echo "Creating new production tag with incremented patch version..."
	@latest_tag=$(git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-release | tail -n1); \
	if [ -z "$latest_tag" ]; then \
		echo "No existing production tag found. Creating first production tag v0.0.1-release"; \
		new_tag="v0.0.1-release"; \
	else \
		echo "Found latest production tag: $latest_tag"; \
		base_version=$(echo $latest_tag | sed 's/-release$//'); \
		major=$(echo $base_version | cut -d. -f1); \
		minor=$(echo $base_version | cut -d. -f2); \
		patch=$(echo $base_version | cut -d. -f3); \
		new_patch=$((patch + 1)); \
		new_tag="v$major.$minor.$new_patch-release"; \
	fi; \
	echo "Creating new tag: $new_tag"; \
	git tag $new_tag; \
	if [ "$(push)" = "push" ]; then \
		echo "Pushing tag $new_tag to origin..."; \
		git push origin $new_tag; \
		echo "Tag $new_tag has been created and pushed to origin"; \
	else \
		echo "New production tag $new_tag created successfully. Run 'git push origin $new_tag' to push the tag to remote."; \
	fi
