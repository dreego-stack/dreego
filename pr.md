---
version: patch
---

- Feat: add typed in-memory event bus (event-bus.1) — generic `EventBus[T]` interface (Publish/Subscribe/Unsubscribe) + `NewInMemoryBus[T]()`, concurrency-safe, panic recovery, ctx-cancellation
- Chore: CI standard-header check — `_tests/go/standard_header_test.go` validates every `test.sh` under `_tests/core/` carries the standard header (`# Using standard: _tests/how-to-test-sh.md` + `# What:`), enforced via `make test` in CI
- Chore: AGENTS.md test rules — Feature Workflow + Bug workflow updated to `_tests/go/*_test.go` + `dreegotest`, `test.sh` reference dropped
- Chore: `dreegotest` parallel-safe — codegen runs via cached CLI subprocess instead of global `os.Chdir`, enabling `t.Parallel()` across `_tests/go`; shutdown tests poll port readiness instead of fixed sleeps; `test.sh` prints test output on failure; cross-compile sets `CGO_ENABLED=0`
- Chore: remove `dreego run -t <seconds>` timer flag — the last shell test (`_tests/core/CLI/run-timer`) and the flaky `bug_run_timer_sigterm` integration test are removed; graceful shutdown (B20) stays covered by `TestDeploymentGracefulShutdown`; `standard_header_test` tolerates an empty/missing `_tests/core`
