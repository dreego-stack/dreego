---
id: csrf.1
title: CSRF Protection (Core-Conditional)
status: 15
phase: v0.0.3
requires:
  - session.1
  - middleware.1
created: 2026-07-26
changed: 2026-07-27
---

CSRF token via session. Middleware validates token on POST/PUT/DELETE. Auto-injection into <form> tags via Codegen. No external backend — session-based (Double Submit Cookie). Core-Conditional: default on, opt-out via config.json.
