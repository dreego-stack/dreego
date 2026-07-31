---
id: dreego-jobs.1
title: dreego-jobs (Background Jobs, Cron, Queue)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - queue-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Target directory: `plugins/jobs/` in this repository.

Plugin for background jobs. Implements queue-interface.1. Redis/NATS/In-Memory backend. Job middleware, batching, retry, delayed dispatch. Cron scheduler. Dashboard (like Laravel Horizon).
