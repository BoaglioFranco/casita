.PHONY: db-create db-up db-down db-force

COMPOSE := docker compose -f ../compose.yml --project-directory ..
DATABASE_URL ?= postgres://groceries:groceries@localhost:5432/notifications?sslmode=disable

# Create the service database once, using the container's PostgreSQL admin user.
db-create:
	$(COMPOSE) exec -T postgres sh -c 'exec createdb -U "$$POSTGRES_USER" notifications'

# Apply pending, versioned migrations using Go's PostgreSQL driver.
db-up:
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action up

# Removes the subscription table and its data from the notifications database.
db-down:
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action down

# Record an already-applied version; does not execute SQL. Verify the schema first.
db-force:
	@test -n "$(VERSION)" || (echo 'VERSION is required'; exit 1)
	DATABASE_URL="$(DATABASE_URL)" go run ./cmd/migrate -action force -version "$(VERSION)"
