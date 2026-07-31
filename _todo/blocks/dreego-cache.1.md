---
id: dreego-cache.1
title: dreego-cache (Caching: Memory, Redis)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - cache-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Target directory: `plugins/cache/` in this repository.

Plugin for caching. Implements cache-interface.1. In-Memory (dev), Redis (prod). Usable by session store, template cache, DB query cache. TTL, Tags, Cache warming.
