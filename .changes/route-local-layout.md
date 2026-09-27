---
version: patch
---

- Bug: a route-local `www/routes/<route>/layouts/default.dreego` no longer collides with the root `www/layouts/default.dreego`; layout renderers are scope-qualified (`DefaultRegistrierung`) so several scopes coexist in the `www/layouts` package
- Bug: `www/routes/<route>/layouts/` is no longer scanned as a route directory, so route-local layouts no longer emit a stray route (for example `/registrierung/layouts`)
