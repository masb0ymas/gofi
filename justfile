# Load environment variables from .env file
set dotenv-load

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

# Command binaries
CMD_AIR := "./bin/air"              # for linux / macos
# CMD_AIR := "./bin/air.exe"        # for windows

# CMD_MIGRATE := "./bin/migrate"    # local golang-migrate
CMD_MIGRATE := "migrate"            # global golang-migrate

CMD_GOSEC := "./bin/gosec"          # local gosec
# CMD_GOSEC := "gosec"              # global gosec

# Build directory
BUILD_DIR := "./cmd"

# Migration path
MIGRATION_PATH := "./src/app/database/migrations"

# Database URL
DATABASE_URL := env_var('DB_CONNECTION') + "://" + env_var('DB_USERNAME') + ":" + env_var('DB_PASSWORD') + "@" + env_var('DB_HOST') + ":" + env_var('DB_PORT') + "/" + env_var('DB_DATABASE') + "?sslmode=disable"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

# help: print this help message
help:
    @just --list

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

# setup: setup the application
setup:
    chmod +x ./script/setup.sh
    bash ./script/setup.sh

# release: release the application
release:
    chmod +x ./script/release.sh
    bash ./script/release.sh

# release-bump: bump the version
release-bump:
    chmod +x ./script/bump-version.sh
    bash ./script/bump-version.sh

# deps-update: update the dependencies
deps-update:
    go get -u && go mod tidy

# lint: lint the code & check typing
lint:
    golangci-lint run

# security: scan the code for security vulnerabilities
security:
    {{CMD_GOSEC}} ./...

# dev: run the application in development mode
dev:
    {{CMD_AIR}} server --port {{env_var('APP_PORT')}}

# clean: clean the build directory
clean:
    rm -rf {{BUILD_DIR}}

# build: build the application
build: clean
    CGO_ENABLED=0 go build -ldflags="-w -s" -o {{BUILD_DIR}}/{{env_var('APP_NAME')}} ./src/main.go

# start: start the application
start: build
    {{BUILD_DIR}}/{{env_var('APP_NAME')}}

# db-migration-version: print the migration version
db-migration-version:
    {{CMD_MIGRATE}} --version

# db-migration-create: create a new migration
db-migration-create:
    #!/usr/bin/env bash
    read -p "Enter migration name: " name
    {{CMD_MIGRATE}} create -ext sql -dir {{MIGRATION_PATH}} -seq "$(date +%Y%m%d%H%M%S)_$name"

# db-migration-up: run the migration up
db-migration-up:
    {{CMD_MIGRATE}} -path {{MIGRATION_PATH}} -database "{{DATABASE_URL}}" -verbose up

# db-migration-down: run the migration down
db-migration-down:
    {{CMD_MIGRATE}} -path {{MIGRATION_PATH}} -database "{{DATABASE_URL}}" -verbose down

# db-migration-clear-dirty-flag: clear dirty migration flag
db-migration-clear-dirty-flag:
    #!/usr/bin/env bash
    read -p "Version to clear flag: " DIRTY_VERSION
    {{CMD_MIGRATE}} -path {{MIGRATION_PATH}} -database "{{DATABASE_URL}}" force $DIRTY_VERSION

# db-migration-stage: check migration stage
db-migration-stage:
    {{CMD_MIGRATE}} -path {{MIGRATION_PATH}} -database "{{DATABASE_URL}}" version

# db-migration-force-version: force migration to specific version
db-migration-force-version:
    #!/usr/bin/env bash
    read -p "Version to force to: " NEW_VERSION
    {{CMD_MIGRATE}} -path {{MIGRATION_PATH}} -database "{{DATABASE_URL}}" force $NEW_VERSION

# tag-staging: deploy staging by creating a new tag with incremented patch version (use 'push=true' to automatically push)
tag-staging push="false":
    #!/usr/bin/env bash
    echo "Creating new staging tag with incremented patch version..."
    latest_tag=$(git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-staging' | tail -n1)
    if [ -z "$latest_tag" ]; then
        echo "No existing staging tag found. Creating first staging tag v0.0.1-staging"
        new_tag="v0.0.1-staging"
    else
        echo "Found latest staging tag: $latest_tag"
        base_version=$(echo $latest_tag | sed 's/-staging$//')
        major=$(echo $base_version | cut -d. -f1)
        minor=$(echo $base_version | cut -d. -f2)
        patch=$(echo $base_version | cut -d. -f3)
        new_patch=$((patch + 1))
        new_tag="v$major.$minor.$new_patch-staging"
    fi
    echo "Creating new tag: $new_tag"
    git tag $new_tag
    if [ "{{push}}" = "true" ]; then
        echo "Pushing tag $new_tag to origin..."
        git push origin $new_tag
        echo "Tag $new_tag has been created and pushed to origin"
    else
        echo "New staging tag $new_tag created successfully. Run 'git push origin $new_tag' to push the tag to remote."
    fi

# tag-release: deploy production by creating a new tag with incremented patch version (use 'push=true' to automatically push)
tag-release push="false":
    #!/usr/bin/env bash
    echo "Creating new production tag with incremented patch version..."
    latest_tag=$(git tag --sort=version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-release' | tail -n1)
    if [ -z "$latest_tag" ]; then
        echo "No existing production tag found. Creating first production tag v0.0.1-release"
        new_tag="v0.0.1-release"
    else
        echo "Found latest production tag: $latest_tag"
        base_version=$(echo $latest_tag | sed 's/-release$//')
        major=$(echo $base_version | cut -d. -f1)
        minor=$(echo $base_version | cut -d. -f2)
        patch=$(echo $base_version | cut -d. -f3)
        new_patch=$((patch + 1))
        new_tag="v$major.$minor.$new_patch-release"
    fi
    echo "Creating new tag: $new_tag"
    git tag $new_tag
    if [ "{{push}}" = "true" ]; then
        echo "Pushing tag $new_tag to origin..."
        git push origin $new_tag
        echo "Tag $new_tag has been created and pushed to origin"
    else
        echo "New production tag $new_tag created successfully. Run 'git push origin $new_tag' to push the tag to remote."
    fi