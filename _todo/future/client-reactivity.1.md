---
area: architecture
phase: between-v0.1-and-v1
---
# Client islands and hydration research

## Goal
Explore a small client-interactivity model without committing an internal runtime to the stable core prematurely.

## Acceptance criteria
- Experiments live outside the stable core first.
- State serialization, lifecycle, event binding, security, and bundling are evaluated.
- Promotion requires evidence from real applications and preserves a simple Go-first workflow.
