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

Core-Fixed: /health (Liveness) + /ready (Readiness, konfigurierbar). Config-basiert aktivierbar. K8s/Docker-kompatibel. Plugins können Readiness-Checks registrieren (DB, Redis). JSON-Response mit Status.
