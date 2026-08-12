---
version: patch
---

- Feat: add typed in-memory event bus (event-bus.1) — generic `EventBus[T]` interface (Publish/Subscribe/Unsubscribe) + `NewInMemoryBus[T]()`, concurrency-safe, panic recovery, ctx-cancellation
- Chore: CI standard-header check — `_tests/go/standard_header_test.go` validates every `test.sh` under `_tests/core/` carries the standard header (`# Using standard: _tests/how-to-test-sh.md` + `# What:`), enforced via `make test` in CI
- Chore: AGENTS.md test rules — Feature Workflow + Bug workflow updated to `_tests/go/*_test.go` + `dreegotest`, `test.sh` reference dropped
