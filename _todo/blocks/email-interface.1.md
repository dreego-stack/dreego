---
id: email-interface.1
title: Email Sending Interface (SMTP, Resend, Postmark)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core interface for email sending. Abstracts SMTP/Resend/Postmark. Methods: Send(To, Subject, Body). Template-based. Plugin implementations (dreego-mail). Queue-capable via queue-interface.1.
