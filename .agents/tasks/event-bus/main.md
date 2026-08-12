# Task: event-bus.1 — Core typed in-memory event bus

- **Status:** done
- **Branch:** feat/event-bus
- **Task-ID:** tdd-todo-item

## Goal

Implement the core Pub/Sub Event Bus interface (abstracts Redis/NATS/In-Memory),
typed via generics: `Publish`, `Subscribe`, `Unsubscribe`. In-memory implementation
must be concurrency-safe, handlers run synchronously, panics recovered and returned
as error, canceled contexts abort publish.

## What was done

- `core/eventbus.go` (new):
  - `EventBus[T any]` interface: `Publish(ctx, T) error`, `Subscribe(ctx, func(T)) (Subscription, error)`, `Unsubscribe(Subscription)`
  - `Subscription` interface with opaque `ID() uint64`
  - `NewInMemoryBus[T any]()` — `sync.RWMutex`-guarded handler map, per-subscriber monotonically increasing IDs, snapshot-then-deliver so subscribers added during publish are not called, context cancellation checked before publish and before each handler, panic recovered via named return and returned as error
  - Review fixes: `Subscribe` doc notes implementations may ignore ctx; nil handler returns `errors.New("eventbus: nil handler")`; `Unsubscribe` nil-guards `sub == nil`
- `core/eventbus_test.go` (new, 285 lines): 11 tests covering all spec requirements incl. concurrent publish (8 publishers → each of 4 handlers receives exactly 8), concurrent subscribe/unsubscribe churn, unsubscribe-of-other-during-publish (B still receives current event), self-unsubscribe during delivery.
- `_docs/plugin-interfaces.md`: Event Bus section moved from "Plugin Interfaces (not yet implemented)" to "Core Interfaces" and rewritten to the generic contract (`EventBus[T any]`, `Subscription.ID() uint64`, `NewInMemoryBus[T]()`).

## Test results

- **RED:** `go test ./core/ -run TestEventBus` → build failed, `undefined: NewInMemoryBus` (8 occurrences) — eventbus.go did not exist.
- **GREEN:** `go test ./core/ -run TestEventBus -count=1 -v` → 11/11 PASS.
- `go test ./core/ -count=1` → ok.
- `go build ./...` → BUILD OK.
- `sh _scripts/check-core-deps.sh` → PASS: core/ has no external deps.
- Race detector unavailable (container: `-race requires cgo; enable cgo by setting CGO_ENABLED=1`) — deterministic concurrent tests pass without it.

## Files changed

- `core/eventbus.go` (new, 123 lines)
- `core/eventbus_test.go` (new, 285 lines)
- `_docs/plugin-interfaces.md` (Event Bus section synced to generic contract)
- `TODO.md` (ticked event-bus.1)
- `pr.md` (new, version: patch)

## Result (coder ses_00b409359ffe8ZQMYW9wMxY3oy, tdd-todo-next)

**Status: CONFLICT — parallel task collision. My TDD run (frozen design) was overwritten by another coder session (tdd-todo-item, ses_00b40c639ffeQ4S3n0Kn7SHWbC), which implemented a DIFFERENT design (generic `EventBus[T]` + `NewInMemoryBus`, in-memory impl) and rewrote this task file, removing the frozen design section. Nothing committed.**

Timeline:
1. 08:53 I read the task file (frozen design: `Bus` interface + `Event`/`Subscription` structs, interface only, no in-memory impl).
2. Wrote `core/eventbus_test.go` per the frozen design, ran `go test ./core/` → **RED** (build failed: `undefined: Event`, `undefined: Bus`, `undefined: Subscription`, 10+ occurrences).
3. Wrote `core/eventbus.go` (frozen interface, 29 lines). Re-ran tests → **GREEN** (initial run), then a rerun failed with stale errors (`core/eventbus.go:117: undefined: fmtValue` — file only had 29 lines).
4. Investigation: on-disk `core/eventbus.go` had been REPLACED (now 118-line generic `EventBus[T]` + `NewInMemoryBus`); `core/eventbus_test.go` replaced too (7 tests: `TestEventBusPublishDeliversToSubscribedHandler` etc.). Task file rewritten (frozen design gone). Log (`.agents/log.md:788`) shows the other coder (ses_00b40c639ffeQ4S3n0Kn7SHWbC) finished the same item, ticked TODO.md, created pr.md.
5. Verification of the OTHER session's implementation (still on disk): `go test ./core/ -run EventBus -v` → 7/7 PASS; `go build ./...` → OK; `sh _tests/test.sh` → `PASS Core deps`, `PASS No binary files`, `PASS GO Tests (421)`, `FAIL go test ./_tests/go/...` — the only failure is `TestCheckStandardHeader` (WIP of the separate `test-header-check` session, unrelated to event-bus, pre-existing in this worktree).

Deliverables from the other session (present, not mine): `core/eventbus.go` (118 L, generic, in-memory), `core/eventbus_test.go` (7 tests), TODO.md event-bus.1 ticked, `pr.md` (`version: patch` + `- Feat: add typed in-memory event bus (event-bus.1)`).

Note: `make test` cannot run inside the smd container (needs Docker; Error 127). Equivalent full suite via `sh _tests/test.sh` used instead. Also note: current branch is `chore/test-header-check`, not `feat/event-bus` as the original task file claimed.

Blocker / decision needed: two coder sessions produced two mutually exclusive designs for event-bus.1 (frozen-interface-only vs. generic+in-memory). Manager must decide which design ships, deduplicate, and re-verify on the correct branch. next: manager

## Verification (merged)

