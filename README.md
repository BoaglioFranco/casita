Casita is a household app for shared grocery lists, calendar events, and push notifications. The push notifications are very important.
We have 2 iphones, i don't intend to share this app with anyone and i don't have a developer account. So i came to the conclusion that the best way to implement this is with a PWA. 
The backend is a small Go server backed by PostgreSQL, and the UI is built with React. The PWA push notification flow is the core feature.

## Local development

From the repository root:

```bash
make dev        # start PostgreSQL and all app servers, then attach to tmux
make dev-start  # start without attaching
make dev-stop   # stop app servers; leave PostgreSQL and its data running
```

The `casita-dev` tmux session has `api`, `notifications`, and `ui` windows.
Open <http://localhost:5173>. Use `Ctrl-b`, then `n` to switch windows and
`Ctrl-b`, then `d` to detach while leaving the servers running. Run `make dev`
again to reattach; inside tmux it switches to that session. In a noninteractive
terminal, startup prints the attach command instead.

The API uses Air for Go reloads; the UI uses Vite. Restart the notification
service after changing its Go code, or restart the stack with `make dev-stop`
followed by `make dev`. Failed service windows stay visible for diagnosis.
Restart the stack after changing environment files.

The Compose project is `casita`, with API database `casita`, notification database
`notifications`, and PostgreSQL role `casita_app`. Existing installations must
follow the [Casita cutover procedure](docs/database-migrations.md#required-cutover-for-the-casita-rename)
before restarting; the renamed Compose project initializes fresh storage.

First-time prerequisites:

- Install Docker with Compose, Go (matching `api/go.mod`), Node 22+, Bash,
  tmux, and Air (`go install github.com/air-verse/air@latest`). The launcher finds
  Air on PATH, in GOBIN, or in GOPATH/bin.
- Copy `.env.example` to `.env` and `api/.env.example` to `api/.env`. Configure
  matching PostgreSQL credentials, a shared S2S token, and a VAPID keypair.
- Run `npm --prefix ui ci` once to install UI dependencies.
- Start PostgreSQL with `make -C api dev-db`, then create the notification
  database once with `make -C notifications db-create`.

`dev.sh` loads `api/.env` for the API. Notifications inherit its database
connection and shared keys, with the database name changed to `notifications`;
optional values in `notifications/.env` override those defaults. Shared keys must
match. Compose's `postgres` hostname is translated to loopback for host processes.
Local ports are fixed at API 2345, notifications 2346, and UI 5173. Environment
files are sourced by Bash: use `KEY=value` without spaces around `=`, and quote
values containing spaces or shell metacharacters. Changes take effect after a
stack restart.

Pending migrations run automatically on service startup. Normal development
startup never resets the database or recreates users. See the reset procedure
below when upgrading an installation that still has the old identity schema.

Launcher tests: `bash tests/dev_test.sh`.

Both services apply pending database migrations on startup. See
[database migrations](docs/database-migrations.md) for setup, schema resets,
and recovery instructions.

Users have one canonical UUID and belong to exactly one household. Account creation
requires a display name and an existing household; `api/cmd/create-user` defaults to
the seeded Casita household. The authenticated API derives household scope from the
session and exposes `/me`, `/lists`, and `/calendar`.

Existing installations must follow the [single-user identity reset procedure](docs/database-migrations.md#required-reset-for-the-single-user-identity-refactor)
before running this version. Both databases are reset together because stored push
subscriptions reference the users being recreated.
