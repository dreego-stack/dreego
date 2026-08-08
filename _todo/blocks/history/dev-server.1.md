---
id: dev-server.1
title: Dev Server with Hot Reload
status: 49
phase: v0.0.25
requires:
  - transpiler.1
  - cli.1
created: 2026-07-31
changed: 2026-08-08
---

`dreego dev` command with file watcher, auto-regenerate, and server restart. Provides first-class DX without requiring external tools like Air. Replaces the rejected hot-reload.1/live-reload.1 blocks with a built-in solution.

Done in v0.0.25: `dreego dev` watches `.dreego` files, regenerates + rebuilds, graceful restart.
