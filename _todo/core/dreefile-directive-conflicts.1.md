# Conflicting Dreefile header directives are accepted silently

- Area: compiler
- Phase: Dreefile grammar
- Goal: reject contradictory or duplicate kind directives in a file header.
- Gap: `ParseFileHeaderStrict` (`internal/transpiler/lexer/lexer_header.go:14-93`) applies `DREEFILE` (`:23-29`), `LAYOUT` (`:32-36`) and the legacy `Component` line (`:56-64`) top down without conflict detection; `DREEFILE layout` plus a later legacy `Component X (..)` silently yields `KindComponent` with `Layout` still set, and two `LAYOUT` lines overwrite each other.
- Acceptance: conflicting directives fail at generate with `file:line:col`.
- Depends on: nothing.
