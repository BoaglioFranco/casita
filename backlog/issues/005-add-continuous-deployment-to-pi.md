---
id: ISSUE-005
status: open
priority: medium
type: feature
created: 2026-10-23
files:
  - compose.yml
  - .github/workflows/
---

# Add controlled continuous deployment to the Raspberry Pi

## Problem

Deployments currently require logging into the Raspberry Pi, pulling source,
building images locally, and running Docker Compose. Introduce a deliberate
GitHub Actions release and deployment flow after the PostgreSQL migration is
stable. The initial idea was to deploy whenever a `prod` branch is updated, but
release tags or a manual promotion workflow may provide clearer versioning and
rollback without maintaining a second long-lived branch.

## Evidence

- `compose.yml` currently builds the API and UI from local source on the Pi.
- The production stack runs on a home-network Raspberry Pi behind Cloudflare
  Tunnel, so a GitHub-hosted runner cannot directly assume inbound access to
  the host.
- GitHub Actions supports branch pushes, tag pushes, and manual
  `workflow_dispatch` events as deployment triggers.
- GitHub deployment environments can track deployments, restrict sources, and
  gate access to secrets, but availability and protection features depend on
  repository visibility and the GitHub plan:
  https://docs.github.com/en/actions/concepts/workflows-and-actions/deployment-environments
- Database migrations and rollback behavior must be defined before unattended
  deployment is safe.

## Expected outcome

A tested commit is built into immutable ARM64 and AMD64 container images, and a
clearly identified release can be deployed to the Pi without compiling there.
The deployed version is observable and can be rolled back deliberately. The
trigger should be chosen after comparing a `prod` branch, version tags, and
manual promotion; release tags plus a production environment are the current
preferred design.

## Acceptance criteria

- [ ] CI tests the API and UI before publishing images
- [ ] API and UI images are published for `linux/arm64` and `linux/amd64`
- [ ] Images use immutable commit or release tags rather than relying on `latest`
- [ ] `prod` branch, release-tag, and manual-promotion triggers are compared and documented
- [ ] The mechanism by which a workflow reaches or signals the Pi is threat-modeled
- [ ] Production secrets remain on the Pi and are not baked into images
- [ ] Only one production deployment can run at a time
- [ ] Database backup and migration steps fail safely before application replacement
- [ ] A health check verifies the new deployment
- [ ] A documented rollback restores the previous application image and accounts for schema compatibility
- [ ] The deployed commit or release is visible without inspecting container internals

## Notes

Prefer an incremental rollout: first publish images in CI and deploy them
manually on the Pi; automate the final deployment only after that path and its
rollback have been exercised. A deploy-only self-hosted runner on the Pi is one
possible outbound-only control channel, but access to the Docker socket is
effectively host-level authority and must be treated accordingly.
