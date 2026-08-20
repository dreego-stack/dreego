---
area: dx
phase: pre-v0.1
---
# Dev watcher skips and restart timeout

## Goal
The dev watcher skips `.git`, `node_modules`, and build directories. Add a kill
timeout to `restartServer`.

## Acceptance criteria
- The watcher skips `.git`, `node_modules`, and build dirs.
- `restartServer` has a kill timeout.
- A watcher test covers the skips.
