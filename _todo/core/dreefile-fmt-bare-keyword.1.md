# `dreego fmt` pushes a bare `DREEFILE` line into the body

- Area: cli
- Phase: Dreefile grammar (slice 3a)
- Goal: make `dreego fmt` and `dreego generate` agree on a bare `DREEFILE` line without a value.
- Gap: the parser recognises a bare `DREEFILE` line (`internal/transpiler/lexer/lexer_header.go:23`) and then rejects it as an invalid value, but `Format` only matches `DREEFILE ` with a trailing space (`internal/transpiler/fmt.go:44`); a bare `DREEFILE` therefore still counts as a body line, so `dreego fmt` moves it out of the header where generate errors on the same input.
- Acceptance: both commands treat the line the same way; add a round-trip test in `internal/transpiler/fmt_dreefile_test.go`.
- Depends on: nothing.
