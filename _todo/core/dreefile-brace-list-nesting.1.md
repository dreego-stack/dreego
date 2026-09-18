# `parseBraceList` ignores nesting and trailing content

- Area: compiler
- Phase: Dreefile grammar
- Goal: make brace-list parsing reject or handle nesting and trailing text.
- Gap: `parseBraceList` (`internal/transpiler/lexer/lexer_header_directives.go:143-177`) scans for the first `}` and splits on commas, so `COMPONENT "x" IMPORT { A, B } extra` silently drops the trailing token and a nested `{ A, { B } }` is flattened by the comma split.
- Acceptance: nested or extra content is either handled or rejected with a diagnostic.
- Depends on: nothing.
