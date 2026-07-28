---
id: route-hooks.1
title: Plugin Route Registration
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

Core entry point for plugin routes. Plugins register own URL paths (e.g. /admin/*, /api/auth/*). No filesystem-based registration — programmatic via Plugin API. Generated dreego/gen/dree.go collects all plugin routes.
