---
id: security-headers.1
title: Security Headers (nosniff, frame, referrer, permissions)
status: 30
phase: v0.0.14
requires:
  - middleware.1
created: 2026-07-26
changed: 2026-07-29
---

Core-Conditional Middleware. X-Content-Type-Options: nosniff, X-Frame-Options: DENY, Referrer-Policy: strict-origin-when-cross-origin, Permissions-Policy: none. Default restrictive, applied before compression in middleware chain.
