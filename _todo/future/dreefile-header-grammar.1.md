# Dreefile header grammar: remaining hardening and layout chaining

- Area: compiler
- Phase: Dreefile grammar
- Goal: finish the Dreefile grammar capability after the parser, migration and legacy rejection landed.
- Gap: slices 3a (parser and fmt), 3b (migration of the real `.dreego` files, public documentation and the `_tests/go` inline declarations) and 3c (legacy `Component` / `import` / `from` forms are hard generate errors) are LANDED and merged into `stage/dreefile`. The migration claim is scoped to `_tests/go`; the internal transpiler tests intentionally keep legacy header strings to cover `dreego fmt` preservation and the generate rejection diagnostic. Explicit `LAYOUT` resolution with cycle and missing-target rejection is LANDED; a chain longer than one layout is now a hard generate error instead of silently dropping inner layouts. What remains is per-item hardening tracked by the individual `_todo/core/dreefile-*` and `head-*` items, plus multi-level (fragment) layout rendering, which is still open and needs its own design because a layout is a full document (`<html>` must not nest).
- Acceptance: the hardening items are closed and multi-level layout rendering is designed and implemented.
- Depends on: Dreefile grammar slices 3a, 3b and 3c (landed); the fragment-rendering design; the per-finding `_todo/core/dreefile-*` and `head-*` hardening items.
