---
id: hot-reload.1
title: Hot Reload (Dev Server + SSE)
status: planned
phase: v0.0.x
requires:
  - cli.1
  - routing.1
created: 2026-07-26
changed: 2026-07-26
---

File change → dreego generate → Browser reload via SSE. State preservation across reloads. Only re-render affected components. File watcher in Go (fsnotify) or via CLI polling.
