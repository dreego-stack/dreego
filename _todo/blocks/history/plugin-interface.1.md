---
id: plugin-interface.1
title: Plugin Interface (Frozen for v1)
status: 44
phase: v0.0.25
requires:
  - context-refactoring.1
  - middleware.1
created: 2026-07-26
changed: 2026-08-08
---

Finalize capability-based plugin system. Plugin interface: OnStart, OnShutdown, Middleware, Routes. Validate plugin configuration. ADR plugin-interface.md exists. First Release = Final Contract.

Done in v0.0.25: frozen v1 `core.Plugin` contract (`Name`, `RegisterRoutes`, `Middlewares`, `Assets`, `OnStart`, `OnShutdown`), `core.UsePlugin(p)`.
