---
id: middleware-hooks.1
title: Plugin Middleware Hooks (app.Use FIFO)
status: 45
phase: v0.0.25
requires:
  - plugin-interface.1
  - middleware.1
created: 2026-07-26
changed: 2026-08-08
---

Core entry point for plugin middleware. app.Use() FIFO registration. Plugins inject middleware into the core stack. Order fixated on first ListenAndServe. Constraint sorting (before/after) in V2.

Done in v0.0.25: plugin `Middlewares()` appended FIFO to runtime chain.
