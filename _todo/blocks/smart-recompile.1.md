---
id: smart-recompile.1
title: Smart Recompile (Text-vs-Go Detection)
status: planned
phase: v0.1.0
requires:
  - transpiler.1
  - hot-reload.1
created: 2026-07-27
changed: 2026-07-27
---

Templ-inspiriert: Generator erkennt, ob sich nur HTML-Text oder Go-Expressions geändert haben. Text-Only: Runtime liest aktualisierte Literale aus `_dreego.txt` ohne `go build`. Go-Changes: volles Recompile. Ermöglicht <1s Hot-Reload für reine Template-Änderungen. Kern-Mechanismus aus templs Watch-Mode.
