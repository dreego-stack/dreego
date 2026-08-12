---
version: patch
---

- Feat: add core Queue interface (queue-interface.1) — `Queue` interface (Dispatch/DispatchAfter/DispatchBatch/Worker/Use) + `Job`/`JobHandler`/`JobMiddleware`, like `database/sql`, interface only, plugins implement (Redis/NATS/In-Memory)
