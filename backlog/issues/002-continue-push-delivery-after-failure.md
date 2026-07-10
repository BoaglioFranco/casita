---
id: ISSUE-002
status: open
priority: high
type: bug
created: 2026-07-10
files:
  - api/internal/notifications/service.go
---

# Continue push delivery after an individual failure

## Problem

Sending a notification to a household stops as soon as delivery to one push
subscription fails. A stale or otherwise invalid subscription can prevent all
subsequent household members from receiving the notification.

## Evidence

`NotifyMembers` returns immediately when `sendNotification` returns an error
inside the subscription loop.

## Expected outcome

The service attempts delivery to every matching subscription even when one or
more individual deliveries fail, while retaining enough failure information for
the caller or logs to expose delivery problems.

## Acceptance criteria

- [ ] Delivery is attempted for every matching subscription
- [ ] One failed subscription does not block later subscriptions
- [ ] Delivery failures remain observable
- [ ] Relevant tests pass
