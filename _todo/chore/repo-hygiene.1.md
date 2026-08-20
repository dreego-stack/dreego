---
area: chores
phase: pre-v0.1
---
# Repository hygiene

## Goal
- Remove the dead code path in `_tests/test.sh` (legacy `_tests/core` loop,
  ~line 114).
- Fix the `AGENTS.md` file-structure sketch (`_tests/core/<Category>/` →
  `_tests/go/` + `_tests/fixtures/`).
- Add `tmp/` to `.gitignore` if not present.
- Delete the empty `bin/` directory.

## Acceptance criteria
- `make test` is green.
- `AGENTS.md` matches reality.
