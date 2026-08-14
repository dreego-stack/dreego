---
area: runtime
phase: v0.1-blocker
depends_on: app-runtime.1
---
# Wait for graceful shutdown to finish

## Goal
Keep the process alive until active requests drain or the shutdown deadline is reached.

## Acceptance criteria
- `Listen` waits for shutdown completion and propagates shutdown failures.
- Signal subscriptions are released and repeated server lifecycles do not leak goroutines.
- An in-flight request completes before process exit when it fits within the deadline.
- A request exceeding the deadline produces a deterministic, observable failure.
