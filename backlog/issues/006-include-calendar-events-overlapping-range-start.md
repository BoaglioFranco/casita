---
id: ISSUE-006
status: open
priority: medium
type: bug
created: 2026-11-02
files:
  - api/internal/calendar/calendar_store.go
  - api/internal/calendar/calendar_service.go
---

# Include calendar events overlapping the start of a requested range

## Problem

Calendar range queries select events and occurrences by their start timestamps.
An event starting before `from` but ending after `from` is omitted, even though
it is still in progress during the requested calendar window. This affects both
one-time events and recurring occurrences.

## Evidence

- `CalendarStore.GetEvents` in `api/internal/calendar/calendar_store.go`
  requires `starts_at >= $2` for one-time events, where `$2` is `from`.
  This excludes overlapping events before timeline generation can see them.
- `CalendarService.GenerateEventTimeline` in
  `api/internal/calendar/calendar_service.go` expands recurring starts using
  `rule.Between(from, to, true)`. Occurrences starting before `from` are not
  expanded, regardless of their duration.
- Reproduction example: request `from=2026-11-02T00:00:00Z` through
  `to=2026-11-03T00:00:00Z`. An event starting at
  `2026-11-01T23:30:00Z` with a 120-minute duration ends at
  `2026-11-02T01:30:00Z` and should be visible, but the current start-based
  filtering excludes it. The same applies to a recurring occurrence with those
  timestamps. This example is derived from code inspection, not a recorded
  live API test.

## Expected outcome

Calendar results include one-time events and recurring occurrences that start
before `from` and continue into the requested window. Candidate selection and
recurrence expansion account for duration while retaining the original
occurrence timestamps and the recurring series anchor.

## Acceptance criteria

- [ ] One-time events starting before `from` and ending after `from` are returned.
- [ ] Recurring occurrences starting before `from` and ending after `from` are returned.
- [ ] Events spanning the entire requested window are returned.
- [ ] Events ending strictly before `from` or starting after `to` are excluded.
- [ ] Existing inclusive start boundaries are preserved: occurrences starting
      exactly at `from` or `to` remain eligible.
- [ ] Decide and document whether an event ending exactly at `from` is included;
      this end-boundary behavior has not yet been specified.
- [ ] Simple regression tests cover overlap for one-time and recurring events,
      non-overlapping events, and the agreed boundary behavior.
- [ ] Household isolation and chronological result ordering are preserved.

## Notes

The user explicitly deferred implementation. Do not implement this ticket
without a later request. Assume events have a duration for this work, as agreed
for the current timeline implementation; defining missing-duration behavior is
outside this ticket. Recurrence expansion already sets
`rule.DTStart(e.StartsAt.UTC())`; the missing behavior is overlap selection, not
the series anchor.
