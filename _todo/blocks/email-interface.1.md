---
id: email-interface.1
title: Email-Sending Interface (SMTP, Resend, Postmark)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Interface für E-Mail-Versand. Abstrahiert SMTP/Resend/Postmark. Methoden: Send(To, Subject, Body). Template-basiert. Plugin-Implementierungen (dreego-mail). Queue-fähig via queue-interface.1.
