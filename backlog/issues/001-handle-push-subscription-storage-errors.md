---
id: ISSUE-001
status: open
priority: high
type: bug
created: 2026-07-10
files:
  - api/internal/notifications/handlers.go
  - api/internal/notifications/store.go
---

# Handle push subscription storage errors

## Problem

The subscription handler returns HTTP 200 even when storing the browser push
subscription fails. The UI can therefore report that notifications are enabled
when the subscription was not persisted.

## Evidence

`subscribeHandler` calls `s.store.StoreSubscription(...)` without checking its
returned error and then unconditionally responds with HTTP 200.

## Expected outcome

The handler checks the store error and returns an appropriate failure response
when the subscription cannot be persisted.

## Acceptance criteria

- [ ] Errors returned by `StoreSubscription` are handled
- [ ] A failed write does not return HTTP 200
- [ ] A successful write still returns a success response
- [ ] Relevant tests pass
