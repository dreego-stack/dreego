---
area: transpiler
phase: v0.1-blocker
---
# Preserve valid Go and HTML source semantics

## Goal
Prevent the lexer, parser, and filter syntax from changing valid source text.

## Acceptance criteria
- `<` comparisons and delimiter-like strings inside Go blocks do not close or corrupt the block.
- `>` inside quoted HTML attributes does not terminate the tag.
- Script and style content use explicit delimiter rules.
- Go operators such as `||` are not parsed as template filters.
- Unknown filters fail with a source position instead of being ignored.
- Regression tests assert generated and rendered behavior for every case.
