#!/usr/bin/env bash
set -euo pipefail

test_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
source "$test_root/dev.sh"
fixture_root=$(mktemp -d "${TMPDIR:-/tmp}/casita-dev-test.XXXXXX")
trap 'rm -rf -- "$fixture_root"' EXIT

reset_fixture() {
    project_root="$fixture_root/project with spaces"
    mkdir -p "$project_root/api" "$project_root/notifications"
    rm -f -- "$project_root/notifications/.env"
    cat > "$project_root/api/.env" <<'EOF'
DATABASE_URL='postgres://user:password@postgres:5432/casita?sslmode=disable'
S2S_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
EOF
    cat > "$project_root/notifications/.env" <<'EOF'
VAPID_PUBLIC_KEY=public
VAPID_PRIVATE_KEY='private#quoted-value'
EOF
}

test_api_environment() {
    reset_fixture
    rm "$project_root/notifications/.env"
    export VAPID_PUBLIC_KEY=stale-public VAPID_PRIVATE_KEY=stale-private
    load_environment api
    [[ "$DATABASE_URL" == 'postgres://user:password@127.0.0.1:5432/casita?sslmode=disable' ]]
    [[ "$NOTIFICATIONS_SERVICE_URL" == http://localhost:2346 && "$PORT" == 2345 && "$ENV" == dev ]]
    [[ ! -v VAPID_PUBLIC_KEY && ! -v VAPID_PRIVATE_KEY ]]
    [[ $(export -p) == *'DATABASE_URL='* ]]
}

test_notification_defaults() {
    reset_fixture
    load_environment notifications
    [[ "$DATABASE_URL" == 'postgres://user:password@127.0.0.1:5432/notifications?sslmode=disable' ]]
    [[ "$PORT" == 2346 && "$VAPID_PUBLIC_KEY" == public ]]
    [[ "$VAPID_PRIVATE_KEY" == 'private#quoted-value' ]]
}

test_notification_database_override() {
    reset_fixture
    printf '%s\n' 'DATABASE_URL=postgres://user:password@localhost:5432/push_dev' >> "$project_root/notifications/.env"
    load_environment notifications
    [[ "$DATABASE_URL" == 'postgres://user:password@localhost:5432/push_dev' ]]
}

test_query_string_with_slashes() {
    reset_fixture
    printf '%s\n' 'DATABASE_URL=postgres://user:password@localhost/casita?application_name=local/dev' >> "$project_root/api/.env"
    load_environment notifications
    [[ "$DATABASE_URL" == 'postgres://user:password@localhost/notifications?application_name=local/dev' ]]
}

test_shared_secret_mismatch() {
    reset_fixture
    printf '%s\n' 'S2S_TOKEN=do-not-print-this-value' >> "$project_root/notifications/.env"
    if load_environment notifications > "$fixture_root/error" 2>&1; then
        printf 'Expected shared secret mismatch to fail.\n' >&2; return 1
    fi
    [[ $(< "$fixture_root/error") == *'S2S_TOKEN must be configured identically'* ]]
    [[ $(< "$fixture_root/error") != *'do-not-print-this-value'* ]]
}

test_same_database_rejected() {
    reset_fixture
    printf '%s\n' 'DATABASE_URL=postgres://other:password@localhost/casita' >> "$project_root/notifications/.env"
    if load_environment notifications > "$fixture_root/error" 2>&1; then
        printf 'Expected identical database targets to fail.\n' >&2; return 1
    fi
    [[ $(< "$fixture_root/error") == *'must use different databases'* ]]
}

test_notification_env_required() {
    reset_fixture
    rm "$project_root/notifications/.env"
    if load_environment notifications > "$fixture_root/error" 2>&1; then
        return 1
    fi
    [[ $(< "$fixture_root/error") == *'Create notifications/.env'* ]]
}

test_notification_keys_required() {
    local missing_key
    for missing_key in VAPID_PUBLIC_KEY VAPID_PRIVATE_KEY; do
        reset_fixture
        sed -i "/^${missing_key}=/d" "$project_root/notifications/.env"
        export VAPID_PUBLIC_KEY=stale-public VAPID_PRIVATE_KEY=stale-private
        if load_environment notifications > "$fixture_root/error" 2>&1; then
            return 1
        fi
        [[ $(< "$fixture_root/error") == *"$missing_key is required in notifications/.env."* ]]
    done
}

test_existing_session_reused() {
    reset_fixture
    tmux() { [[ "$*" == 'has-session -t =casita-dev' ]]; }
    docker() { printf 'Docker must not run when reusing a session.\n' >&2; return 1; }
    start_services --detach > "$fixture_root/output"
    [[ $(< "$fixture_root/output") == *'Attach with: tmux attach -t casita-dev'* ]]
}

for test_case in test_api_environment test_notification_defaults test_notification_database_override \
    test_notification_env_required test_notification_keys_required \
    test_query_string_with_slashes test_shared_secret_mismatch test_same_database_rejected test_existing_session_reused; do
    ("$test_case")
    printf 'PASS %s\n' "$test_case"
done
