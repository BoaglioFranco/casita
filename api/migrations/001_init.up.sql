CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    password_hash TEXT NOT NULL
);

CREATE TABLE sessions (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
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
    created_by_user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE TABLE activity_events (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    household_id UUID NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    list_id UUID REFERENCES lists (id) ON DELETE SET NULL,
    item_id UUID REFERENCES items (id) ON DELETE SET NULL,
    actor_user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    event_type TEXT NOT NULL,
    item_name_snapshot TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX users_household_id_idx ON users (household_id);
CREATE INDEX lists_household_id_idx ON lists (household_id);
CREATE INDEX items_list_id_created_at_idx ON items (list_id, created_at DESC);
CREATE INDEX activity_events_household_id_created_at_idx
ON activity_events (household_id, created_at DESC);
