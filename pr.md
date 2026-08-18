---
version: patch
---

- Docs: add compatibility policy defining the breaking-change policy and the v0.1 stability promise
- Docs: document that the plugin contract stays explicitly provisional until v1 and is excluded from the v0.1 stability promise
- Fix: remove the dead exported FieldError type (never used by generated code or applications)
- Docs: correct stale MiddlewareProvider reference in middleware.md to the actual app.Use() plugin model
