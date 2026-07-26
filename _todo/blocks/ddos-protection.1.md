---
id: ddos-protection.1
title: DDoS-Schutz (PoW + Rate-Limiting) — Plugin
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - middleware-hooks.1
created: 2026-07-26
changed: 2026-07-26
---

Plugin (nicht Core) — braucht Redis/Backend. Rate-Limiting pro IP/User. Proof-of-Work Challenge für kritische Endpoints. Token-Bucket Algorithm. Config-basierte Schwellwerte. DDoS = Infrastruktur-abhängig → Plugin nach MIDDLEWARE_SYSTEM Decision.
