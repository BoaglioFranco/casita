# Events and notification model

Date: 2026-10-29. Design notes only; implementation is pending.

## 1. Decouple notification selection from devices

Remove `push_subscription_id` from `list_notification_subscriptions`.
Selection belongs to a household member, independently of their devices:

```text
list_notification_subscriptions
  list_id    FK -> lists
  member_id  FK -> members
  PRIMARY KEY (list_id, member_id)
```

A member can have multiple `push_subscriptions`. Delivery resolves all devices
for subscribed members. Adding or removing a device must not change list preferences.

The current schema in `api/migrations/001_init.up.sql` includes
`push_subscription_id`; current delivery in
`api/internal/notifications/notification_store.go` queries household subscriptions
without consulting list selection. Both need to follow the intended model.

Acceptance criteria:

- [ ] List selection no longer references an individual push subscription.
- [ ] Each list/member pair can occur only once.
- [ ] List delivery honors selection and resolves all selected members' devices.
- [ ] Adding or removing devices preserves list selection.
- [ ] Update the existing initial migration in place.

The existing initial migration can be edited in place; a separate migration is
not required for this change.

## 2. Proposed event model

```text
events
  id
  household_id              FK -> households
  created_by_member_id      FK -> members
  title
  description               nullable
  starts_at                 TIMESTAMPTZ
  duration_minutes          positive integer
  recurrence_rule           nullable
  reminder_offset_minutes   default 10; nonnegative
  created_at
  updated_at

event_notification_settings
  event_id                  FK -> events
  member_id                 FK -> members
  enabled                   default true
  PRIMARY KEY (event_id, member_id)
```

- Events are shared and visible across the household.
- Store and evaluate schedules in UTC. No timezone or original offset field;
  no daylight-saving adjustments. Display occurrences in the viewer's timezone.
- No recurrence rule means a one-time event. Daily and fortnightly schedules
  repeat from `starts_at` every 24 hours and 14 days respectively.
- End time is derived from start plus duration, for each occurrence.
- Everyone is selected for notifications by default. Authorized editors can
  toggle members on creation or edit; members can toggle themselves.
- Event selection applies across all of a member's devices, as with lists.
- One-way export to Google Calendar is desired; its model remains undecided.

Recurrence syntax, duration default, and behavior for newly joined members remain
to be decided. Delivery workers and synchronization architecture are outside these notes.
