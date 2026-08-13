---
area: plugins
phase: after-v0.1
---
# OpenAPI plugin

## Goal
Generate an OpenAPI 3 specification from typed route metadata and optionally serve Swagger UI.

## Acceptance criteria
- The implementation lives in a separate plugin module.
- The plugin serves `/openapi.json` from typed route metadata.
- Swagger UI remains optional.
