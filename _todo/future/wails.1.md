---
area: architecture
phase: wails-host
---
# Expanded Wails integration

## Goal
Add the first-party Wails host after the render foundation and the TypeScript
and Lua client processors prove the shared JavaScript output pipeline. See
`_plan/phase-wails.md`.

## Acceptance criteria
- Components render without assuming an HTTP request.
- Desktop-specific APIs remain explicit and do not leak into target-neutral or
  SSR behavior.
- A reference desktop application verifies the supported workflow.
- No hidden localhost HTTP server or developer-managed npm pipeline is required.
