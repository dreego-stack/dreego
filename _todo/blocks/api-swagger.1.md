---
id: api-swagger.1
title: Swagger/OpenAPI Auto-Generation
status: planned
phase: v0.0.x
requires:
  - api-json.1
created: 2026-07-26
changed: 2026-07-26
---

Auto-generierte OpenAPI 3.0 Spec aus Go-Struct-Tags und API-Routen. `c.Swagger()` endpoint. Struct-Tags: `api:"description"`, `validate:"required"`. CodeGen erzeugt /openapi.json Route. Swagger UI optional einbindbar. Keine manuelle Spec-Datei nötig.
