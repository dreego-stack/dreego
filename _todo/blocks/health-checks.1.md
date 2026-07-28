---
id: health-checks.1
title: /health + /ready Endpoints
status: planned
phase: v0.0.x
requires:
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Fixed: /health (Liveness) + /ready (Readiness, configurable). Config-based activation. K8s/Docker-compatible. Plugins can register readiness checks (DB, Redis). JSON response with status.
