.PHONY: dev dev-db db-up db-down db-reset db-shell create-user test-integration

SHELL := /bin/bash
DATABASE_URL ?= postgres://groceries:groceries@localhost:5432/groceries?sslmode=disable
HOUSEHOLD_ID ?= 1

# Run the API with hot reloading
dev:
	air

# Start the local PostgreSQL service and wait until it is healthy
dev-db:
	docker compose -f ../compose.yml --project-directory .. up -d --wait postgres

# Apply database migrations
db-up:
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f migrations/001_init.up.sql
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f migrations/002_seed_demo.up.sql

# Roll back database migrations
db-down:
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f migrations/002_seed_demo.down.sql
	psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f migrations/001_init.down.sql

# Recreate the local database schema
db-reset: db-down db-up

# Open the local database shell
db-shell:
	psql "$(DATABASE_URL)"

# Create a user and membership, prompting for the password
create-user:
	@set -a; \
	. ./.env; \
	set +a; \
	go run ./cmd/create-user -email "$(EMAIL)" -name "$(NAME)" -household-id "$(HOUSEHOLD_ID)"

test-integration:
	TEST_DATABASE_URL="$(DATABASE_URL)" go test ./...
