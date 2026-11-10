.PHONY: dev dev-db db-up db-down db-force db-reset db-shell create-user test-integration

SHELL := /bin/bash
DATABASE_URL ?= postgres://groceries:groceries@localhost:5432/groceries?sslmode=disable
HOUSEHOLD_ID ?=

# Run the API with hot reloading
dev:
	air

# Start the local PostgreSQL service and wait until it is healthy
dev-db:
	docker compose -f ../compose.yml --project-directory .. up -d --wait postgres

# Apply database migrations
db-up:
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action up

# Roll back database migrations
db-down:
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action down

# Record an already-applied version; does not execute SQL. Verify the schema first.
db-force:
	@test -n "$(VERSION)" || (echo 'VERSION is required'; exit 1)
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action force -version "$(VERSION)"

# Recreate the local database schema
db-reset: db-down db-up

# Open the local database shell
db-shell:
	psql "$(DATABASE_URL)"

# Create a user in one household, prompting for the password
create-user:
	@set -a; \
	. ./.env; \
	set +a; \
	go run ./cmd/create-user -email "$(EMAIL)" -name "$(NAME)" $(if $(strip $(HOUSEHOLD_ID)),-household-id "$(HOUSEHOLD_ID)")

test-integration:
	TEST_DATABASE_URL="$(DATABASE_URL)" go test ./...
