# `dreego fmt` accepts headers that `dreego generate` rejects

- Area: cli
- Phase: Dreefile grammar (slice 3c)
- Goal: align `dreego fmt` with `dreego generate` on files that still carry a legacy header.
- Gap: `Format` copies legacy `Component` / `import` / `from` lines verbatim (`internal/transpiler/fmt.go:37-43`, `internal/transpiler/fmt.go:97-102`), so `dreego fmt --check` reports such a file as already formatted and exits 0 (`cmd/dreego/fmt.go:52-77`) even though `ParseFileHeaderStrict` rejects it at generate (`internal/transpiler/lexer/lexer_header.go:56-69`). The ADR documents the lenient behaviour as intentional (`_docs/decisions/explicit-dreefile-header-grammar.md:59-61`), so the two commands disagree about the same file by design.
- Acceptance: decide and implement one policy — either `fmt` errors like `generate`, or the ADR explicitly states why `fmt` stays lenient — and pin it with a test.
- Depends on: nothing.
