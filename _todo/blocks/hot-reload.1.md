---
id: hot-reload.1
title: Hot Reload (Dev-Server + SSE)
status: planned
phase: v0.0.x
requires:
  - cli.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

Datei-Änderung → dreego generate → Browser Reload via SSE. State-Erhalt über Reloads. Nur betroffene Komponenten neu rendern. File-Watcher in Go (fsnotify) oder via CLI-Polling.
