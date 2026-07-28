---
id: storage-interface.1
title: File Storage Interface (S3, R2, Local)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core interface for file storage. Abstracts S3/R2/Local. Methods: Put, Get, Delete, List, URL. Plugin implementations (dreego-storage). No core provider — interface only, like database/sql.
