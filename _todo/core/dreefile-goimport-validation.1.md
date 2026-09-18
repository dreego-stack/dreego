# `GOIMPORT` and `LAYOUT` accept empty or unquoted values

- Area: compiler
- Phase: Dreefile grammar
- Goal: validate `GOIMPORT` list entries and `LAYOUT` path syntax at generate time.
- Gap: `parseGoImport` (`internal/transpiler/lexer/lexer_header_directives.go:122-141`) accepts an empty list (`GOIMPORT {}` returns zero paths) and arbitrary strings, so an unquoted multi-word entry like `encoding json` passes as one path.
- Gap: `parseLayoutLine` (`internal/transpiler/lexer/lexer_header_directives.go:77-79`) delegates to `parseQuotedValue` (`:69-75`), which returns `""` for unquoted input; `LAYOUT www/layouts/x.dreego` is therefore stored as empty and the layout is silently ignored.
- Acceptance: empty lists and empty or unquoted paths fail at generate with a diagnostic naming the directive.
- Depends on: nothing.
