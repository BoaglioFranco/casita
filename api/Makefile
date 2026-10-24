.PHONY: dev db-up db-down db-reset db-shell create-user test-integration

SHELL := /bin/bash
DATABASE_URL ?= postgres://groceries:groceries@localhost:5432/groceries?sslmode=disable
HOUSEHOLD_ID ?= 1

# Run the API with hot reloading
dev:
	air

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
	go run ./cmd/create-user -database-url "$(DATABASE_URL)" -email "$(EMAIL)" -name "$(NAME)" -household-id "$(HOUSEHOLD_ID)"

test-integration:
	TEST_DATABASE_URL="$(DATABASE_URL)" go test ./...
