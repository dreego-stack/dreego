---
id: health-checks.1
title: /health + /ready Endpoints
status: 29
phase: v0.0.14
requires:
  - routing.1
created: 2026-07-26
changed: 2026-07-29
---

Core-Fixed: /health (Liveness, returns "ok") + /ready (Readiness, returns "ok"/"not ready"). Configurable via core.SetReady(bool). Plugins can register readiness checks (DB, Redis). JSON response with status.
