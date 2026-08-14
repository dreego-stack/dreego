---
area: hardening
phase: v0.1-blocker
depends_on:
  - app-runtime.1
  - session-security.2
---
# Race detector and global-state audit

## Goal
Detect races in routes, sessions, CSP configuration, plugins, and other mutable package state.

## Acceptance criteria
- CI runs the relevant tests with `-race`.
- Detected shared-state races are fixed or removed by design.
