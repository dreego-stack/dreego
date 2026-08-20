---
area: code
phase: pre-v0.1
---
# Remove dead private helpers and empty files

## Goal
Remove dead private helpers and empty files:
- `core/middleware_health.go` (empty)
- `core/session_opts.go` (`opt[T]` unused)
- `core/static.go` `staticPattern`
- `core/forms.go` `actionFuncName`
- `core/parser.go` `isSectionTag`
- `core/parser_section_head/script/style` `parse*Section` (if truly unused — verify against `parser.go` `parseNonDivSection`)
- `core/lexer_token.go` `TokenComponentHeader`/`TokenImport` (if never emitted)
- `core/codegen_template.go` `genTemplateNodeTo`
- `core/codegen_layout.go` `genLayoutNode`
- `core/codegen_text_section.go` `compTextWithAttrs`

## Acceptance criteria
- Verify each candidate is truly unused before removal.
- Keep any helper still referenced by tests in a meaningful way; otherwise
  remove the test too.
- `go test ./core/...` is green.
