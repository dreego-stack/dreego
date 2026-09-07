---
area: architecture
phase: dreejs-and-live-updates
---
# DreeJS client behavior

## Goal
Implement optional modular client behavior after SSR, target-neutral rendering,
the client language processors, and Wails establish their contracts. See
`_plan/phase-dreejs.md` and `_plan/phase-live.md`.

## Acceptance criteria
- Server-only components emit no runtime.
- Local presentation state, serialization, lifecycle, event binding, security,
  accessibility, and minimal module generation are proven first.
- Private fetch islands, poll, stream, and live behavior are implemented independently in that
  order and preserve a Go-first workflow without project-owned npm tooling.
