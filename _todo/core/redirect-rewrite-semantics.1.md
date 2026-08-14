---
area: routing
phase: v0.1-blocker
depends_on: app-runtime.1
---
# Define segment-safe redirects and rewrites

## Goal
Prevent a rule for `/api` or `/api/*` from unexpectedly matching `/apiary`.

## Acceptance criteria
- Exact and wildcard rules have separate documented semantics.
- Wildcards respect path-segment boundaries and canonical path handling.
- Invalid patterns, targets, loops, and redirect status codes fail during configuration.
- Negative HTTP tests cover near-prefix paths and malformed rules.
