---
area: dx
phase: pre-v0.1
---
# Dev watcher restart kill timeout

## Goal
`restartServer` must not hang the watcher when a server ignores SIGTERM.

## Acceptance criteria
- `restartServer` has a kill timeout.
- A watcher test covers the timeout.
