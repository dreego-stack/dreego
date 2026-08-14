---
version: patch
---

- Breaking: replace all package-global runtime state with an explicit App instance (New, Register, Build, Handler, Listen)
- Breaking: Plugin.RegisterRoutes() now receives *App instead of using global registration
- Breaking: generated code emits gen.Register(app) instead of init() with dreego.Register
- Feat: multiple App instances are isolated — two apps can run concurrently with different routes, middleware, and sessions