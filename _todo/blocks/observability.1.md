---
id: observability.1
title: Observability (Prometheus, OpenTelemetry, Request-ID)
status: planned
phase: v0.0.x
requires:
  - middleware.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Fixed: Request-ID (X-Request-ID header, in Context + Logs). Core-Conditional: /metrics (Prometheus format). Plugin: OpenTelemetry tracing. Structured logging via slog. Every request gets trace_id + request_id.
