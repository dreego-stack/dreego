# Future grammar todo overclaims that no test string uses legacy forms

- Area: docs
- Phase: Dreefile grammar
- Goal: keep the future grammar todo scoped to what is actually true.
- Gap: `_todo/future/dreefile-header-grammar.1.md:7` states that no test string uses the legacy forms, but the internal transpiler tests still carry them: `Component` in `internal/transpiler/lexer/lexer_header_directives_test.go:159`, `internal/transpiler/fmt_dreefile_test.go:162` and `internal/transpiler/fmt_test.go:140`. The migration holds only for `_tests/go`.
- Acceptance: narrow the claim to the `_tests/go` scope (or the scope that is true after the recheck).
- Depends on: nothing.
