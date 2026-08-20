---
area: tests
phase: pre-v0.1
---
# Extend dreegotest with assertion and session helpers

## Goal
Extend `dreegotest` with:
- Assertion helpers: `MustStatus`, `MustContainBody`, `MustHeader`,
  `MustEqual`/`MustNotEqual`.
- A typed `Response` struct with headers/cookies.
- Session helpers: `SessionVal`, CSRF token helper.
- `MustScaffold(t)` wrapping the init+generate+build ritual.
- Export `freePort`/`waitForPort` (rename usage in `quick_start_test.go`).

## Acceptance criteria
- No external deps (`go.mod` stays zero-dep; no `x/net/html`).
- `dreegotest` self-tests pass.
- At least 5 existing integration tests are refactored to the new helpers.
- `make test` is green.
