# Dreefile header grammar: remaining hardening and layout chaining

- Area: compiler
- Phase: Dreefile grammar
- Goal: finish the Dreefile grammar capability after the parser, migration and legacy rejection landed.
- Gap: slices 3a (parser and fmt), 3b (migration of the real `.dreego` files, public documentation and the `_tests/go` inline declarations) and 3c (legacy `Component` / `import` / `from` forms are hard generate errors) are LANDED and merged into `stage/dreefile`. The migration claim is scoped to `_tests/go`; the internal transpiler tests intentionally keep legacy header strings to cover `dreego fmt` preservation and the generate rejection diagnostic. What remains is per-item hardening tracked by the individual `_todo/core/dreefile-*` and `head-*` items, plus the layout-chaining slice (a layout declares `LAYOUT` and is selected by another layout, with cycle and missing-target detection) described in `_plan/phase-dreefile.md:149-153`.
- Acceptance: the hardening items are closed and a layout can be chained through an explicit `LAYOUT` reference with cycles and missing targets rejected at generate.
- Depends on: Dreefile grammar slices 3a, 3b and 3c (landed); the per-finding `_todo/core/dreefile-*` and `head-*` hardening items.
