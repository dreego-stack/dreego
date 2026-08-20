---
area: security
phase: pre-v0.1
---
# Bind wraps r.Body in http.MaxBytesReader with a configurable limit

## Goal
`Bind` must wrap `r.Body` in `http.MaxBytesReader` with a configurable limit.

## Acceptance criteria
- `Bind` reads the body through `http.MaxBytesReader`.
- The limit is configurable.
- An oversized JSON POST returns a 413-ish error.
- A test covers the oversized body path.
