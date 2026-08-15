---
title: component-correctness.1 — Implement the selected component call contract
status: plan-review
---

# Goal
Make props, nested components, and slots behave exactly as selected in the accepted named-prop and lexical-slot contract from `_todo/core/component-correctness.1.md`.

# Acceptance criteria (from todo)
1. Props are named and order-independent without silent positional binding.
2. Missing, unknown, and duplicate props produce source-aware errors where applicable.
3. Nested component output is preserved for self-closing and child-bearing calls.
4. Slot state is scoped to one component invocation and cannot leak to siblings.
5. HTTP-level tests assert rendered output instead of compilation alone.

# Sub-increments (each: tests → code → review → docs → commit)

## A. Named props and source-aware prop errors
- Tests in `_tests/go/component_props_test.go`:
  - order-independent named props render correctly;
  - missing prop → `dreego generate` error with source location;
  - unknown prop → `dreego generate` error with source location;
  - duplicate prop in same call → `dreego generate` error with source location.
- Code in `core/` to enforce named binding and report errors.
- Review.
- Docs: `_docs/components.md` prop rules + `_docs/testing.md` table update.
- Commit.

## B. Nested component output preservation
- Tests in `_tests/go/component_nested_test.go`:
  - self-closing nested component call renders inner output;
  - child-bearing nested component call renders inner output inside outer slot.
- Code in `core/` to preserve nested component output in both forms.
- Review.
- Docs: `_docs/components.md` nested usage.
- Commit.

## C. Lexical slot scoping
- Tests in `_tests/go/component_slots_test.go`:
  - slot content is bound to the intended component invocation;
  - sibling component invocations do not see each other's slots;
  - named slots do not leak into default slot of another invocation.
- Code in `core/` to keep slot state per component invocation.
- Review.
- Docs: `_docs/components.md` slot scoping rules.
- Commit.

## D. HTTP-level rendered output assertions
- Tests in `_tests/go/component_http_test.go`:
  - `dreegotest.NewApp(app).Get` asserts full HTML body contains expected component output;
  - covers self-closing, child-bearing, nested, and slot cases.
- Code: ensure generated components render over HTTP (likely already works, tests prove it).
- Review.
- Docs: `_docs/testing.md` HTTP assertion examples.
- Commit.

# Constraints
- Max 300 lines per file; split if needed.
- All code/comments/tests in English.
- No external dependencies in `core/`.
- Every bug-like behavior gets a regression-style test.
- Use `dreegotest` helpers as documented in `_docs/testing.md`.
- Each increment must leave `make test` green.
- After final increment: delete `_todo/core/component-correctness.1.md` in the last PR step (do not delete now).

# Next
Human reviews this plan; then increment A starts with test writing.
