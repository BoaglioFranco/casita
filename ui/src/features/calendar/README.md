# Calendar UI

Calendar is a household feature at `/$householdId/calendar`. It does not depend on
a list. The main navigation also exposes a lists overview at `/$householdId/lists`.
Both routes check household membership through the existing authentication guard.

## API contract

- `GET /households/:householdId/calendar?from=YYYY-MM-DD&to=YYYY-MM-DD` returns
  `{ from, to, occurrences }`. The backend interprets date-only boundaries in its
  local timezone and returns timestamp boundaries. Both endpoints are inclusive.
- Each occurrence has its own `startsAt` and `endsAt`, plus the underlying event.
  A recurring event's `event.startsAt` is the series anchor, not necessarily the
  start displayed in the current week. Expansion stays on the server.
- `POST /households/:householdId/calendar/events` accepts `name`, `description`,
  `startsAt`, `durationMinutes`, `notificationOffsetMinutes`, and `recurrenceRule`.
  Success is an empty `201`; the UI invalidates the household's calendar queries.
- Local date/time inputs are converted to an ISO timestamp for writes. Occurrences
  are displayed in the viewer's timezone. Recurrence follows the backend's UTC
  rules and does not adjust for daylight saving changes.
- The form sends `null` for an empty description or a non-repeating event. It
  sends an explicit reminder offset, including `0` for at-start reminders. The
  initial offset is the backend default of 10 minutes. A positive duration is
  required by the UI because timeline generation currently assumes one exists.
- Reminder offsets are saved, but scheduled event push delivery is a later backend
  phase.

The creation form shows a name input and compact summaries for the schedule,
description, and reminder. Mantine Accordion controls expand each row into its
editor. The schedule summary includes the actual date, start/end times, duration,
and repeat choice. Duration uses a native select in 15-minute steps up to 24 hours,
with longer presets for multi-day events. Mantine's modal transition animates
opening and closing; the draft resets after the exit transition. Routine success
uses a 2.5-second Mantine toast. Layout uses Mantine props and its Styles API.

Reminder choices are At start, 5/10/20/30/45 minutes, and 1 through 8 hours.
Picker icons use the same neutral theme color; field focus and the selected
calendar view use the primary brand color.

Date fields use the shared Mantine `CalendarDateInput` in both the toolbar and
event editor. They open a themed popover on desktop and a date-selection dialog
on phones. The entire field, including its icon, activates the picker. Time uses
Mantine `TimePicker` with a clock trigger and quarter-hour choices. Nested pickers
handle Escape without closing the event form. Selected values remain local date
and `HH:mm` strings until the existing form converts them to the API timestamp.

## Fetching and rendering

The initial query spans 28 days: 14 days before today through 14 days after today.
Navigation prefetches another 28-day batch when the visible range approaches
within three days of either loaded boundary. TanStack Query retains all batches
in that window, including earlier events that continue into later days.

Adjacent batches share their inclusive boundary. Merge by event ID **and occurrence
start instant** to remove duplicates while retaining separate recurring instances.
A jump more than 28 days beyond the loaded range opens a new centered four-week
window; it does not fetch all intervening dates. Inactive windows remain cached
for 30 minutes. Queries use the shared application's staleness and retry policy.

Below Mantine's `sm` breakpoint, the time grid offers Day and 3 days and defaults
to Day. Desktop keeps Day and Week, defaulting to Week. Each layout remembers its
own selection when the viewport changes. The three-day grid starts on the selected
date, fits the phone width, and advances three days at a time. Day view shows only
the selected day; Week starts on Monday. Event positions use
the occurrence timestamps, clip across local midnights, and allocate lanes for
overlapping events. Short events have a minimum visible height and full details
are available by selecting them. Time slots support half-hour selection with a
pointer; keyboard users can create an event using the focused 09:00 slot or the
New event button.

Initial load and the Today button center the timeline on the current local time,
including when Today is already selected. Other date navigation preserves the
user's vertical scroll position.

Date navigation slides the date columns in the direction of travel while keeping
the time gutter and vertical scroll position fixed. The browser animates existing
elements so rapid navigation replaces the current animation without delaying date
selection or fetching. Mantine's reduced-motion hook disables this effect when
requested by the device settings.

The backend's initial-window overlap limitation remains deferred in
`backlog/issues/006-include-calendar-events-overlapping-range-start.md`. Events
that start before the earliest fetched boundary cannot be reconstructed by the UI.

## Verification

Run `npm test` from `ui/` with Node 22.6+ to check date boundaries, DST-safe day
arithmetic, inclusive-batch deduplication, recurring instance identity, overlapping
lanes, and multi-day clipping. `npm run build` checks types and the production bundle.

For UI checks, exercise day/week navigation, both fetch boundaries, a distant date
jump, creation failure and retry without losing the draft, empty calendars, and
refresh failure with cached events still visible. Check phone widths down to
320px as well as desktop. Headless Chromium checks covered mocked API responses
down to 320px and the running Go API with PostgreSQL at desktop and phone sizes.
Live checks used the supplied test login and verified event field round trips,
recurrence expansion, rejected recurrence rules with draft retention, overnight
events, fetching in both directions, household isolation, and event deletion.
All events created by these checks were removed afterward.
