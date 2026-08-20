---
area: api
phase: pre-v0.1
---
# App.Build() returns an error instead of panicking

## Goal
`App.Build()` must return an `error` instead of panicking on invalid
configuration (session store validation, conflicting route patterns).

## Acceptance criteria
- `Build()` returns an error for invalid session store configuration.
- `Build()` returns an error for conflicting/overlapping route patterns.
- The error is propagated through `Listen()` and `Handler()`.
- No panic paths remain in `Build()` for the covered invalid inputs.
- Tests cover a bad session store and overlapping route patterns.
