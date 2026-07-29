---
id: observability.1
title: Observability (Request-ID, Metrics, Tracing)
status: planned
phase: v0.0.x
requires:
  - middleware.1
created: 2026-07-26
changed: 2026-07-29
---

Split into Core + Plugins:
- **Core-Fixed** (`request-id.1`): X-Request-ID header middleware, inject into context + logs
- **Plugin** (`dreego-metrics`): Prometheus /metrics endpoint — blocked on plugin-interface.1
- **Plugin** (`dreego-tracing`): OpenTelemetry spans — V2, blocked on plugin-interface.1
