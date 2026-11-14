# Database migrations

The API and notification service run golang-migrate before accepting HTTP requests.
SQL is embedded in each binary. Each service uses `DATABASE_URL` and owns the
migration history of its separate database. Never point both services at the same
database. PostgreSQL 18 is required by the API's `uuidv7()` defaults.

The API database is `casita`, the notification database is `notifications`, and
the shared PostgreSQL role is `casita_app`. The databases must already exist and
credentials must permit schema changes.
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

## Required cutover for the Casita rename

The Compose project is now `casita` (previously `ui`). Its PostgreSQL volume is
`casita_postgres-data`, so the first startup creates fresh storage rather than
reusing `ui_postgres-data`. This is an intentional pre-production reset: accounts
and application data must be recreated. The old volume remains available until
explicitly removed. The service JWT issuer also changes to `casita-api`; rebuild
and restart the API and notifications together.

On each existing deployment, including the production host:

1. Stop the old application and PostgreSQL before starting the new project. From
   the repository root, run `docker compose -p ui down`. Do not add `-v`: retaining
   the old volume makes the previous data available for recovery. If development
   servers are running in the old tmux session, stop them with
   `tmux kill-session -t groceries-dev` as well.
2. Update the uncommitted environment files on that host. In `.env`, `api/.env`,
   and `api/.env.prod` where present, set `POSTGRES_DB=casita` and
   `POSTGRES_USER=casita_app`. Change API `DATABASE_URL` values to use user
   `casita_app` and database `casita`; notification URLs use user `casita_app`
   and database `notifications`. Preserve the existing passwords, S2S token,
   VAPID keys, hostnames, and connection options. Compose URLs use host `postgres`;
   services running on the host use `localhost`. Remove or update any external
   `COMPOSE_PROJECT_NAME` setting or `-p ui` launch argument so normal commands
   use the new project name.
3. Initialize the new PostgreSQL instance and notification database:

   ```bash
   docker compose up -d --wait postgres
   make -C notifications db-create
   ```

   PostgreSQL uses the new `.env` values only when initializing empty storage;
   editing these values does not rename roles or databases in an existing volume.
4. For production, rebuild and start the full stack, then recreate each account:

   ```bash
   docker compose up -d --build
   docker compose exec api /app/create-user -email 'you@example.com' -name 'Your name'
   ```

   For local development, use `make dev` and `make -C api create-user` with `EMAIL`
   and `NAME` instead. Startup applies both services' migrations automatically.
5. Sign in again and re-register device notifications as described in step 5 of
   the identity reset below. Fresh databases already satisfy that reset; do not
   drop their schemas a second time.

Do not point the new project at the old volume without a separate database/role
migration and the identity schema reset below. Old data is not imported by this
cutover. Only remove `ui_postgres-data` after verifying the new deployment and
making any desired backup.

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
