---
id: security-cookie.1
title: Harden Session and CSRF Cookie Flags
status: 35
phase: pre-v1.0
requires:
  - session.1
  - csrf.1
created: 2026-07-31
changed: 2026-07-31
---

Review and tighten cookie security defaults:
- Set Secure flag on session cookie when TLS is active
- Review SameSite defaults for both session and csrf_token cookies
- Ensure HttpOnly is applied correctly per cookie purpose

Keep behavior core-conditional and configurable via `core.Options` or config.json.
