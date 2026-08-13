---
area: transpiler
phase: v0.1-blocker
depends_on: pre-v0.1-product-decisions.1
---
# Align frontmatter implementation and documentation

## Goal
Implement or remove public frontmatter according to the pre-v0.1 product decision.

## Acceptance criteria
- Public documentation and the generator expose one consistent contract.
- If retained, generation parses frontmatter before the template and reports source-aware errors.
- If removed, public examples and unused exported APIs are deleted without compatibility shims.
- Black-box tests prove the selected behavior in a real generated route.
