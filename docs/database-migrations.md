# Database migrations

The API and notification service run golang-migrate before accepting HTTP requests.
SQL is embedded in each binary, independent of the working directory. Each service
uses `DATABASE_URL` and owns the migration history of its separate database.
The database must already exist and credentials must permit schema changes.

PostgreSQL can stay in Docker: the Go driver connects through its published port.
Use `localhost:5432` from the host and `postgres:5432` inside the Compose network.
Make targets can be run from either service directory:

```bash
make db-up       # pending migrations only; also happens automatically on startup
make db-down     # explicitly roll back ALL versions, deleting application data
```

Override `DATABASE_URL` for non-default connection details. Never point both
services at the same database. `notifications/Makefile` provides `db-create` for
creating its database once on the existing local PostgreSQL container.

## Existing databases: baseline once

Previously, Compose initialization scripts or Make/psql applied SQL without tracking
versions. Compose no longer mounts those initialization scripts; startup owns them.
Existing volumes/data are preserved, but an untracked schema will fail on startup
because its tables already exist. There is no automatic baseline or reset.

Before starting the new binaries, inspect the schema and confirm exactly which
migration files have already been fully applied. Back up any data you need. Then
record that version using the relevant service's Makefile, for example:

```bash
# API: ONLY if migrations 001, 002, and 003 are already fully applied:
make -C api db-force VERSION=3

# Notifications: ONLY if its 001 migration is already fully applied:
make -C notifications db-force VERSION=1
```

These examples use each Makefile's local database URL; supply the correct
`DATABASE_URL` for the target environment. Force records a version without running
SQL. Never use it on an empty database or to conceal an unresolved migration failure.
An earlier schema must be baselined at its actual applied version instead.

## Deployment and recovery

Fresh databases migrate on first startup. Existing versioned databases apply only
pending files; concurrent runners use PostgreSQL advisory locking. Lock waits are
bounded, so a competing startup may fail and need restarting after the runner finishes.

A migration failure marks the database dirty and stops startup. Inspect and repair
the database before explicitly forcing the verified version. Do not automatically
force or run down migrations on failure. Treat applied migration files as immutable;
add numbered up/down pairs for future changes.

Startup migrations also run in production. Use backward-compatible changes while
older application versions may still be running. The API image includes
`/app/migrate` for explicit administration, with `-action up`, `-action down`, or
`-action force -version N`, using `DATABASE_URL` from the environment.

The API's existing demo seed remains migration 002 and therefore still runs on fresh
databases; this change does not alter its contents or move seed data to a separate flow.
