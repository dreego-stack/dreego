---
id: route-hooks.1
title: Plugin Route Registration
status: 46
phase: v0.0.25
requires:
  - plugin-interface.1
  - routing.1
created: 2026-07-26
changed: 2026-08-08
---

Core entry point for plugin routes. Plugins register own URL paths (e.g. /admin/*, /api/auth/*). No filesystem-based registration — programmatic via Plugin API. Generated dreego/gen/dree.go collects all plugin routes.

Done in v0.0.25: plugin calls `core.Register(...)` in `RegisterRoutes()`, reachable via `core.ServeMux()`.
