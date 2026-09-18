# Conflicting Dreefile header directives are accepted silently

- Area: compiler
- Phase: Dreefile grammar
- Goal: reject contradictory or duplicate kind directives in a file header.
- Gap: `ParseFileHeaderStrict` applies `DREEFILE` and `LAYOUT` top down without conflict detection; two `LAYOUT` lines overwrite each other silently, and a `DREEFILE` kind plus another `DREEFILE` kind is not rejected. The legacy `Component` line is no longer part of this: it now returns a hard legacy-header error before a conflict can form.
- Acceptance: conflicting directives fail at generate with `file:line:col`.
- Depends on: nothing.
