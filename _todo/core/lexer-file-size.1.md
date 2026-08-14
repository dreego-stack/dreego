---
area: transpiler
phase: pre-v0.1
---
# Split the lexer into focused files

## Goal
Bring `core/lexer.go` below the project limit of 300 lines without changing
lexer behavior.

## Acceptance criteria
- Each resulting file owns one clear lexer responsibility and stays at or below
  300 lines.
- Existing lexer and integration tests remain green.
- The refactor does not change accepted syntax, diagnostics, or source locations.
