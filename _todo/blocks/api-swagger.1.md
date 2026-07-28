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

Auto-generated OpenAPI 3.0 Spec from Go struct tags and API routes. `c.Swagger()` endpoint. Struct-Tags: `api:"description"`, `validate:"required"`. CodeGen generates /openapi.json route. Swagger UI optionally embeddable. No manual spec file needed.
