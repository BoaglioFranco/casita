CREATE TABLE calendar_events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID REFERENCES households (id) ON DELETE CASCADE,
    created_by_member_id UUID REFERENCES members (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    duration_minutes INTEGER,
    notification_offset_minutes INTEGER NOT NULL,
    recurrence_rule TEXT, -- RRULE
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE calendar_event_notification_settings (
    event_id UUID NOT NULL REFERENCES calendar_events (id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (event_id, member_id)
);
