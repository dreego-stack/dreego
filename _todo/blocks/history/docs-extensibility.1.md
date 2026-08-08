---
id: docs-extensibility.1
title: Extensible dreego docs Command
status: 47
phase: v0.0.25
requires:
  - cli.1
  - plugin-interface.1
created: 2026-07-31
changed: 2026-08-08
---

Adapt `dreego docs` so it can consume documentation from plugins. Official plugins now live under `plugins/<name>/` in this repository, so local `plugins/<name>/_docs/` directories are the primary source. Open design question: how to discover plugin docs elegantly (embed? local path? also support external repos for community plugins?). Must not depend on any specific plugin at compile time.

Done in v0.0.25: `dreego docs` resolves plugin docs from local `plugins/<name>/_docs/`.
