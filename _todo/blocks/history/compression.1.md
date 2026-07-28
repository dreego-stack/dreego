---
id: compression.1
title: Gzip Compression Middleware
status: 31
phase: v0.0.14
requires:
  - middleware.1
created: 2026-07-26
changed: 2026-07-29
---

Core-Conditional Middleware for response compression. Gzip via compress/gzip wrapping ResponseWriter. Accept-Encoding header check. Applied after security headers, before logging.
