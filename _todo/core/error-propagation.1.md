---
area: correctness
phase: v0.1-blocker
---
# Propagate boundary errors

## Goal
Stop treating failed file access, parsing, form handling, session operations, and generated rendering as successful empty results.

## Acceptance criteria
- Filesystem walk, directory read, and source read failures abort generation with the affected path and wrapped cause.
- Form parsing and binding failures remain distinguishable from valid empty values.
- Session reads, writes, deletes, destruction, and CSRF persistence failures reach an application error path.
- Generated render calls do not discard returned errors.
- Public APIs return errors where callers can recover; unrecoverable startup errors fail loudly.
- Production responses use generic error messages while internal causes remain available to application error handling and structured logs.
- Tests prove that filesystem paths, database details, Go type errors, and other internal causes are not disclosed to clients.
- Regression tests cover each previously silent failure path.
