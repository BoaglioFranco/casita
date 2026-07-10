# Project backlog

This directory stores deferred work as repository-local, agent-friendly tickets.
Each ticket is a Markdown file in `issues/`, named with a sequential identifier
and a short kebab-case description, for example:

```text
issues/001-handle-expired-push-subscriptions.md
```

Use the next available number when creating a ticket. Keep one independently
actionable problem per file.

## Ticket template

```markdown
---
id: ISSUE-001
status: open
priority: medium
type: bug
created: YYYY-MM-DD
files:
  - path/to/relevant-file
---

# Short, actionable title

## Problem

Describe what was observed and why it matters.

## Evidence

Record errors, reproduction details, relevant behavior, or code locations.

## Expected outcome

Describe the desired behavior.

## Acceptance criteria

- [ ] The root cause is addressed
- [ ] Relevant tests pass
- [ ] Related behavior does not regress

## Notes

Add useful constraints, uncertainties, or implementation context. Omit this
section when it adds no value.
```

## Field values

- `status`: `open`, `in-progress`, `blocked`, or `closed`
- `priority`: `low`, `medium`, `high`, or `critical`
- `type`: `bug`, `feature`, `maintenance`, `performance`, `security`, or `research`

