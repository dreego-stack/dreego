---
area: dx
phase: pre-v0.1
---
# Startup configuration warnings

## Goal
Warn in `Build()` when CSRF is enabled but no session store is configured.
Log a warning when `dreego/config.json` is invalid (`loadSettings` is currently
silent).

## Acceptance criteria
- `Build()` warns when CSRF is enabled without a session store.
- An invalid `dreego/config.json` produces a warning.
- Both warnings are covered by tests.
