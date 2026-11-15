#!/usr/bin/env bash
set -euo pipefail

project_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
session_name=casita-dev

fail() { printf 'dev: %s\n' "$*" >&2; return 1; }
session_exists() { tmux has-session -t "=$session_name" 2>/dev/null; }

host_database_url() {
    local url=$1 scheme rest authority hostport
    [[ "$url" == postgres://*/* || "$url" == postgresql://*/* ]] ||
        { fail 'DATABASE_URL must be a PostgreSQL URL with a database name.'; return 1; }
    scheme=${url%%://*}
    rest=${url#*://}
    authority=${rest%%/*}
    hostport=${authority##*@}
    case "$hostport" in
        postgres|postgres:*) authority="${authority%"$hostport"}127.0.0.1${hostport#postgres}" ;;
    esac
    printf '%s://%s/%s\n' "$scheme" "$authority" "${rest#*/}"
}

database_target() {
    local rest=${1#*://} hostport database
    hostport=${rest%%/*}
    hostport=${hostport##*@}
    database=${rest#*/}
    database=${database%%\?*}
    if [[ "$hostport" =~ ^(localhost|127\.0\.0\.1|\[::1\])(:([0-9]+))?$ ]]; then
        hostport="localhost:${BASH_REMATCH[3]:-5432}"
    elif [[ "$hostport" != *:* ]]; then
        hostport+=:5432
    fi
    printf '%s/%s\n' "$hostport" "$database"
}

load_environment() {
    local service=$1 api_database_url api_s2s_token key
    [[ -f "$project_root/api/.env" ]] ||
        { fail 'Create api/.env from api/.env.example first.'; return 1; }
    set -a
    # These are trusted local shell assignment files, like the existing Make targets.
    source "$project_root/api/.env"
    set +a
    [[ -n "${DATABASE_URL:-}" ]] || { fail 'DATABASE_URL is required in api/.env.'; return 1; }
    DATABASE_URL=$(host_database_url "$DATABASE_URL")
    api_database_url=$DATABASE_URL
    # VAPID belongs to notifications, including when inherited from the shell.
    unset VAPID_PUBLIC_KEY VAPID_PRIVATE_KEY
    [[ -n "${S2S_TOKEN:-}" ]] || { fail 'S2S_TOKEN is required in api/.env.'; return 1; }
    api_s2s_token=$S2S_TOKEN
    [[ ${#S2S_TOKEN} -ge 32 ]] || { fail 'S2S_TOKEN must be at least 32 characters long.'; return 1; }
    export ENV=dev PORT=2345 NOTIFICATIONS_SERVICE_URL=http://localhost:2346
    if [[ "$service" == notifications ]]; then
        # Keep the connection authority and query string; change only the database.
        DATABASE_URL=${api_database_url%%\?*}
        DATABASE_URL="${DATABASE_URL%/*}/notifications"
        [[ "$api_database_url" != *\?* ]] || DATABASE_URL+="?${api_database_url#*\?}"
        [[ -f "$project_root/notifications/.env" ]] ||
            { fail 'Create notifications/.env from notifications/.env.example first.'; return 1; }
        set -a
        source "$project_root/notifications/.env"
        set +a
        for key in VAPID_PUBLIC_KEY VAPID_PRIVATE_KEY; do
            [[ -n "${!key:-}" ]] || { fail "$key is required in notifications/.env."; return 1; }
        done
        DATABASE_URL=$(host_database_url "$DATABASE_URL")
        [[ $(database_target "$DATABASE_URL") != "$(database_target "$api_database_url")" ]] ||
            { fail 'API and notifications must use different databases.'; return 1; }
        [[ "${S2S_TOKEN:-}" == "$api_s2s_token" ]] ||
            { fail 'S2S_TOKEN must be configured identically for both services.'; return 1; }
        export ENV=dev PORT=2346
    fi
    export DATABASE_URL
}

find_air() {
    local gobin gopath directory
    local go_paths=()
    if command -v air >/dev/null; then command -v air; return; fi
    gobin=$(go env GOBIN)
    if [[ -n "$gobin" && -x "$gobin/air" ]]; then printf '%s\n' "$gobin/air"; return; fi
    gopath=$(go env GOPATH)
    IFS=: read -r -a go_paths <<< "$gopath"
    for directory in "${go_paths[@]}"; do
        if [[ -x "$directory/bin/air" ]]; then printf '%s\n' "$directory/bin/air"; return; fi
    done
    fail 'Air is required: go install github.com/air-verse/air@latest'
}

show_session() {
    printf 'Development session: %s. UI: http://localhost:5173\n' "$session_name"
    if [[ "${1:-}" == --detach || ! -t 0 || ! -t 1 ]]; then
        printf 'Attach with: tmux attach -t %s\n' "$session_name"
    elif [[ -n "${TMUX:-}" ]]; then
        tmux switch-client -t "$session_name"
    else
        tmux attach-session -t "$session_name"
    fi
}

start_services() {
    if session_exists; then show_session "${1:-}"; return; fi
    local dependency port service service_command
    local cache_env=()
    for dependency in docker go npm; do
        command -v "$dependency" >/dev/null || { fail "$dependency is required."; return 1; }
    done
    (load_environment notifications)
    find_air >/dev/null
    [[ -d "$project_root/ui/node_modules" ]] ||
        { fail 'Install UI dependencies first: npm --prefix ui ci'; return 1; }
    for port in 2345 2346 5173; do
        if (exec 3<>"/dev/tcp/127.0.0.1/$port") 2>/dev/null; then
            fail "Port $port is already in use outside the $session_name session."
            return 1
        fi
    done
    docker compose -f "$project_root/compose.yml" --project-directory "$project_root" up -d --wait postgres
    [[ -z "${GOCACHE:-}" ]] || cache_env=(-e "GOCACHE=$GOCACHE")
    for service in api notifications ui; do
        if [[ "$service" == api ]]; then
            tmux new-session -d -s "$session_name" -n "$service" -c "$project_root" "${cache_env[@]}"
        else
            tmux new-window -d -t "$session_name:" -n "$service" -c "$project_root"
        fi
        tmux set-window-option -t "$session_name:$service" remain-on-exit on
        printf -v service_command '%q ' bash "$project_root/dev.sh" service "$service"
        tmux respawn-pane -k -t "$session_name:$service" "$service_command"
    done
    tmux select-window -t "$session_name:ui"
    show_session "${1:-}"
}

stop_services() {
    if ! session_exists; then printf 'No %s session is running.\n' "$session_name"; return; fi
    local pane dead states attempt
    while read -r pane dead; do
        [[ "$dead" == 1 ]] || tmux send-keys -t "$pane" C-c
    done < <(tmux list-panes -s -t "$session_name" -F '#{pane_id} #{pane_dead}')
    # Terminal interrupts reach the whole foreground process group, including
    # go run/npm children. Air handles its own child process groups on shutdown.
    for ((attempt=0; attempt<80; attempt++)); do
        session_exists || break
        states=$(tmux list-panes -s -t "$session_name" -F '#{pane_dead}')
        [[ "$states" == *0* ]] || break
        sleep 0.1
    done
    if session_exists; then tmux kill-session -t "$session_name"; fi
    printf 'Development servers stopped. PostgreSQL is still running; data is preserved.\n'
}

run_service() {
    case "$1" in
        api) load_environment api; cd "$project_root/api"; exec "$(find_air)" ;;
        notifications) load_environment notifications; cd "$project_root/notifications"; exec go run ./cmd ;;
        ui) cd "$project_root/ui"; exec npm run dev -- --host 127.0.0.1 --port 5173 --strictPort ;;
        *) fail 'Unknown service.' ;;
    esac
}

main() {
    case "${1:-}" in
        start)
            [[ $# == 1 || ( $# == 2 && $2 == --detach ) ]] || { fail 'Usage: dev.sh start [--detach]'; return 1; }
            command -v tmux >/dev/null || { fail 'tmux is required.'; return 1; }
            start_services "${2:-}" ;;
        stop)
            [[ $# == 1 ]] || { fail 'Usage: dev.sh stop'; return 1; }
            command -v tmux >/dev/null || { fail 'tmux is required.'; return 1; }
            stop_services ;;
        service)
            [[ $# == 2 ]] || { fail 'Usage: dev.sh service api|notifications|ui'; return 1; }
            run_service "$2" ;;
        *) fail 'Usage: dev.sh start [--detach] | stop | service api|notifications|ui' ;;
    esac
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then main "$@"; fi
