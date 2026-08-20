---
version: none
---

- Refactor: split the runtime into `core/internal/` subpackages (server, session, middleware, context, validate); `core/` re-exports the public API, import path unchanged
