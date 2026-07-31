---
id: servemux-cache.1
title: Cache Built Middleware/Router Stack
status: planned
phase: pre-v1.0
requires:
  - middleware.1
  - routing.1
created: 2026-07-31
changed: 2026-07-31
---

`core.ServeMux()` currently rebuilds the entire middleware and routing stack on every call. Build the handler once at first use (or during `core.Listen`) and cache it. Improves performance and makes the stack inspectable.
