# Database migrations

The API and notification service run golang-migrate before accepting HTTP requests.
SQL is embedded in each binary. Each service uses `DATABASE_URL` and owns the
migration history of its separate database. Never point both services at the same
database. PostgreSQL 18 is required by the API's `uuidv7()` defaults.

The databases must already exist and credentials must permit schema changes.
Use `localhost:5432` from the host and `postgres:5432` inside Compose.
`make -C notifications db-create` creates the notification database on the local
Compose PostgreSQL instance if it does not exist yet.

From either service directory:

```bash
make db-up       # pending migrations; also happens automatically on startup
make db-down     # roll back ALL versions, deleting application data
```

Override `DATABASE_URL` for other connection details. API migration 002 still seeds
the demo household (`Casita`) and list (`Compras`); it does not create accounts.

## Required reset for the single-user-identity refactor

This pre-production refactor rewrites API migrations 001 and 003 in place. The old
`members` table is removed; users now have a required name and household. An old
version-3 database is incompatible even though the migration version is unchanged.
**Reset existing databases before running the updated app.** Do not use `db-force`
to claim the old schema matches the new one. The rewritten down migrations also
expect the new schema; use the schema reset below for an old installation.

Reset both service databases together: stored notification subscriptions reference
user UUIDs from the API database, and recreated users receive new UUIDs. This deletes
accounts, sessions, lists, items, activity, events, and stored device subscriptions.
It leaves the PostgreSQL instance and other databases intact.

1. Stop the API, notification service, and UI so nothing writes during the reset.
   For Compose, use `docker compose stop api notifications ui`.
2. Set `API_DATABASE_URL` and `NOTIFICATIONS_DATABASE_URL` to the intended databases,
   using connection addresses reachable from the host. Verify their database names
   differ. Run these commands as a role that owns the schemas/objects:

   ```bash
   psql "$API_DATABASE_URL" -v ON_ERROR_STOP=1 -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'
   psql "$NOTIFICATIONS_DATABASE_URL" -v ON_ERROR_STOP=1 -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'
   make -C api db-up DATABASE_URL="$API_DATABASE_URL"
   make -C notifications db-up DATABASE_URL="$NOTIFICATIONS_DATABASE_URL"
   ```

3. Recreate each account, prompting for its password:

   ```bash
   go -C api run ./cmd/create-user -database-url "$API_DATABASE_URL" -email 'you@example.com' -name 'Your name'
   ```

   The command defaults to the seeded household. Supply `-household-id UUID` for
   another existing household. Each account is a single row in `users`; there is
   no separate membership or membership ID.
4. Rebuild and restart all three application services with the updated code
   (`docker compose up -d --build api notifications ui` for Compose).
5. Fully close and reopen the PWA/browser, sign in again, and open a list or profile.
   Old sessions are invalid. An existing browser push subscription with granted
   permission is automatically registered to the new user ID. Enable notifications
   in profile settings if the device has no subscription. Keep existing VAPID keys
   during the reset so installed devices can reuse their browser subscriptions.

No migration from old records or legacy URL aliases is provided. API resources
now use `/me`, `/lists`, `/calendar`, and `/notifications/subscribe`; browser routes
use `/lists`, `/lists/:listId`, `/calendar`, and `/profile`.

## Deployment and recovery

Fresh databases migrate on first startup. Existing versioned databases apply only
pending files; concurrent runners use PostgreSQL advisory locking. Lock waits are
bounded, so a competing startup may need restarting after the runner finishes.

A migration failure marks the database dirty and stops startup. Inspect and repair
the database before explicitly forcing its verified version. Force records a version
without executing SQL. Never automatically force or roll back on startup failure.
The single-user reset above is an explicit pre-production exception to immutable
migrations. For future changes, add numbered up/down pairs and use compatible changes
while older application versions may still run.

The API image includes `/app/migrate`, supporting `-action up`, `-action down`, and
`-action force -version N`, using the environment's `DATABASE_URL`.

## Validation

Use disposable databases, apply migrations first, and set `TEST_DATABASE_URL` for
each service's Go suite. The API route integration tests and migration lifecycle
tests create and clean isolated schemas. The database transaction test expects the
base API schema to have been migrated already.
