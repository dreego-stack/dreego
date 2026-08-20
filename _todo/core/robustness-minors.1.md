---
area: code
phase: pre-v0.1
---
# Minor robustness fixes

## Goal
- Add `responseWriter.Unwrap()` for `http.ResponseController`.
- `BindForm` returns an error instead of panicking on a non-pointer.
- (Optional) meta-description dedupe attribute-order-insensitive.

## Acceptance criteria
- `responseWriter` supports `http.ResponseController` via `Unwrap()`.
- `BindForm` returns an error on a non-pointer argument instead of panicking.
- Tests cover the `Unwrap()` and `BindForm` error paths.
