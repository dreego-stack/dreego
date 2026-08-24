---
area: architecture
phase: planned-v0.5
---
# Expanded Wails integration

## Goal
Add the first-party Wails target after the v0.2 render foundation. It is
sequenced after SSG in the roadmap but does not depend on an SSG-specific API.
See `_plan/v0.5-wails-target.md`.

## Acceptance criteria
- Components render without assuming an HTTP request.
- Desktop-specific APIs remain explicit and do not leak into target-neutral or
  SSR behavior.
- A reference desktop application verifies the supported workflow.
- No hidden localhost HTTP server or developer-managed npm pipeline is required.