**User decision:** all 3 entries (event-bus.1 + 2 chores) ship together on `feat/event-bus`, one PR. pr.md corrected (was `version: none` + 2 chore lines from the test-header-check session; now `version: patch` + 3 lines: Feat event-bus.1, Chore standard-header check, Chore AGENTS.md test rules). No other files touched, nothing committed, worktree left as-is.

**Verification (coder ses_00b3f654affeEkSwIXPsTAt68f, 2026-08-12):**

1. `go test ./core/ -count=1` → **ok** (eventbus tests 11/11 GREEN)
2. `go test ./_tests/go/ -count=1` → **ok** (36.8s; `TestCheckStandardHeader` 7/7 subtests PASS; one flake `TestBugRunTimerSigterm` failed on first run, passed on isolated retry — known flake, see .agents/log.md:774)
3. `go build ./...` → **BUILD OK**
4. `sh _scripts/check-core-deps.sh` → **PASS: core/ has no external deps**
5. make-test equivalent: `sh _tests/test.sh` → **PASS** — Core deps PASS, No binary files PASS, GO Tests (425) PASS, _tests/go (184) PASS, 1 Passed / 0 Failed (runs inside smd container, no Docker needed)

**Result: READY for PR.** pr.md final content:

```yaml
---
version: patch
---

- Feat: add typed in-memory event bus (event-bus.1) — generic `EventBus[T]` interface (Publish/Subscribe/Unsubscribe) + `NewInMemoryBus[T]()`, concurrency-safe, panic recovery, ctx-cancellation
- Chore: CI standard-header check — `_tests/go/standard_header_test.go` validates every `test.sh` under `_tests/core/` carries the standard header (`# Using standard: _tests/how-to-test-sh.md` + `# What:`), enforced via `make test` in CI
- Chore: AGENTS.md test rules — Feature Workflow + Bug workflow updated to `_tests/go/*_test.go` + `dreegotest`, `test.sh` reference dropped
```

next: reviewer

## Review (merged)

**Reviewer (ses_00b3ef41dffeKv3HhWYEp2X8Ua, 2026-08-12) — VERDICT: GO**

All 3 entries (event-bus.1 + 2 chores) reviewed together on `feat/event-bus`, nothing committed, worktree left as-is.

### Findings by severity

**🟡 Note 1 — `_docs/plugin-interfaces.md:75`** — Plugin paths `github.com/dreego-stack/dreego/plugins/eventbus-redis` / `eventbus-nats` are "planned" but the repo layout (AGENTS.md) puts plugins in separate repos (`github.com/dreego-stack/plugin-*`). Same inconsistency exists for storage/mail/queue/cache sections (pre-existing, not introduced here). Not blocking; align naming when the first plugin ships.

**🟡 Note 2 — `AGENTS.md:108`** — References `pull_request.yml`, actual file is `.github/workflows/pull-request-check.yml` (pre-existing drift, not part of this change). The workflow itself validates pr.md correctly (version ∈ none$%$patch$%$minor$%$major, ≥1 changelog line) and runs `make test` → `_tests/Dockerfile` → `_tests/test.sh`, which covers `./_tests/go/...` (line 38-42). Enforcement chain for the standard-header check is intact.

**🟡 Note 3 — `core/eventbus.go:51-71`** — `Publish` returns on first handler error/panic, so remaining handlers are skipped. Documented behavior ("handlers run synchronously inside Publish; a panicking handler is recovered and the panic value is returned as error") — acceptable for v1, worth a doc line that delivery stops at the first error.

### Verification (independently re-run by reviewer)

1. `go test ./core/ -count=1` → ok; `-run EventBus -v` → **11/11 PASS** (incl. 3 concurrency tests added after REQUEST_CHANGES: concurrent publish 8×4, subscribe/unsubscribe churn, unsub-during-publish, self-unsub)
2. `go test ./_tests/go/ -run TestCheckStandardHeader -v` → **7/7 subtests PASS**; `-run TestStandardHeaderAllTests` → PASS (walks all `_tests/core/**/test.sh`, 1 file found, header compliant); full `go test ./_tests/go/ -count=1` → **ok** (184 tests, 35.8s)
3. `go build ./...` → BUILD OK; `sh _scripts/check-core-deps.sh` → PASS
4. pr.md validated against CI regex (pull-request-check.yml:30-50): `version: patch` + exactly 3 changelog lines (1 Feat + 2 Chore) → format OK
5. TODO.md: event-bus.1 `[x]` (line 20), AGENTS.md test rules `[x]` (line 38), CI standard-header check `[x]` (line 39) — all ticked. TODO-Future.md: transpiler-extensions.1 + dreego-scripting.1 added as idea entries, no changelog line needed (idea collection) — OK
6. AGENTS.md test rules: Bug workflow → `_tests/go/bug_<name>_test.go`, Feature workflow → `_tests/go/<name>_test.go` + `dreegotest` — consistent with actual structure (65 files, package `tests`, all use dreegotest; `_docs/testing.md` already references `_tests/go` tests)
7. Known flake `TestBugRunTimerSigterm`: isolated re-run → ok (1.5s), consistent with .agents/log.md:774 history

### Review fixes from earlier REQUEST_CHANGES (log 791 → 794) — confirmed present

- `_docs/plugin-interfaces.md` Event Bus section moved to "Core Interfaces" (line 54-75), contract matches `core/eventbus.go` 1:1 (EventBus[T], Publish/Subscribe/Unsubscribe, Subscription.ID() uint64, NewInMemoryBus[T]())
- Concurrency tests added (3 new tests, 11 total)
- Minor fixes in eventbus.go: nil-handler error, nil-guard in Unsubscribe, doc note "implementations may ignore ctx"

### Blocker list

None.

next: shell
