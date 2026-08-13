---
area: components
phase: v0.1-blocker
depends_on: pre-v0.1-product-decisions.1
---
# Implement the selected component call contract

## Goal
Make props, nested components, and slots behave exactly as selected in the
pre-v0.1 product decisions.

## Acceptance criteria
- Prop binding follows the selected named or positional contract without silent misbinding.
- Missing, unknown, and duplicate props produce source-aware errors where applicable.
- Nested component output is preserved for self-closing and child-bearing calls.
- Slot state is scoped to one component invocation and cannot leak to siblings.
- HTTP-level tests assert rendered output instead of compilation alone.
