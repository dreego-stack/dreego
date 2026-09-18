# Use the `ir.HeadPlaceholder` constant everywhere

- Area: compiler
- Phase: head dedupe cleanup
- Goal: use the single `ir.HeadPlaceholder` constant instead of hardcoded placeholder literals.
- Gap: `ir.HeadPlaceholder` (`internal/transpiler/ir/head_prefix.go:5`) is defined but unused; `internal/transpiler/codegen_layout.go:113,118,134,143` and `SplitHeadPlaceholder` in `internal/transpiler/html/output/templ.go:129-134` hardcode `{#head}`.
- Acceptance: the constant is the single source for the placeholder; no hardcoded `{#head}` remains in non-test code.
- Depends on: nothing.
