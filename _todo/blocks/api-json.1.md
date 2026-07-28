---
id: api-json.1
title: API Routes + JSON Responses
status: planned
phase: v0.0.x
requires:
  - routing.1
  - context-refactoring.1
created: 2026-07-26
changed: 2026-07-26
---

`dreego/routes/api/` — API routes with JSON response. `c.JSON(data)` + `c.JSONError(code, msg)`. Content-Type auto-switch. API routes without layout (JSON only). Struct-to-JSON via encoding/json. Method-based Routing: POST → Request-Body Parsing.
