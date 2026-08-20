---
area: runtime
phase: pre-v0.1
---
# Detect transitive redirect/rewrite cycles at Build() time

## Goal
Detect transitive redirect/rewrite cycles at `Build()` time. A cycle such as
`/a → /b` plus `/b → /a` must fail registration.

## Acceptance criteria
- A 2-hop redirect/rewrite cycle fails at `Build()`.
- A 3-hop redirect/rewrite cycle fails at `Build()`.
- Tests cover both the 2-hop and 3-hop cases.
