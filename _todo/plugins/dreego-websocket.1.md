---
area: plugins
phase: future
---
# WebSocket plugin

## Goal
Provide WebSocket upgrades and connection lifecycle management without adding a WebSocket dependency or protocol-specific state to core.

## Acceptance criteria
- The plugin lives in a separate repository with its own `go.mod`, tests, releases, and CI.
- It registers only the routes and middleware it owns on a specific `App`.
- Connection limits, shutdown, origin checks, authentication hooks, and backpressure are documented and tested.
- Any route-level server-limit exceptions do not weaken unrelated application routes.
