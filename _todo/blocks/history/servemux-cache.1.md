---
id: servemux-cache.1
title: Cache Built Middleware/Router Stack
status: 38
phase: v0.0.22
requires:
  - middleware.1
  - routing.1
created: 2026-07-31
changed: 2026-08-08
---

`core.ServeMux()` currently rebuilds the entire middleware and routing stack on every call. Build the handler once at first use (or during `core.Listen`) and cache it. Improves performance and makes the stack inspectable.

Done in v0.0.22: `core/runtime.go` caches the built middleware/router stack.
