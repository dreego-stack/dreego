---
version: patch
---

- Breaking: replace all package-global runtime state with an explicit App instance (New, Register, Build, Handler, Listen)
- Breaking: remove the central Plugin interface in favor of package-owned Register(app, typedOptions) error functions
- Breaking: generated code emits gen.Register(app) instead of init() with dreego.Register
- Feat: multiple App instances are isolated — two apps can run concurrently with different routes, middleware, and sessions
- Bug: reject configuration after Build and reject duplicate or reserved routes
- Bug: generate compilable registration code for projects without routes
- Breaking: make dreegotest request clients own an explicit App handler
- Docs: migrate public examples from removed package-global APIs to App methods
