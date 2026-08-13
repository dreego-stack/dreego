---
area: security
phase: v0.1-blocker
depends_on: app-runtime.1
---
# HTTP server limits and timeouts

## Goal
Protect the default server from slow requests and unbounded request bodies.

## Acceptance criteria
- The server configures read-header, read, write, and idle timeouts.
- `MaxHeaderBytes` and request-body limits have documented defaults.
- `App` provides secure defaults for every route.
- Connection and header timeouts remain app-wide server policy for v0.1.
- Request-body limits can be narrowed or deliberately raised for a specific route through normal handler composition without weakening unrelated routes.
- Streaming and upload exceptions are designed only when a real plugin proves the requirement.
- Documentation explains the consequences of disabling or raising a limit.
