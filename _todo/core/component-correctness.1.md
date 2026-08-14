---
area: components
phase: v0.1-blocker
---
# Implement the selected component call contract

## Goal
Make props, nested components, and slots behave exactly as selected in the
accepted named-prop and lexical-slot contract.

## Acceptance criteria
- Props are named and order-independent without silent positional binding.
- Missing, unknown, and duplicate props produce source-aware errors where applicable.
- Nested component output is preserved for self-closing and child-bearing calls.
- Slot state is scoped to one component invocation and cannot leak to siblings.
- HTTP-level tests assert rendered output instead of compilation alone.
