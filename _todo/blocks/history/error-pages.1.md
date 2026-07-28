---
id: error-pages.1
title: Custom Error Pages (404.dreego + 500.dreego)
status: 10
phase: v0.0.2
requires:
  - routing.1
  - recovery.1
created: 2026-07-25
changed: 2026-07-26
---

404.dreego as catch-all per directory (Go Mux pattern precedence). 500.dreego via Recovery middleware. No layout wrapping. Per-directory cascading.
