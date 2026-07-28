---
id: cache-interface.1
title: Caching Interface (Memory, Redis)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core interface for caching. Abstracts In-Memory/Redis. Methods: Get, Set, Delete, Has, Remember. TTL-based. Plugin implementations (dreego-cache). Usable by session, templates, DB queries.
