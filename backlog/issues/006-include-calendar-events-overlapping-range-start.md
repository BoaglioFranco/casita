---
id: ISSUE-006
status: open
priority: medium
type: research
created: 2026-11-02
files:
  - api/internal/calendar/calendar_store.go
  - api/internal/calendar/calendar_service.go
---

# Evaluate calendar events overlapping the start of a requested range

## Problem

Calendar range queries select events and occurrences by their start timestamps.
An event starting before `from` but ending after `from` is omitted, even though
it is still in progress during the requested calendar window. This affects both
one-time events and recurring occurrences. Start-based selection is not
necessarily a bug: its suitability depends on the UI's fetching and retention
behavior described below.

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
  `2026-11-02T01:30:00Z` would overlap the window, but the current start-based
  filtering excludes it. The same applies to a recurring occurrence with those
  timestamps. This example is derived from code inspection, not a recorded
  live API test.

## Expected outcome

Evaluate whether the rolling UI window makes start-based selection sufficient.
Backend overlap selection is a possible enhancement, not an approved requirement.
If adopted, calendar results would include one-time events and recurring occurrences that start
before `from` and continue into the requested window. Candidate selection and
recurrence expansion account for duration while retaining the original
occurrence timestamps and the recurring series anchor.

## Acceptance criteria

- [ ] Evaluate the initial-window edge case against the rolling UI fetching plan
      before deciding whether to implement backend overlap selection.

The following criteria apply only if backend overlap selection is chosen:

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

### UI fetching plan

- Initially fetch a rolling four-week window: two weeks before now and two
  weeks after now, rather than a calendar month.
- The UI implementation in `ui/src/features/calendar/CalendarPage.tsx` fetches
  adjacent 28-day batches when the visible range comes within three days of either
  loaded boundary. It retains previously fetched events so ongoing events from
  earlier batches stay available. A date jump more than 28 days beyond the loaded
  range opens a new centered four-week window; inactive windows remain cached
  for 30 minutes. The initial-load overlap edge case also applies to these jumps.
- Avoid showing duplicate occurrences when combining batches. Both range
  endpoints are currently inclusive, so adjacent batches sharing a boundary
  can return the same occurrence. Event ID plus occurrence start identifies an
  occurrence; event ID alone cannot distinguish a recurring series' instances.
- The remaining initial-load edge case is an event that starts before the
  initial `from` (two weeks before now) and continues into the loaded window.
  This is deferred and is not a blocker for the current approach.

The user explicitly deferred implementation. Do not implement this ticket
without a later request. Assume events have a duration for this work, as agreed
for the current timeline implementation; defining missing-duration behavior is
outside this ticket. Recurrence expansion already sets
`rule.DTStart(e.StartsAt.UTC())`; the missing behavior is overlap selection, not
the series anchor.
