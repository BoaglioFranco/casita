.PHONY: dev db-up db-down db-reset db-shell create-user

SHELL := /bin/bash
DB_PATH ?= groceries.db
HOUSEHOLD_ID ?= 1

# Run the API with hot reloading
dev:
	air

# Apply database migrations
db-up:
	sqlite3 $(DB_PATH) ".read migrations/001_init.up.sql"
	sqlite3 $(DB_PATH) ".read migrations/002_seed_demo.up.sql"

# Roll back database migrations
db-down:
	sqlite3 $(DB_PATH) ".read migrations/002_seed_demo.down.sql"
	sqlite3 $(DB_PATH) ".read migrations/001_init.down.sql"

# Recreate the local database schema
db-reset: db-down db-up

# Open the local database shell
db-shell:
	sqlite3 $(DB_PATH)

# Create a user and membership, prompting for the password
create-user:
	go run ./cmd/create-user -db "$(DB_PATH)" -email "$(EMAIL)" -name "$(NAME)" -household-id "$(HOUSEHOLD_ID)"
