---
id: golden-tests-core.1
title: Golden Code Tests for Generator Output
status: planned
phase: pre-v1.0
requires:
  - dreegotest.1
created: 2026-07-31
changed: 2026-07-31
---

Add golden-file tests that compare generated `gen/routes.go` and `gen/components.go` against expected output for representative `.dreego` inputs. Catches CodeGen regressions at `go test` time, before integration tests run.
