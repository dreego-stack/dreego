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

dreego/static/ folder: Files are read during generate and registered as inline `[]byte` in `core.RegisterStatic(path, mime, bytes)`. MIME via extension. Collision check: same URL path as route → Error.
