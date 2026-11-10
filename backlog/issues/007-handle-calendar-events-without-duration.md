---
id: ISSUE-007
status: open
priority: high
type: bug
created: 2026-11-10
files:
  - api/internal/calendar/calendar_types.go
  - api/internal/calendar/calendar_service.go
  - api/internal/calendar/calendar_store.go
  - api/migrations/003_events.up.sql
---

# Handle calendar events without a duration

## Problem

The API accepts a calendar event without `durationMinutes`, but reading a range
containing that event panics. The current UI submits a duration, but direct API
clients and existing nullable records can trigger the failure.

## Evidence

Discovered during user identity refactor integration tests on PostgreSQL 18.6:

1. Authenticate and POST `/calendar/events` with
   `{"name":"Test","startsAt":"2026-11-11T12:00:00Z"}`. Creation succeeds (201).
2. GET `/calendar?from=2026-11-11&to=2026-11-12`.
3. `CalendarService.GenerateEventTimeline` dereferences `*e.DurationMinutes`
   when computing `EndsAt`, causing a nil-pointer panic. Both recurring and
   one-time occurrence generation use this unchecked dereference.

`CreateEventPayload.DurationMinutes` is optional (`*int`, `omitempty,gt=0`).
Create/update persist the pointer without applying a default. The migration also
allows a NULL `calendar_events.duration_minutes`. Production Gin recovery converts
the panic into an error response; tests without recovery observe the panic.

## Expected outcome

Choose and document a consistent duration policy across create/update, storage,
and timeline generation. Decide whether omitted durations receive a default or
are rejected. The product default remains undecided. Existing NULL records must
be handled without panicking, including recurring events.

## Acceptance criteria

- [ ] Define omitted/null duration behavior for creation and updates.
- [ ] Align request validation, persistence, and schema with that policy.
- [ ] Range reads safely handle existing NULL durations.
- [ ] Cover omitted/null durations for one-time and recurring events.
- [ ] Positive explicit durations retain their existing behavior.

## Notes

The user explicitly deferred this fix while implementing the one-user/one-household
refactor. Its integration fixtures use explicit durations to isolate identity and
authorization behavior. Do not infer a default duration from those fixtures.
