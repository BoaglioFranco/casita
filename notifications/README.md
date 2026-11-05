# Notifications

## Service authentication

`POST /notifications/subscribe` requires `Authorization: Bearer <JWT>`. The main
backend signs the JWT with HS256 using the shared `S2S_TOKEN`; it derives the user
ID from the authenticated session. The token must contain:

- `sub`: a nonzero user UUID.
- `iss`: `groceries-api`.
- `aud`: `notifications` (or an audience array containing it).
- `exp`: expiration as a Unix timestamp; issue short-lived tokens, e.g. five minutes.

The service validates the signature, algorithm, issuer, audience, expiration, and
`nbf` when present, then sets the verified UUID in Gin's `user_id` context.
Identity headers or body fields do not establish identity. Invalid requests return
401 without exposing token details. `/health` and `/vapid-public-key` remain public.

Both services must load the same `S2S_TOKEN` into their environment. The notification
service refuses startup if it is missing or shorter than 32 bytes. `.env` files are
not automatically loaded by `go run`; export them or use your development runner.
The main backend's token issuance/forwarding is not implemented yet.

## Database migrations

Use the existing PostgreSQL container, with a separate `notifications` database.
Run from this directory, after starting PostgreSQL from the repository root with
`docker compose up -d --wait postgres`:

```bash
make db-create  # once, to create the database
make db-up     # optional: startup also applies pending migrations
```

Both services apply embedded SQL with golang-migrate before serving HTTP. Each
database has its own `schema_migrations` version table. Repeated startup is safe;
PostgreSQL advisory locks coordinate concurrent migration runners. A failed or
dirty migration prevents startup. Startup never rolls back or forces a version.

Set `DATABASE_URL` when starting the service. It must point to the existing
notifications database, not the main backend database. From the host, use the
published PostgreSQL port (`localhost:5432`); from another Compose container,
use `postgres:5432`. The Make targets have a local development URL that can be
overridden with `DATABASE_URL`. No local `psql` or migrate CLI is required.

`db-create` uses the PostgreSQL container's admin user. The service's startup
credentials need DDL permissions in its own database because it runs migrations.

`make db-down` explicitly rolls back all migrations and deletes subscription data.
For shared deployment/baseline instructions, see [database-migrations.md](../docs/database-migrations.md).
