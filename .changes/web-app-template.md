---
version: patch
---

- Feat: add the `web-app` project template — a full SSR application starter with an accessible app shell, `Nav` and `Card` components, a nested `dashboard` route, and a server-rendered note form
- Feat: register `web-app` in the template embed list and cover it with registry, install, and end-to-end CLI tests; `web-minimal` remains the default template
- Fix: give each `web-app` route its own `<title>` (WCAG 2.4.2) and drop the layout's duplicate title so every page renders exactly one
- Fix: synchronise the `web-app` notes list with a `sync.Mutex` so `go test -race` stays clean under concurrent submissions
- Fix: remove the inert `csrf_token` hidden field from the `web-app` form instead of implying CSRF protection that the minimal starter does not configure
- Test: add serve-based `web-app` coverage for titles, `aria-current` navigation, the dashboard table, the form round-trip, and race-safe concurrent note posts
