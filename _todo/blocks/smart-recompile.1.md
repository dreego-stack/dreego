---
id: smart-recompile.1
title: Smart Recompile (Text vs Go Detection)
status: planned
phase: v0.1.0
requires:
  - transpiler.1
  - hot-reload.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspired: Generator detects whether only HTML text or Go expressions changed. Text-only: Runtime reads updated literals from `_dreego.txt` without `go build`. Go changes: full recompile. Enables <1s hot reload for pure template changes. Core mechanism from templ's watch mode.
