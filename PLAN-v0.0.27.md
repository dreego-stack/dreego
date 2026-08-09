# Plan v0.0.27 — Chore Batch: Platform Migration + Repo Restructure

**Type:** chore-only (no features, no breaking behavior changes for end users)
**Status:** in progress
**Scope:** platform move, single-module restructure, plugin extraction, PR-based CI, docs, test migration, cleanup

## Decisions (locked by user)

1. **Platform:** GitHub, org `dreego-stack`. Main repo: `github.com/dreego-stack/dreego`
2. **Plugins:** Each official plugin gets its own repo under `github.com/dreego-stack/<plugin-name>`
3. **v0.0.27 = chores only** — no new features, no behavior changes
4. **Releases are PR-driven:** every change lands via a PR with `pr.md` (version bump + changelog lines); the tag is created by CI after merge. No local tags, no `_scripts/release.sh`.
5. **No secrets in workflows:** Dependabot handles plugin dependency updates (auto PRs), no PAT/App needed.

---

## Phase 1 — Platform Migration (codeberg → github.com/dreego-stack) ✅ DONE

- Global module path rename `codeberg.org/dreego/dreego` → `github.com/dreego-stack/dreego` (245 files)
- `origin` → `https://github.com/dreego-stack/dreego.git`, codeberg removed
- All 36 tags pushed to GitHub
- AGENTS.md + LICENSE updated
- Commit: `d0a378c`

## Phase 2 — Single Root Module ✅ DONE

- Root `go.mod` (`module github.com/dreego-stack/dreego`), deleted `core/go.mod`, `cmd/dreego/go.mod`, `dreegotest/go.mod`, `go.work`
- Import paths stay identical; one tag per release
- 182 shell tests: `require github.com/dreego-stack/dreego` + `replace => $realrepo`
- `_tests/how-to-test-sh.md` template updated
- `demo/` stays standalone module (`replace => ..`)
- `release.sh` → single tag (later removed entirely, see Phase 4)
- `new.go` scaffold: require/replace root module, `findLocalRepo()`
- `_scripts/check-core-deps.sh` enforces core has no external deps (replaces module boundary)
- Fresh clone builds without go.work
- Commits: `ac5adec`, `683f56f`

## Phase 3 — Plugin Extraction ✅ DONE

- `plugins/sample` → `github.com/dreego-stack/plugin-example`
- `plugins/sse` → `github.com/dreego-stack/plugin-sse`
- Each plugin repo: own `go.mod` requiring `github.com/dreego-stack/dreego`, own git history, pushed to GitHub (public)
- `plugins/` removed from main repo; `_tests/test.sh` no longer scans `_tests/plugins`
- `cmd/dreego/go.sum` removed (no longer a separate module)
- Commits: `6357014`, `683f56f`

## Phase 4 — PR-based Release System ✅ DONE (replaces original "Cross-Repo CI" plan)

**Original plan (repository_dispatch + PAT) was replaced by a simpler, token-free design:**

### 4.1 Core repo workflows

- `pull_request.yml` — checks pr.md exists + format valid + `make test`. Runs automatically for in-repo PRs, manually via `workflow_dispatch` for fork PRs (saves Actions minutes).
- `release-prep.yml` — manual `workflow_dispatch` with PR number. Applies pr.md (changelog + version) on the PR branch, removes pr.md, pushes. No branch-protection bypass needed (pushes to feature branch, not main).
- `release.yml` — after merge (`pull_request: closed` + `merged`), creates the tag if VERSION was bumped.

### 4.2 Plugin repo workflows

- Same `pull_request.yml` + `release-prep.yml` + `release.yml`
- `dependabot.yml` — daily gomod updates for `github.com/dreego-stack/dreego`. When core releases a new tag, Dependabot auto-creates an update PR in each plugin repo (tests run via pull_request.yml, then merge updates the require).
- `pull_request.yml` drops the local `replace` directive and resolves the real dreego version.

### 4.3 pr.md format

```yaml
---
version: patch        # none | patch | minor | major
---

- Bug: fix X
- Feat: add Y
```

- `version: none` — no version bump, changelog lines only (e.g. dependabot updates)
- `version: patch` — `0.0.x` +1; `minor` — `0.x.0` +1; `major` — `x.0.0` +1
- `_scripts/release-prep.py` applies pr.md → CHANGELOG + VERSION, removes pr.md
- `pr.md.example` + `.github/PULL_REQUEST_TEMPLATE.md` in every repo

