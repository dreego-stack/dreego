# Dreefile header grammar slices 3b and 3c

- Area: compiler
- Phase: Dreefile grammar
- Goal: complete the grammar migration and retire the legacy header forms.
- Gap: slice 3b (migrate the 7 real `.dreego` files and the 17 `_tests/go` files off the legacy `Component` / `from` / `import` forms) is still open; slice 3c (turn the legacy forms into hard errors with a migration diagnostic) also remains, and touches roughly 35 inline component declarations inside Go test strings.
- Acceptance: no `.dreego` file or test string uses the legacy forms; the legacy forms fail at generate with a migration diagnostic; the grammar is documented.
- Depends on: Dreefile grammar slices 3a (parser and fmt) and the per-finding header hardening items.
