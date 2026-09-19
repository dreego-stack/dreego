---
version: patch
---

- Fix: scaffolded `web-minimal` and `web-app` layouts now emit `<!DOCTYPE html>`, so generated pages no longer render in quirks mode
- Feat: `dreego generate` warns when the route `<body>` tag carries HTML attributes (for example `x-data`/`x-init`) that the layout's `<body>` cannot apply, instead of dropping them silently
- Feat: `dreego generate` warns when a route uses Alpine.js directives under the default CSP, which blocks Alpine because it needs `'unsafe-eval'`
- Docs: document per-library copy-paste CSP examples (HTMX, Alpine.js) in the progressive-enhancement guide and state that the default CSP has no `unsafe-eval`
- Fix: the `expected root section` error now points to the `<body>` placement for a `<!DOCTYPE>` or `<html>` at the top of a file
- Fix: head dedupe now also removes duplicated `<meta charset>` and `<meta name="viewport">`, not only `<title>` and meta description
- Fix: the missing-session-store CSRF warning is suppressed for pure read-only apps that register no state-changing route
- Docs: consolidate all first-party documentation into the single root `_docs` tree with module index pages, and update `dreego docs`, the sitemap, and all cross-links
- Feat: scaffolded projects ship `Dockerfile` and `docker-compose.yml` in the `_common` template overlay
