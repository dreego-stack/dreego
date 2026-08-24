---
area: architecture
phase: planned-v0.4
---
# Static site generation

## Goal
Generate static HTML, JavaScript, and CSS for hosts such as Cloudflare Pages and GitHub Pages.

Implementation begins only after the v0.2 render foundation is complete. See
`_plan/v0.4-ssg-target.md`.

## Acceptance criteria
- Static rendering reuses proven `.dreego` semantics without weakening SSR.
- Dynamic routes and content collections have explicit generation rules.
- Deployment is documented for at least Cloudflare Pages and GitHub Pages.
- Server-dependent DreeJS capabilities require an explicit external endpoint.
