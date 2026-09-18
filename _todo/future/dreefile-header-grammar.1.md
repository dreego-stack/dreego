# Dreefile header grammar slices 3b and 3c

- Area: compiler
- Phase: Dreefile grammar
- Goal: complete the grammar migration and retire the legacy header forms.
- Gap: slices 3b (migrate the 7 real `.dreego` files, the public documentation, and the `_tests/go` inline declarations) and 3c (reject the legacy `Component` / `import` / `from` forms with a migration diagnostic) are landed. What remains is follow-up hardening only: the `_todo/core` items for parser validation, alias scope, and `dreegotest` support of `DREEFILE` components.
- Acceptance: no `.dreego` file or test string uses the legacy forms (done); the legacy forms fail at generate with a migration diagnostic (done); the grammar is documented (done). Remaining follow-up is tracked by the individual `_todo/core/dreefile-*` and `head-*` items.
- Depends on: Dreefile grammar slices 3a (parser and fmt) and the per-finding header hardening items.
