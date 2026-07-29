---
id: request-id.1
title: Request-ID Middleware (X-Request-ID)
status: planned
phase: v0.0.17
requires:
  - middleware.1
created: 2026-07-29
changed: 2026-07-29
---

Core-Fixed middleware. Every request gets a unique `X-Request-ID` header. If client sends one, it's accepted. Otherwise, a UUIDv4 is generated. The ID is injected into the request context and included in all log output. Available via `c.Get("request_id")` in templates and Go blocks.
