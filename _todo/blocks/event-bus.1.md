---
id: event-bus.1
title: Pub/Sub Event Bus (Core Interface)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core interface for event bus (Pub/Sub). Abstracts Redis/NATS/In-Memory. Methods: Publish, Subscribe, Unsubscribe. Typed via Go generics. Plugin implementations. Usable by notifications, SSE, real-time features.
