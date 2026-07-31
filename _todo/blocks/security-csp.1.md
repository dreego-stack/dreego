---
id: security-csp.1
title: Add Content-Security-Policy Header
status: planned
phase: pre-v1.0
requires:
  - security-headers.1
created: 2026-07-31
changed: 2026-07-31
---

Add a sensible Content-Security-Policy default to the security headers middleware. Must allow HTMX/Alpine.js and inline styles from scoped CSS. Should be overridable via dreego/config.json.
