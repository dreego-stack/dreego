---
id: event-bus.1
title: Pub/Sub Event-Bus (Core-Interface)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
created: 2026-07-26
changed: 2026-07-26
---

Core-Interface für Event-Bus (Pub/Sub). Abstrahiert Redis/NATS/In-Memory. Methoden: Publish, Subscribe, Unsubscribe. Typisiert via Go-Generics. Plugin-Implementierungen. Nutzbar von Notifications, SSE, Real-Time-Features.
