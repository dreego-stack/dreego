---
id: storage-interface.1
title: File-Storage Interface (S3, R2, Local)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Interface für File-Storage. Abstrahiert S3/R2/Local. Methoden: Put, Get, Delete, List, URL. Plugin-Implementierungen (dreego-storage). Kein Core-Provider — Interface only, wie database/sql.
