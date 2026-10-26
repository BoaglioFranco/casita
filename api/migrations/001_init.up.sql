CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    password_hash TEXT NOT NULL
);

CREATE TABLE sessions (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    UNIQUE (user_id, household_id)
);

CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    icon TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    list_id UUID NOT NULL REFERENCES lists (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    created_by_member_id UUID REFERENCES members (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE push_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE list_notification_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    list_id UUID NOT NULL REFERENCES lists (id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
    push_subscription_id UUID NOT NULL REFERENCES push_subscriptions (
        id
    ) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    UNIQUE (list_id, member_id, push_subscription_id)
);

CREATE TABLE activity_events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    list_id UUID REFERENCES lists (id) ON DELETE SET NULL,
    item_id UUID REFERENCES items (id) ON DELETE SET NULL,
    actor_member_id UUID REFERENCES members (id) ON DELETE SET NULL,
    event_type TEXT NOT NULL,
    item_name_snapshot TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX members_household_id_idx ON members (household_id);
CREATE INDEX lists_household_id_idx ON lists (household_id);
CREATE INDEX items_list_id_created_at_idx ON items (list_id, created_at DESC);
CREATE INDEX push_subscriptions_household_id_idx ON push_subscriptions (
    household_id
);
CREATE INDEX activity_events_household_id_created_at_idx
ON activity_events (household_id, created_at DESC);
