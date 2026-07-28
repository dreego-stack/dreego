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

Plugin for distributed deployments. Node discovery, shared session store, PubSub sync, distributed cache. Valkey/Redis as backend. Load balancer-compatible — every node sees same state. No Kubernetes required.
