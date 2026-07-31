---
id: security-session.1
title: Document or Encrypt Session Payload
status: planned
phase: pre-v1.0
requires:
  - session.1
created: 2026-07-31
changed: 2026-07-31
---

Session payload is currently HMAC-signed but readable. Document this clearly or add optional encryption for sensitive values. Encryption must remain optional to keep the core simple.
