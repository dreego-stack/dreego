# Dead legacy header parsers and formatters kept alive by tests

- Area: compiler
- Phase: Dreefile grammar (slice 3c)
- Goal: remove the legacy header parsing and formatting helpers that no production path calls.
- Gap: `parseFromImport` (`internal/transpiler/lexer/lexer_header.go:90-120`), `parseComponentHeader` (`internal/transpiler/lexer/lexer_header.go:122-154`) and `parseImportLine` (`internal/transpiler/lexer/lexer_header.go:182-198`) have no production callers; `ParseFileHeaderStrict` rejects every legacy form before they are reached. Likewise `formatCompHeader` (`internal/transpiler/fmt.go:160-186`) and `formatImport` (`internal/transpiler/fmt.go:188-198`) are unreachable from `Format`.
- Gap: their only references are tests that pin removed behaviour — `TestParseComponentHeader*` and `TestParseImportLine*` (`internal/transpiler/lexer/lexer_header_test.go:5-103`), `TestFormatCompHeader*` and `TestFormatImport*` (`internal/transpiler/fmt_test.go:8-65`). `parseProps` is still used by the new directive parser (`internal/transpiler/lexer/lexer_header_directives.go:50`) and must stay.
- Acceptance: delete the functions and their tests, or mark them deprecated with a reason if a caller is intended.
- Depends on: nothing.
