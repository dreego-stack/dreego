---
id: dreego-mail.1
title: dreego-mail (Email SMTP/Resend/Postmark)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - email-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Target directory: `plugins/mail/` in this repository.

Plugin for email sending. Implements email-interface.1. SMTP, Resend, Postmark backends. HTML/Text templates. Queue integration for async sending. Usable by dreego-auth (Verify, Reset) and dreego-notify.
