---
id: dreego-cluster.1
title: dreego-cluster (Multi-Node, Distributed State)
status: planned
phase: v0.x.0
requires:
  - plugin-interface.1
  - session.1
  - cache-interface.1
  - event-bus.1
created: 2026-07-26
changed: 2026-07-26
---

Plugin für verteilte Deployments. Node-Discovery, Shared-Session-Store, PubSub-Sync, Distributed-Cache. Valkey/Redis als Backend. Loadbalancer-kompatibel — jeder Node sieht gleichen State. Kein Kubernetes nötig.
