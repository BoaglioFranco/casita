# Events and notification model

Created: 2026-10-29. Updated: 2026-10-30.
Design notes only; implementation is pending. Keep these decisions in this file.

## Implementation order

1. Develop the events feature in the main backend and UI, functional without
   notification delivery. Include one-time and recurring events, duration,
   household access, and calendar date-range queries.
2. Introduce the notification service and migrate existing list notification
   delivery into it. Use immediate list notifications to exercise the outbox,
   Debezium, Kafka, and consumer pipeline before adding recurring delivery.
3. Add scheduled and recurring event notifications through the service.
4. Consider one-way Google Calendar export later; it is outside the MVP.

The decoupling design is recorded first below; migrating delivery waits until
the events feature is working. Exact timing of the schema cleanup remains open.

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
to be decided. Whole-series versus individual-occurrence edits also remains open.

Calendar date-range queries expand recurrence rules for the requested range;
displaying an event years ahead does not require persisting years of notification
jobs. Calendar rendering and notification scheduling must use consistent UTC
recurrence semantics, potentially through a shared library.

## 3. Notification service architecture

This design is implementable. Its additional infrastructure serves the project's
backend and infrastructure learning goals, rather than an expected load of only
5–10 users. Kafka is an intentional learning choice, not a throughput requirement.

```text
Main backend transaction
  event/list change + explicit outbox message
                    |
                    v
Debezium through Kafka Connect
  captures committed outbox inserts from PostgreSQL's change log
                    |
                    v
                  Kafka
                    |
                    v
Notification service + its own database
  notification schedules -> upcoming jobs -> delivery workers -> push providers
```

Debezium runs alongside PostgreSQL; it is not installed as application logic in
the database. It replaces a custom polling outbox relay. Kafka transports and
retains messages; the notification service handles notification timing.

### Ownership

| Main backend | Notification service |
| --- | --- |
| Events, lists, household membership and permissions | Device push subscriptions |
| Event recurrence and calendar rendering | Notification schedule copies and recurrence expansion for delivery |
| Entity-specific recipient selection | Generation of concrete upcoming jobs |
| Notification intent and schedule changes | Job claiming, leases, sending, retries, delivery history |

The service accepts notification intent for immediate delivery, delivery once at
a timestamp, or recurring delivery. List changes normally request immediate
delivery. Event requests include recurrence, anchor, reminder offset, recipients,
and content. Exact message schemas and identifiers are not yet defined.

The service should not query the main backend's tables. PostgreSQL is the current
database recommendation for both services, with separate databases and credentials;
they can initially share one PostgreSQL instance. MongoDB was considered but no
workload advantage has been established.

### Synchronization and delivery

- Save each domain change and its outgoing outbox message in one transaction.
- Consume messages durably and tolerate redelivery using message IDs.
- Use stable schedule IDs and increasing versions to ignore outdated updates.
  Retain cancellation state/version so an old message cannot restore a schedule.
- Apply incoming changes and record their processing atomically in the
  notification service database. Update/cancel obsolete pending jobs accordingly.
- Accept eventual consistency: updates and cancellations can lag, and an old
  reminder might be sent before the change arrives. Already-sent pushes cannot
  reliably be retracted.
- Persist only a limited horizon of upcoming jobs and extend it periodically.
  Seven days was an example, not a settled value.
- Workers claim due jobs with short database transactions and time-limited leases;
  send outside the transaction. Claim tokens prevent stale workers from recording
  results after reassignment. PostgreSQL `FOR UPDATE SKIP LOCKED` is a candidate.
- Expired leases allow crash recovery. Exactly-once push delivery is not promised:
  a crash after sending but before recording success can result in a duplicate.

### Decisions to resume with

- Message contract, acceptance semantics, Kafka keys/order, consumer recovery,
  and safe handling of retries and stale versions.
- Whether a job represents an occurrence, member, or device; tracking partial
  success without unnecessarily retrying successful device deliveries.
- Generation horizon, polling interval, lease duration/renewal, retry policy,
  and how late a reminder may be delivered.
- Content and recipient changes, household departures, device changes, and
  ownership of any future category-wide preferences.
- Migration of existing device subscriptions and list delivery without losing
  subscriptions or unintentionally running both delivery paths.

## 4. Deployment and resource validation

Start on the existing Raspberry Pi 5 with 4 GB RAM. At idle, the existing stack
showed roughly 385 MiB used, 3.6 GiB available, and no swap usage. Docker reported
`0B / 0B` memory for containers, so those per-container readings are unusable.

Run services as independently deployable processes/containers on the same host.
Start with one Kafka broker; this is a learning deployment with a shared host
failure point. Measure before assuming either that it fits or needs new hardware.

Kafka and Connect maximum JVM heaps of 512 MiB each were suggested as experimental
starting settings, not validated requirements or total process memory limits.
Check host/process memory, CPU, disk growth, startup, snapshots, backlog recovery,
and restarts. Bound Kafka retention and account for PostgreSQL WAL retained while
CDC is behind. Confirm ARM64 image support for selected versions.

Later, migrate the containers to AWS. A single Lightsail or EC2 host is an option;
managed Kafka and multiple machines are not required for this project's traffic.
Choose capacity from measured usage and recheck pricing at that time.
