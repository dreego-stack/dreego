---
id: api-json.1
title: API-Routen + JSON Responses
status: planned
phase: v0.0.x
requires:
  - routing.1
  - context-refactoring.1
created: 2026-07-26
changed: 2026-07-26
---

`dreego/routes/api/` — API-Routen mit JSON-Response. `c.JSON(data)` + `c.JSONError(code, msg)`. Content-Type auto-switch. API-Routen ohne Layout (nur JSON). Struct-to-JSON via encoding/json. Method-based Routing: POST → Request-Body Parsing.
