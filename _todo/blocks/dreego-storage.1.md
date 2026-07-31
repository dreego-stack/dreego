---
id: dreego-storage.1
title: dreego-storage (File Uploads, Progress, Resize)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - storage-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Target directory: `plugins/storage/` in this repository.

Plugin for file uploads. Implements storage-interface.1. S3, R2, Local backends. Progress tracking via SSE. Image resize/thumbnail generation. Integration with dreego-stripe (invoice uploads).
