# Raspberry Pi deployment plan

## Purpose

Build a repeatable deployment process for the groceries application while
learning each infrastructure component in small steps. The first host will be a
Raspberry Pi to avoid ongoing hosting costs, but the resulting container images
and runtime configuration should remain portable to a cloud VM.

The build must not depend on manually compiling binaries, building the React
application, or copying generated files after each source change.

## Agreed decisions

- The UI and API remain separate services with one Dockerfile each.
- The UI image uses Node to build React and nginx to serve the generated files.
- nginx is the application's single HTTP entry point. It serves the UI and
  proxies `/api/*` requests to the API service.
- The Go service remains API-only; it does not serve the React application.
- nginx removes the `/api` prefix before forwarding requests because the current
  Go routes do not include that prefix.
- SQLite remains a file owned by the API rather than a separate container.
- The SQLite file is stored in persistent mounted storage so container rebuilds
  and replacements do not remove application data.
- Docker Compose defines the local and Raspberry Pi runtime.
- Only nginx is reachable outside the Compose network. The API is internal.
- Cloudflare Tunnel is added after the production-like stack works locally and
  on the Pi's local network.
- `cloudflared` uses Cloudflare's published container image and therefore does
  not need a project Dockerfile.
- Public HTTPS terminates at Cloudflare. nginx and `cloudflared` communicate over
  the private Docker network using HTTP.
- A named tunnel and a `boaglio.dev` subdomain will be used for the alpha; Quick
  Tunnels are only suitable for temporary experiments.
- Automated image publishing and deployment come after the manual Compose flow
  is understood and verified.

## Target architecture

