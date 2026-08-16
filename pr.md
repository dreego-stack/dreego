---
version: patch
---

- Feat: dreegotest generation parity test — `_tests/go/parity_test.go` compares CLI `generate` output with the `dreegotest` generation path for the same fixtures; any divergence fails CI with a diff
- Chore: remove the obsolete shell test-header contract — `_tests/go/standard_header_test.go` and `_tests/how-to-test-sh.md` are deleted; Go integration tests are discoverable through normal Go conventions
- Chore: test runner fails when intended Go test packages are missing or tests are skipped unexpectedly
