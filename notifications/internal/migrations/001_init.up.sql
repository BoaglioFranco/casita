CREATE TABLE push_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    icon TEXT
);

CREATE TABLE notification_schedules (
    schedule_id UUID PRIMARY KEY,
    schedule_version INTEGER NOT NULL CHECK (schedule_version > 0),
    deleted BOOLEAN NOT NULL DEFAULT false,
    starts_at TIMESTAMPTZ NOT NULL,
    recurrence_rule TEXT,
    reminder_offset_minutes INTEGER NOT NULL CHECK (reminder_offset_minutes >= 0),
    notification_id UUID NOT NULL REFERENCES notifications(id)
);

CREATE TABLE notification_schedule_recipients (
    user_id UUID NOT NULL,
    schedule_id UUID NOT NULL REFERENCES notification_schedules(schedule_id),
    PRIMARY KEY (schedule_id, user_id)
);

CREATE TABLE notification_jobs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    schedule_id UUID NOT NULL REFERENCES notification_schedules(schedule_id),
    schedule_version INTEGER NOT NULL CHECK (schedule_version > 0),
    runs_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'expanded', 'cancelled')),
    UNIQUE (schedule_id, schedule_version, runs_at)
);

CREATE TABLE notification_deliveries (
    job_id UUID NOT NULL REFERENCES notification_jobs(id),
    push_subscription_id UUID NOT NULL REFERENCES push_subscriptions(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'cancelled')),
    attempt_count INTEGER NOT NULL DEFAULT 0 CHECK (attempt_count >= 0),
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    last_error TEXT,
    claim_token UUID,
    claimed_until TIMESTAMPTZ,
    CHECK ((claim_token IS NULL) = (claimed_until IS NULL)),
    PRIMARY KEY (job_id, push_subscription_id)
);

CREATE INDEX notification_jobs_pending_runs_at_idx
ON notification_jobs (runs_at) WHERE status = 'pending';

CREATE INDEX notification_deliveries_pending_next_attempt_at_idx
ON notification_deliveries (next_attempt_at) WHERE status = 'pending';

CREATE INDEX notification_deliveries_processing_claimed_until_idx
ON notification_deliveries (claimed_until) WHERE status = 'processing';

CREATE INDEX notification_deliveries_push_subscription_id_idx
ON notification_deliveries (push_subscription_id);

CREATE INDEX push_subscriptions_user_id_idx
ON push_subscriptions (user_id);
