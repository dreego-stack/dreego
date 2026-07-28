---
id: ddos-protection.1
title: DDoS Protection (PoW + Rate-Limiting) — Plugin
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - middleware-hooks.1
created: 2026-07-26
changed: 2026-07-26
---

Plugin (not Core) — needs Redis/Backend. Rate limiting per IP/User. Proof-of-Work challenge for critical endpoints. Token bucket algorithm. Config-based thresholds. DDoS = infrastructure-dependent → Plugin per MIDDLEWARE_SYSTEM decision.
