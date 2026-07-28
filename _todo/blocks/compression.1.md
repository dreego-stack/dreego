---
id: compression.1
title: Gzip/Brotli Compression Middleware
status: planned
phase: v0.0.x
requires:
  - middleware.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Conditional Middleware for response compression. Gzip (standard) + Brotli (optional, modern). Content-Type-based selection. No external package — evaluate Go's compress/gzip + andybalholm/brotli.
