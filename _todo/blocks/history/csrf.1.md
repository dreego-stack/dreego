---
id: csrf.1
title: CSRF-Schutz (Core-Conditional)
status: 15
phase: v0.0.3
requires:
  - session.1
  - middleware.1
created: 2026-07-26
changed: 2026-07-27
---

CSRF-Token via Session. Middleware validiert Token bei POST/PUT/DELETE. Auto-Injection in <form>-Tags per Codegen. Kein externes Backend — Session-basiert (Double Submit Cookie). Core-Conditional: default an, opt-out via config.json.
