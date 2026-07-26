---
id: middleware-hooks.1
title: Plugin-Middleware-Hooks (app.Use FIFO)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - middleware.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Eintrittspunkt für Plugin-Middleware. app.Use() FIFO-Registrierung. Plugins injizieren Middleware in den Core-Stack. Reihenfolge-Fixierung beim ersten ListenAndServe. Constraint-Sortierung (before/after) in V2.