### 4.4 Security

- No secrets, no `pull_request_target`, no cache (`cache: false` in setup-go)
- Branch protection on `main` in all 3 repos: 1 required review, no direct pushes
- `make test` is all-in-one: core-deps check + go test + 185 integration tests (Docker)

### 4.5 Flaky test fix

- `run-timer-sigterm`: `os.Exit(0)` after `fmt.Println` lost the buffered stdout line under parallel load. Fixed with `return` (Go flushes on normal main exit). Commits: `651b24a`, `3adee68`

### 4.6 AGENTS.md + docs updated

- PR-based commit convention, file structure (go.mod, pr.md.example, .github/workflows), feature workflow
- `_scripts/release.sh` removed (obsolete)
- Commits: `46875fb`, `fd46e90`

---

## Phase 5 — Docs Simplification (partially done, decision changed)

**Original plan (remove embedded mirror) was changed:** the embedded docs stay (they work offline, no HTTP fetch exists — `fetchDocFallback` already points to `fetchDocEmbedded`). Only stale references were corrected:

- `_docs/cli.md`: removed `plugins/<name>/_docs/` priority note (plugins are in separate repos now)
- `_docs/config.md`: `plugins/logging` → `plugin-logging`
- `_docs/plugins.md` + `_docs/plugin-interfaces.md`: layout updated to separate repos (done in Phase 2)
- Root `CLI.md` still exists — decide: keep as duplicate or delete (see Phase 8)

## Phase 6 — `_todo` / Blockwebchain Simplification (open)

- Drop integer-chain status, keep only `planned | in-progress | done | rejected`
- Done blocks stay in `blocks/` with `status: done` (no move to `history/`)
- Simplify `process.py`: keep dependency graph + cycle detection + "available next", drop chain validation
- Update `_todo/Blockwebchain.skill.md` + AGENTS.md

## Phase 7 — Test Migration (shell → _test.go) (open, incremental)

- New behavior tests: `_test.go` using `dreegotest` + `httptest` + golden files
- Shell tests: keep only for true CLI/subprocess E2E (`dreego new`, `dreego init`, `dreego dev`, `dreego fmt`, deployment)
- Migrate existing shell tests incrementally (not all at once)
- Inline test.sh rules into AGENTS.md + CI check for standard header

## Phase 8 — Cleanup Chores (open)

- `_tests/TODO.nd` — typo file, delete
- Root `CLI.md` — duplicate of `_docs/cli.md`, decide: delete
- `.gitignore`: remove `go.work`/`go.work.sum` entries (no longer relevant)
- `.review/` — confirm ignored
- `demo/` — keep for now; consider moving to `_docs/examples/` or own repo later

---

## Execution Order Summary

```
Phase 1 (platform)    ✅ done — 245 files, scripted replace, verified
Phase 2 (single mod)  ✅ done — root go.mod, one tag, fresh clone builds
Phase 3 (plugins out) ✅ done — 2 repos created + pushed
Phase 4 (PR-based CI) ✅ done — pull_request/release-prep/release + dependabot
Phase 5 (docs)        🟡 partial — stale refs fixed, embedded stays
Phase 6 (_todo)       ⬜ open
Phase 7 (tests)       ⬜ open (incremental, spans versions)
Phase 8 (cleanup)     ⬜ open
```

**Phases 1–4 are complete.**
Phases 5–6 can be v0.0.27 or v0.0.28.
Phase 7 is ongoing across multiple versions.
Phase 8 is batch cleanup at the end of v0.0.27.

---

## Risks

- **Phase 1:** 333-file replace verified with `go build` + `make test`. Old codeberg tags stay frozen on the codeberg repo (read-only archive) for pinned users.
- **Phase 2:** core-deps enforcement shifted from compile-time module boundary to `_scripts/check-core-deps.sh` in CI.
- **Phase 3:** Plugin extraction is irreversible (git history split). Main repo tags remain as frozen reference.
- **Phase 4:** Dependabot needs ~24h for the first update check. Branch protection requires a manual approve + `release-prep` run per PR.
- **Phase 5:** Embedded docs stay — no offline regression. Root `CLI.md` duplication remains until Phase 8.
