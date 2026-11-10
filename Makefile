.PHONY: dev dev-start dev-stop

PROJECT_ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Start the development stack and attach (or switch) to its tmux session.
dev:
	@bash "$(PROJECT_ROOT)dev.sh" start

# Start without attaching, useful from scripts or another terminal.
dev-start:
	@bash "$(PROJECT_ROOT)dev.sh" start --detach

# Stop application processes; keep PostgreSQL and its data.
dev-stop:
	@bash "$(PROJECT_ROOT)dev.sh" stop
