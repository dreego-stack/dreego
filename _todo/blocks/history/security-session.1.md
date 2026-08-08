---
id: security-session.1
title: Document or Encrypt Session Payload
status: 40
phase: v0.0.22
requires:
  - session.1
created: 2026-07-31
changed: 2026-08-08
---

Session payload is currently HMAC-signed but readable. Document this clearly or add optional encryption for sensitive values. Encryption must remain optional to keep the core simple.

Done in v0.0.22: optional AES-256-GCM session encryption (`core.Options{Encrypt: true}`), encrypt-then-MAC.
