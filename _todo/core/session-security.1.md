---
area: security
phase: v0.1-blocker
---
# Authentication-state invalidation for cookie sessions

## Goal
Define how stateless cookie sessions invalidate or replace authentication state
across login, privilege changes, logout, key rotation, and replay scenarios.

## Acceptance criteria
- The documentation does not promise identifier rotation when the selected store has no server-side identifier.
- Auth plugins can replace or invalidate pre-authentication state on login and privilege changes.
- Logout invalidates the complete authentication state rather than deleting one value with weaker cookie options.
- Key rotation, replay limitations, and store-specific guarantees are documented.
- Tests cover login, privilege changes, logout, replayed cookies, rotated keys, and invalid sessions.
