---
version: patch
---

- Feat: configure read-header, read, write, and idle timeouts with documented defaults on every App
- Feat: add MaxBodyReader middleware for per-route request body limits without weakening unrelated routes
- Feat: add ServerConfig and SetServerConfig for tuning connection, header, and shutdown timeouts before build
- Fix: Listen now waits for graceful shutdown completion and propagates shutdown failures
- Fix: signal subscriptions are released on Listen return so repeated server lifecycles do not leak goroutines
- Docs: document server timeouts, limits, body reader, and shutdown behavior in runtime.md