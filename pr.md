---
version: patch
---

- Fix: generation aborts with affected path and wrapped cause on filesystem walk, directory read, and source read failures
- Fix: form parsing and binding failures stay distinguishable from valid empty values via c.FormError() and c.Errors("_form")
- Fix: session read, write, delete, destroy, and CSRF persistence failures reach an application error path via c.SessionError() and fail loudly in CSRF middleware
- Fix: generated render calls no longer discard returned errors (self-closing components panic, form re-render paths propagate)
- Fix: production responses use generic 500 messages; filesystem paths, database details, and Go type errors are no longer disclosed to clients
- Fix: Compress middleware no longer commits a 200 status when a handler panics with gzip enabled
