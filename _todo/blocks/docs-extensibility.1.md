---
id: docs-extensibility.1
title: Extensible dreego docs Command
status: planned
phase: v0.x.0
requires:
  - cli.1
  - plugin-interface.1
created: 2026-07-31
changed: 2026-07-31
---

Adapt `dreego docs` so it can consume documentation from plugins and external repos, not only this repository. Open design question: how to discover plugin docs elegantly (embed? local path? Codeberg URL?). Must not depend on any specific plugin at compile time.
