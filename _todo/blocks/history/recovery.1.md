---
id: recovery.1
title: Recovery-Middleware (Panic → 500)
status: 8
phase: v0.0.2
requires:
  - middleware.1
created: 2026-07-25
changed: 2026-07-26
---

recover() in Core-Fixed middleware. Loggt Panic mit Stack-Trace via slog.NewJSONHandler. Optionaler 500-Error-Handler via runtime.SetErrorHandler.
