---
id: queue-interface.1
title: Background Job Queue Interface
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core interface for background jobs. Abstracts Redis/NATS/In-Memory. Job middleware, batching, chaining, delayed dispatch. Plugin implementations (dreego-jobs). Like Laravel Queues.