```text
iPhones
   |
   | HTTPS
   v
Cloudflare edge and DNS
   |
   | outbound tunnel connection established by the Pi
   v
cloudflared container
   |
   | private Compose network, HTTP
   v
nginx UI container
   |-- /              -> compiled React application
   |-- /sw.js         -> PWA service worker
   `-- /api/*         -> Go API container, with /api stripped
                              |
                              v
                       persistent SQLite file
```

Cloudflare Tunnel is the only public route to the application. The home router
does not need port forwarding, and the API is never published directly.

## Planned repository layout

```text
compose.yaml
api/
  Dockerfile
  ...
ui/
  Dockerfile
  nginx.conf
  ...
plans/
  README.md
  raspberry-pi-deployment.md
```

Deployment secrets and machine-specific configuration must remain uncommitted.
Examples may be committed with placeholder values.

## Implementation stages

Each stage should be completed and explained before beginning the next one.

### 1. Build the UI image

- Add `ui/Dockerfile`.
- Use a Node build stage with `npm ci` and `npm run build`.
- Use nginx as the runtime stage.
- Copy the generated `dist` directory from the build stage into nginx.
- Add `ui/nginx.conf` to:
  - serve static files;
  - fall back to `index.html` for React client-side routes;
  - proxy `/api/*` to the Compose API hostname;
  - strip the `/api` prefix when proxying.
- Add an appropriate Docker ignore file so local dependencies and build output
  are not sent into the build context.

Completion check: the UI image builds without using the repository's existing
`ui/node_modules` or `ui/dist` directories.

### 2. Build the API image

- Add `api/Dockerfile`.
- Use a Go build stage that downloads locked dependencies and compiles the API.
- Use a minimal runtime stage containing the binary and anything it genuinely
  needs at runtime.
- Run as a non-root user where practical.
- Keep runtime values such as `PORT`, `DB_PATH`, `ENV`, and VAPID keys outside
  the image.
- Add an appropriate Docker ignore file so the local database, environment
  file, temporary binaries, and other development artifacts are excluded.

Completion check: the API image builds from tracked source and starts without a
host-installed Go toolchain.

### 3. Define the local production stack

- Add `compose.yaml` with `ui` and `api` services.
- Give both services a private Compose network.
- Publish only the nginx port for local verification, for example
  `localhost:8080`.
- Do not publish the API port to the host.
- Set `ENV=prod` so production session cookies are secure when the application
  is accessed over HTTPS. Local HTTP testing may require a clearly documented
  development override because browsers do not send secure cookies over HTTP.
- Mount persistent storage at a stable API data path such as `/data` and set
  `DB_PATH` to a file inside it.
- Add restart behavior and health checks once suitable health endpoints exist.

Completion check: `docker compose up --build` builds and starts both services,
and nginx can reach the API using its Compose service name.

### 4. Make database setup repeatable

The current Make targets invoke the host's `sqlite3` executable. A fresh
container instead needs an explicit, container-compatible migration path.

- Choose and document a migration mechanism.
- Prefer an explicit migration command initially so schema changes remain
  visible while learning, for example a command run through the API image.
- Ensure migrations operate on the same mounted database path used at runtime.
- Make initial household and grocery-list creation repeatable.
- Document the one-time commands for creating the two users.
- Generate VAPID keys once, store them as secrets, and preserve them across
  deployments. Replacing the keys would invalidate existing push subscriptions.

Completion check: starting with empty persistent storage, documented commands
can create a usable database and initial users without installing SQLite or Go
on the host.

### 5. Validate the production-like stack locally

Verify the whole browser-to-database path before using the Pi:

- nginx serves the React application.
- Directly refreshing a client-side route returns the application rather than
  an nginx 404.
- `/api` requests reach the Go service with the expected rewritten path.
- Login, listing items, adding an item, and removing an item work.
- Recreating the containers preserves users, sessions, subscriptions, and list
  data.
- Runtime secrets do not appear in images, committed files, or logs.
- The service worker and manifest are served from the expected paths.

Completion check: the application can be exercised through nginx alone and its
data survives `docker compose down` followed by another startup.

### 6. Run the same stack on the Raspberry Pi

- Confirm the Raspberry Pi model, CPU architecture, operating system, and free
  storage before selecting image platforms.
- Install Docker and the Compose plugin.
- Initially build the images on the Pi from the repository.
- Create the Pi's uncommitted environment configuration and persistent data
  directory or named volume.
- Start the stack in detached mode.
- Verify it over the home network before making it public.

Completion check: the same Compose definition runs on the Pi, the application
is usable over the LAN, and a restart of both Docker and the Pi preserves data.

### 7. Add Cloudflare Tunnel

- Ensure the chosen `boaglio.dev` subdomain is managed in Cloudflare DNS.
- Create a named Cloudflare Tunnel.
- Add `cloudflared` to Compose using a pinned official image version.
- Store the tunnel token only in the Pi's uncommitted secret configuration.
- Route the public hostname to `http://ui:80` on the private Compose network.
- Remove any temporary host port publication that is no longer required.
- Verify that the application sees the expected forwarded scheme and host.

Completion check: both iPhones can reach the application at its HTTPS hostname
without Tailscale, a static home IP, or router port forwarding.

### 8. Test the PWA and Web Push flow

- Install the PWA on one iPhone from the production HTTPS hostname.
- Enable notifications through an explicit user interaction.
- Confirm the push subscription is persisted.
- Add or remove an item from the other session and verify notification delivery.
- Repeat on the second iPhone after the single-device flow is reliable.
- Resolve the existing high-priority notification issues before treating push as
  dependable:
  - `backlog/issues/001-handle-push-subscription-storage-errors.md`
  - `backlog/issues/002-continue-push-delivery-after-failure.md`

Completion check: both installed PWAs can authenticate, share the same list,
and receive expected notifications.

### 9. Add operational safety

- Create consistent SQLite backups using SQLite's supported backup mechanism,
  not a blind copy while the database may be writing.
- Copy backups away from the Pi so an SD-card or device failure does not remove
  both the database and its backup.
- Document restore and rollback procedures and test them.
- Add useful container health checks and log inspection commands.
- Pin base image versions and define an update cadence.
- Review login rate limiting, CSRF exposure, secure cookies, error leakage, and
  dependency updates before relying on the public endpoint.

Completion check: a new database can be restored from backup and the previous
application version can be restarted using documented commands.

### 10. Automate builds and later deployments

- Configure a Git remote and choose a container registry.
- Add CI that tests the application and builds versioned `linux/arm64` and
  `linux/amd64` images.
- Reference immutable versions in the deployment rather than relying on
  `latest`.
- Change the Pi workflow from building source to pulling published images.
- Keep deployment automation understandable and reversible; automatic
  production updates are optional and should not be the first iteration.

The eventual update flow should be similar to:

```text
Git commit or tag
  -> CI tests and builds images
  -> registry stores versioned images
  -> Pi pulls the selected version
  -> migrations run explicitly
  -> Compose replaces the containers
```

Completion check: the Pi can deploy or roll back to a selected image version
without compiling or copying application files manually.

## Configuration and secrets

Expected non-secret runtime configuration includes:

- `ENV`
- `PORT`
- `DB_PATH`

Secrets include:

- `VAPID_PRIVATE_KEY`
- the Cloudflare tunnel token
- any future backup credentials

The VAPID public key is not confidential, but it should remain paired with the
private key. Real secrets must never be committed or baked into an image.

## Explicit non-goals for the first deployment

- Kubernetes or another container orchestrator.
- A managed cloud database.
- Multiple API replicas sharing SQLite.
- nginx obtaining public TLS certificates.
- Go serving the frontend.
- Fully automatic deployments before the manual process is understood.
- Replacing the existing Go API, React UI, or SQLite database solely for
  deployment purposes.

## Open decisions

These should be resolved at the stage where they become relevant:

- Raspberry Pi model and whether its operating system is 64-bit.
- Final public hostname under `boaglio.dev`.
- Exact database migration implementation.
- Bind mount versus named volume for Pi persistence and backup ergonomics.
- Backup destination and retention policy.
- Container registry and CI provider.
- Exact pinned Node, nginx, Go, runtime, and `cloudflared` image versions.
