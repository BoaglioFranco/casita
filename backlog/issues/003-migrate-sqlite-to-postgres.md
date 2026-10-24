---
id: ISSUE-003
status: in-progress
priority: medium
type: maintenance
created: 2026-10-21
files:
  - compose.yml
  - api/internal/database/database.go
  - api/migrations/001_init.up.sql
  - api/migrations/002_seed_demo.up.sql
  - api/cmd/create-user/main.go
  - api/internal/auth/store.go
  - api/internal/households/store.go
  - api/internal/lists/store.go
  - api/internal/items/store.go
  - api/internal/activity/store.go
  - api/internal/notifications/store.go
---

# Migrate persistent application data from SQLite to PostgreSQL

## Problem

The first Raspberry Pi deployment will use SQLite for simplicity, but the
application is intended to move to PostgreSQL soon afterward. The current
database initialization, schema, seed migration, and SQL statements are
SQLite-specific, so changing only the database driver or Compose service would
not be sufficient.

## Evidence

- `api/internal/database/database.go` opens the `modernc.org/sqlite` driver,
  limits the connection pool to one connection, and executes `PRAGMA
  foreign_keys = ON`.
- `api/migrations/001_init.up.sql` uses SQLite-specific `INTEGER PRIMARY KEY
  AUTOINCREMENT` and `strftime(...)` expressions.
- `api/migrations/002_seed_demo.up.sql` uses `INSERT OR IGNORE`.
- Store queries and `api/cmd/create-user/main.go` use SQLite `?` parameter
  placeholders; the PostgreSQL driver and query strategy still need to be
  selected.
- The Compose stack does not yet define a PostgreSQL service or persistent
  PostgreSQL data volume.

## Expected outcome

The API uses PostgreSQL in local and Raspberry Pi deployments. This is a clean
pre-production cutover; existing SQLite data does not need to be transferred.

## Acceptance criteria

- [ ] PostgreSQL is defined as an internal Compose service with persistent storage
- [ ] Database credentials are supplied through uncommitted runtime configuration
- [ ] The API connects to PostgreSQL and waits or retries appropriately during startup
- [ ] Migrations create a PostgreSQL-compatible schema from an empty database
- [ ] Queries and transaction behavior work with the selected PostgreSQL driver
- [ ] Relevant API tests pass against PostgreSQL

## Notes

The user confirmed on 2026-10-23 that the application is pre-production and no
SQLite data preservation or transfer path is required.
