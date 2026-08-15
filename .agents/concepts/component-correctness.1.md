---
type: concept
status: draft
related: .agents/concepts/component-call-contract.md
---
# Component Correctness Implementation Plan

This plan implements the contract defined in
[component-call-contract.md](component-call-contract.md).
Each increment adds HTTP-level assertions on rendered output.

## Increments

### A: Named prop validation

- Tests: `_tests/go/component_call_named_props_test.go`
  - order-independent named props
  - missing required prop error
  - unknown prop error
  - duplicate prop error
  - source line accuracy
- Code: `core/` component call parsing and declaration matching
- Review gates: Flash writes, Pro verifies `make test`, ≤300 lines per file,
  English only, no external dependencies in `core/`
- Docs: update `_docs/components.md` with error examples
- Commit: `component-correctness A: named prop validation`

### B: Expression prop type checking

- Tests: `_tests/go/component_call_expr_props_test.go`
  - string literal passed to string prop
  - int expression passed to int prop
  - wrong-type expression prop error
  - source line accuracy
- Code: `core/` expression type extraction and prop type comparison
- Review gates: Flash writes, Pro verifies `make test`, ≤300 lines per file,
  English only, no external dependencies in `core/`
- Docs: update `_docs/components.md` expression prop rules
- Commit: `component-correctness B: expression prop type checking`

### C: Self-closing and slot fallback

- Tests: `_tests/go/component_call_slots_test.go`
  - self-closing call with no body
  - self-closing call with children error
  - default slot fallback (empty children)
  - source line accuracy
- Code: `core/` call body handling and slot fallback
- Review gates: Flash writes, Pro verifies `make test`, ≤300 lines per file,
  English only, no external dependencies in `core/`
- Docs: update `_docs/components.md` slot fallback note
- Commit: `component-correctness C: self-closing and slot fallback`

### D: Named slots and sibling isolation

- Tests: `_tests/go/component_call_named_slots_test.go`
  - named slot render
  - unknown named slot error
  - nested component inside slot
  - sibling component calls do not leak slot content
  - source line accuracy
- Code: `core/` named slot extraction and per-call slot scope
- Review gates: Flash writes, Pro verifies `make test`, ≤300 lines per file,
  English only, no external dependencies in `core/`
- Docs: update `_docs/components.md` named slot examples
- Commit: `component-correctness D: named slots and sibling isolation`

## Final PR wrap-up

1. Create `pr.md` from `pr.md.example` with `version: patch` and changelog lines
   for increments A–D.
2. Update `CHANGELOG.md` with the same entries under the next patch version.
3. Update `.agents/log.md` with a summary of the work.
4. Delete `_todo/core/component-correctness.1.md`.

Run `make test` and confirm all new and existing tests pass before opening the
PR.
