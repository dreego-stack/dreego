---
id: request-id.1
title: Request-ID Middleware (X-Request-ID)
status: 34
phase: v0.0.17
requires:
  - middleware.1
created: 2026-07-29
changed: 2026-07-29
---

Core-Fixed middleware. Every request gets a unique X-Request-ID header (16-char hex). Client-supplied IDs are accepted. The ID is injected into the request context and included in JSONL log output via `rid` field. Available via `c.RequestID()` in templates and Go blocks.
