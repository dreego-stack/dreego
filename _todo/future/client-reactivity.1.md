---
area: architecture
phase: planned-v0.6-v0.7
---
# DreeJS client behavior

## Goal
Implement optional modular client behavior after SSR, target-neutral rendering,
SSG, and Wails establish their contracts. See `_plan/v0.6-dreejs-foundation.md`
and `_plan/v0.7-dreejs-data-live.md`.

## Acceptance criteria
- Static components emit no runtime.
- Local presentation state, serialization, lifecycle, event binding, security,
  accessibility, and minimal module generation are proven first.
- Fetch, poll, stream, and live behavior are implemented independently in that
  order and preserve a Go-first workflow without project-owned npm tooling.
