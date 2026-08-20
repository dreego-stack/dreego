---
area: tests
phase: pre-v0.1
---
# Coverage measurement and gate

## Goal
Introduce coverage measurement (`go test -cover ./core/...`) and a coverage
gate script wired into `Makefile`/`test.sh` with a threshold.

## Acceptance criteria
- A coverage script runs in CI.
- A threshold is defined.
- The first run reports a baseline number.
- README/drunk not affected.
