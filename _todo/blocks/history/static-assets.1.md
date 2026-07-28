---
id: static-assets.1
title: Static Assets (dreego/static/ → inline Handler)
status: 26
phase: v0.0.10
requires:
  - transpiler.1
  - routing.1
created: 2026-07-26
changed: 2026-07-28
---

dreego/static/ Ordner: Files werden beim Generate eingelesen und als inline `[]byte` in `core.RegisterStatic(path, mime, bytes)` registriert. MIME via Extension. Kollision-Check: gleicher URL-Path wie Route → Error.
