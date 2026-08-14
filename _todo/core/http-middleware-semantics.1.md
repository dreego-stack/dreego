---
area: runtime
phase: v0.1-blocker
depends_on: app-runtime.1
---
# Make recovery and compression HTTP-correct

## Goal
Ensure middleware composition cannot emit corrupt, mislabeled, or invalid responses.

## Acceptance criteria
- A panic under compression produces one valid response with the intended status and encoding.
- Compression respects `gzip;q=0`, sets `Vary`, and skips HEAD, 204, 304, and pre-encoded responses.
- `Content-Length` and partial-write behavior are correct.
- Response-writer wrappers preserve supported optional interfaces or expose deliberate limitations.
- Stack-level HTTP tests exercise recovery and compression together.
