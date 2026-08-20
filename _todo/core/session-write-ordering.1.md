---
area: runtime
phase: pre-v0.1
---
# CookieStore.Set must not mutate state before the cookie write succeeds

## Goal
`CookieStore.Set` must not mutate request-local session state before the
cookie write succeeds.

## Acceptance criteria
- On `ErrSessionTooLarge`, `SessionVal` does not return the new value.
- The session state is only updated after the cookie write succeeds.
- A test covers an oversize cookie write and verifies the value is unchanged.
