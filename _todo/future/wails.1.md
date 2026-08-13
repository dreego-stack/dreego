---
area: architecture
phase: v1.x
---
# Expanded Wails integration

## Goal
Improve the existing Wails path after the SSR core reaches v1 stability.

## Acceptance criteria
- Components render without assuming an HTTP request.
- Desktop-specific APIs remain explicit and do not leak into SSR core behavior.
- A reference desktop application verifies the supported workflow.
