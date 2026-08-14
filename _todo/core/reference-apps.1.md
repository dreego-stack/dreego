---
area: testing
phase: v0.1-blocker
depends_on:
  - app-runtime.1
  - component-correctness.1
  - routing-correctness.1
---
# End-to-end reference applications

## Goal
Create a small set of documented applications that also verify the complete public CLI-to-HTTP workflow.

## Acceptance criteria
- Fixtures cover minimal usage, forms and sessions, components, and a plugin.
- CI generates, builds, starts, and verifies each application through HTTP.
- Fixtures teach public APIs rather than internal test shortcuts.
