# Plan v0.0.27 — Chore Batch: Platform Migration + Repo Restructure

**Type:** chore-only (no features, no breaking behavior changes for end users)
**Status:** planning
**Scope:** platform move, single-module restructure, plugin extraction, CI, docs, test migration, cleanup

## Decisions (locked by user)

1. **Platform:** GitHub, org `dreego-stack`. Main repo: `github.com/dreego-stack/dreego`
2. **Plugins:** Each official plugin gets its own repo under `github.com/dreego-stack/<plugin-name>`
3. **v0.0.27 = chores only** — no new features, no behavior changes

---

## Phase 1 — Platform Migration (codeberg → github.com/dreego-stack)

**Why first:** Every other phase depends on the final module path. Doing this last means double-work.

### 1.1 Global module path rename

- `codeberg.org/dreego/dreego` → `github.com/dreego-stack/dreego` everywhere
- Affected files (333 total):
  - 15 `go.mod` files (module declarations + require/replace lines)
  - 30 `.go` files (import paths in codegen output, scaffold templates, docs.go URLs)
  - 254 `.sh` test files (require/replace in heredoc go.mod)
  - 34 `.md` files (docs, embedded docs, README, CHANGELOG)
- Key locations to verify:
  - `cmd/dreego/docs.go:14-16` — `docsBaseURL`, `docsWebBase`, `feedbackURL`
  - `cmd/dreego/new.go` — scaffold templates (go.mod, main.go generated for users)
  - `cmd/dreego/fmt.go`, `cmd/dreego/main.go`, `cmd/dreego/dev.go` — imports
  - `core/generate.go` — codegen output import paths
  - `core/fmt_test.go` — test fixtures
  - `_tests/how-to-test-sh.md:22-23` — test template (require + replace)
  - `demo/go.mod`, `demo/main.go`
  - All `_tests/**/test.sh` heredocs
- **Method:** scripted `find + sed` replace across tracked files, then `go build ./...` + `make test` to verify
- `go.work.sum` is gitignored — will regenerate
- **NOTE:** Phase 1 and Phase 2 must be executed as a **coordinated migration** for shell test files — see Phase 2.1 for the combined sed step that handles both the path rename and the replace-structure change in one pass. Running them as independent passes will break all 254 shell tests.

### 1.2 Git remote + mirror

- Change `origin` to `git@github.com:dreego-stack/dreego.git`
- Optionally keep `codeberg` as a push-mirror remote (user decides)
- Push all branches + existing tags
- **Keep codeberg repo as a read-only frozen archive.** Do not delete the repo or remove old tags. Existing `codeberg.org/dreego/dreego/core@vX` pins resolve against these frozen tags indefinitely.

### 1.3 Tag strategy for the migration

- Existing tags (`core/v0.0.26`, `cmd/dreego/v0.0.26`, `plugins/sample/v0.0.26`) stay as historical artifacts on the old path
- After the single-module restructure (Phase 2), new tags are just `v0.0.27`, `v0.0.28`, ...
- Old module-path users who pinned `codeberg.org/dreego/dreego/core@v0.0.26` keep working against the old tags (frozen). New users use `github.com/dreego-stack/dreego`

---

## Phase 2 — Single Root Module (enables single tag)

**Why second:** Module structure determines how plugins, CI, and tags work. Must be stable before extracting plugins.

### 2.1 Root go.mod

- Create `/go.mod`: `module github.com/dreego-stack/dreego`
- Delete `core/go.mod`, `cmd/dreego/go.mod`, `dreegotest/go.mod`
- Delete `go.work` + `go.work.sum` (single module = no workspace needed)
- Import paths stay identical: `github.com/dreego-stack/dreego/core`, `.../cmd/dreego`, `.../dreegotest`
- `go install github.com/dreego-stack/dreego/cmd/dreego@v0.0.27` works with a single tag
- **Coordinated shell-test migration (Phase 1 + Phase 2 combined):** The 254 shell test files contain heredoc `go.mod` blocks that must be updated for BOTH the path rename and the replace-structure change in a single pass:
  - `require codeberg.org/dreego/dreego/core v0.0.0` → `require github.com/dreego-stack/dreego v0.0.0`
  - `replace codeberg.org/dreego/dreego/core => $realrepo/core` → `replace github.com/dreego-stack/dreego => $realrepo`
  - Update `_tests/how-to-test-sh.md` template accordingly (require + replace + `$realrepo` depth)
  - `$realrepo` path depth stays the same (still points to repo root, where the new root `go.mod` lives)
- **`demo/` decision:** Keep `demo/` as a standalone module with its own `go.mod`. Update in Phase 1: `require github.com/dreego-stack/dreego v0.0.27` (or `v0.0.0` + `replace => ..` for local dev) + update `demo/main.go` import path. Users can copy `demo/` as a starting project template.

### 2.2 Update release.sh

- From: loop over `core cmd/dreego plugins/sample` creating N tags
- To: single `git tag $(cat VERSION)` + push
- Remove all module-prefix tag logic

### 2.3 Update AGENTS.md

- Remove "own module" rules for core/cmd/dreegotest
- Keep "core has no external deps" as a convention (not a module boundary)
- Keep "max 300 lines per file" rule
- Update build instructions: `go build ./cmd/dreego` directly (no go.work)
- **Core-deps enforcement:** Add a CI check that verifies `core/` has no imports outside the standard library and the dreego module itself. CI step: `go list -deps ./core/ | grep -v '^github.com/dreego-stack/dreego' | grep -v std` (or a dedicated `_scripts/check-core-deps.sh`) fails if core picks up external deps. This replaces the compile-time module-boundary enforcement lost by merging `core/go.mod` into the root module.

### 2.4 Verify

- `go build ./...` from root
- `go test ./core/... ./cmd/dreego/...`
- `make test` (integration tests still work — they create their own go.mod with replace)
- Fresh clone test: `git clone ... && go build ./cmd/dreego` works without go.work

---

## Phase 3 — Plugin Extraction

**Why third:** Module path is final (Phase 1), root is single-module (Phase 2). Now plugins can be extracted cleanly.

### 3.1 Create `github.com/dreego-stack/plugin-example`

- Move `plugins/sample/` content to new repo `github.com/dreego-stack/plugin-example`
- `go.mod`: `module github.com/dreego-stack/plugin-example`, `require github.com/dreego-stack/dreego v0.0.27`
- This is the reference plugin — every other plugin copies its structure
- Add `README.md` with: how to use, how to test against core, CI setup reference

### 3.2 Create `github.com/dreego-stack/plugin-sse`

- Move `plugins/sse/` content to `github.com/dreego-stack/plugin-sse`
- Same pattern as plugin-example

### 3.3 Remove plugins/ from main repo

- Delete `plugins/` directory from main repo
- Remove `plugins/` entries from `go.work` (already deleted in Phase 2)
- Update `AGENTS.md`: "Official plugins live in separate repos under `github.com/dreego-stack/`"
- Update `_todo/plan.md`: remove "Official plugins live in `plugins/`" references
- Update `_docs/plugins.md`: point to separate repos

### 3.4 Update docs.go plugin path resolution

- `cmd/dreego/docs.go` currently resolves `plugins/<name>/_docs/` locally
- With plugins in separate repos, local resolution no longer works
- Options:
  - (a) `dreego docs plugins/<name>` fetches from the plugin repo's raw URL
  - (b) Drop plugin-docs resolution entirely; plugin repos have their own `dreego docs`-compatible structure
- Decision needed during implementation (low risk, can be deferred)

---

## Phase 4 — Cross-Repo CI (compatibility tests)

**Why fourth:** Plugins are extracted (Phase 3). Now CI must catch core-breaks-plugin.

### 4.1 Tag-triggered compat tests (primary)

- In **each plugin repo**: GitHub Actions workflow triggered by `repository_dispatch` or `workflow_dispatch`
- Core repo release workflow (on tag push) sends `repository_dispatch` to each plugin repo with `{"version": "v0.0.27"}`
- Plugin repo workflow: `go get github.com/dreego-stack/dreego@v0.0.27 && go test ./...`
- Failure → e-mail notification (GitHub Actions default) + issue auto-created in core repo

### 4.2 Daily compat sweep (safety net)

- In core repo: scheduled workflow (daily, cron `0 6 * * *`)
- Checks out each plugin repo, runs `go get github.com/dreego-stack/dreego@latest && go test ./...`
- Catches breakage from core `main` branch pushes (before a tag is cut)
- Failure → e-mail + issue

### 4.3 Core repo CI

- GitHub Actions on push/PR: `go build ./...` + `go test ./...` + `make test` (Docker-based integration tests)
- On tag push: build + release workflow (triggers plugin compat tests via 4.1)

### 4.4 Token / secret setup

- `DREEGO_DISPATCH_TOKEN` — **fine-grained GitHub PAT** scoped to specific plugin repos (`dreego-stack/plugin-example`, `dreego-stack/plugin-sse`) with `Actions: write` + `Contents: read` permissions only. Do not use classic PATs with broad `repo` scope.
- Or use a dedicated GitHub App (more secure, more setup) — defer to later

---

## Phase 5 — Docs Simplification

**Why fifth:** Low risk, no dependencies on other phases. Can be done anytime after Phase 1.

### 5.1 Remove embedded docs mirror

- Delete `cmd/dreego/embedded/` directory (18 `_docs/` copies + README + CHANGELOG)
- Delete `cmd/dreego/embed.go` (the `//go:embed` directive)
- Delete `_scripts/sync-embedded-docs.sh`
- Delete `cmd/dreego/docs_embed_test.go`
- Binary gets smaller; drift source eliminated

### 5.2 Slim down `dreego docs`

- `dreego docs <path>` → fetches from GitHub raw URL (no local embed fallback)
- `dreego docs --web` → opens GitHub source page (unchanged)
- `dreego docs --dump` → fetches multiple pages from remote
- Offline → prints URL hint ("fetch docs online or visit https://github.com/dreego-stack/dreego")
- Update `cmd/dreego/docs_test.go` to match new behavior (mock remote fetch or skip when offline)

### 5.3 CLI help = source of truth for CLI reference

- `dreego help` / `dreego --help` stays in Go code (always correct)
- Delete root `CLI.md` (duplicate of `_docs/cli.md`)
- Keep `_docs/cli.md` as human-readable reference (generated from or verified against `dreego help`)

### 5.4 Update docs URLs

- `docs.go`: `docsBaseURL` → `https://raw.githubusercontent.com/dreego-stack/dreego/main`
- `docs.go`: `docsWebBase` → `https://github.com/dreego-stack/dreego/blob/main`
- `docs.go`: `feedbackURL` → `https://github.com/dreego-stack/dreego/issues/new`

---

## Phase 6 — `_todo` / Blockwebchain Simplification

**Why sixth:** Agent reliability improvement. No code dependencies.

### 6.1 Drop the integer-chain system

- Remove chain-status integers (1–49) from all block frontmatter
- Status becomes only: `planned | in-progress | done | rejected`
- Done blocks stay in `blocks/` with `status: done` (no move to `history/`)
- Delete `blocks/history/` directory (merge done blocks back into `blocks/` or archive)

### 6.2 Simplify process.py

- Remove: chain validation, gap detection, duplicate-integer checks, next-code computation
- Keep: dependency graph, cycle detection, "available next" (all requires = done), validation (requires resolve)
- Output: available / in-progress / blocked / done / rejected (no chain numbering)

### 6.3 Update AGENTS.md + skill doc

- Remove chain-integer ritual from workflow descriptions
- AGENTS.md: "Mark block `status: done` when complete. Run `python _todo/process.py` to verify graph."
- Update `_todo/Blockwebchain.skill.md` to reflect simplified system

---

## Phase 7 — Test Migration (shell → _test.go)

**Why seventh:** Largest effort, incremental, no blockers. Can span multiple versions.

### 7.1 Principle

- New behavior tests: write as `_test.go` using `dreegotest` + `httptest` + golden files
- Shell tests (`_tests/**/test.sh`): keep only for true CLI/subprocess E2E (`dreego new`, `dreego init`, `dreego dev`, `dreego fmt`, deployment)
- Migrate existing shell tests incrementally (not all 189 at once)

### 7.2 Migration candidates (high value, low effort)

- `Transpiler/*` → `_test.go` in `core/` (codegen assertions, already have golden tests)
- `Components/*` → `_test.go` in `core/` (component rendering)
- `Middleware/*` → `_test.go` in `core/` (middleware behavior via `dreegotest`)
- `Routing/*` → `_test.go` in `core/` (route resolution)
- `Session/*` → `_test.go` in `core/` (session behavior)
- `Layout/*` → `_test.go` in `core/` (layout merge)

### 7.3 Shell tests to keep

- `CLI/*` (init, new, version-flag, docs) — true CLI subprocess tests
- `Deployment/*` — Docker/cross-compile
- `Bugs/*` — keep as regression markers (can add `_test.go` alongside)

### 7.4 Update how-to-test-sh.md + AGENTS.md

- Inline core test.sh rules into `AGENTS.md` (condensed: mktemp, trap, $DREEGO_BIN, port, header)
- Add CI check: validate every `_tests/**/test.sh` has the standard header (`# Using standard:` + `# What:`)

---

## Phase 8 — Cleanup Chores

**Why last:** Small, independent, no risk. Batch together.

### 8.1 Delete stale files

- `_tests/TODO.nd` — typo file (should be .md, content stale), delete
- `cmd/dreego/embedded/CHANGELOG.md` — embedded copy, deleted with Phase 5
- Root `CLI.md` — duplicate, deleted with Phase 5

### 8.2 Fix plugin inconsistencies (pre-extraction)

- `plugins/sse/go.mod` uses `core v0.0.0` + `replace`; `plugins/sample/go.mod` uses `core v0.0.26`
- Moot after Phase 3 extraction (each plugin repo pins its own core version)

### 8.3 Review tracked files

- `.review/` — not tracked (gitignored or empty). Confirm and add to `.gitignore` if needed
- `demo/` — tracked, 8 files. Keep for now; consider moving to `_docs/examples/` or own repo later
- `.DS_Store` — in `.gitignore`, not tracked. Confirm none slip in

### 8.4 .gitignore updates

- Remove `go.work` / `go.work.sum` from `.gitignore` (no longer relevant after Phase 2)
- Add `.review/` if not already ignored

### 8.5 AGENTS.md updates

- Update all codeberg references → github.com/dreego-stack
- Update file structure diagram (no plugins/, no go.work, no embedded/)
- Update coding rules (single module, single tag)
- Update commit convention (single tag via `git tag $(cat VERSION)`)

---

## Execution Order Summary

```
Phase 1 (platform)    →  333 files, scripted replace, verify build+test
Phase 2 (single mod)  →  delete go.work + sub-module go.mods, create root go.mod
Phase 3 (plugins out) →  create 2 repos, move content, delete plugins/
Phase 4 (CI)          →  GitHub Actions: core CI + tag-trigger + daily sweep
Phase 5 (docs)        →  delete embedded/, slim dreego docs, fix URLs
Phase 6 (_todo)       →  drop chain integers, simplify process.py
Phase 7 (tests)       →  incremental shell→_test.go (can span versions)
Phase 8 (cleanup)     →  stale files, .gitignore, AGENTS.md
```

**Phases 1–4 are the critical path for v0.0.27.**
Phases 5–6 can be v0.0.27 or v0.0.28.
Phase 7 is ongoing across multiple versions.
Phase 8 is batch cleanup at the end of v0.0.27.

---

## Risks

- **Phase 1:** 333-file replace is mechanical but must be verified with `go build` + `make test`. A missed import path breaks silently in codegen output (only visible at user build time). Mitigation: full test suite must pass. Phase 1 and Phase 2 must be coordinated for shell test files (see Phase 2.1).
- **Phase 2:** `dreegotest` and `demo` currently use `replace` directives — with single module these become unnecessary for `dreegotest` (merged into root). `demo/` stays standalone with its own `go.mod` (see Phase 2.1 decision). Core-deps enforcement shifts from compile-time module boundary to a CI check (see Phase 2.3).
- **Phase 3:** Plugin extraction is irreversible (git history split). Mitigation: use `git filter-repo` or just copy files (plugins have no deep history). Keep main repo tags as frozen reference.
- **Phase 4:** Cross-repo CI needs a fine-grained PAT or GitHub App. Fine-grained PAT is simpler and scoped; App is more secure. Defer App to post-v0.1.
- **Phase 5:** Removing embedded docs breaks offline `dreego docs`. Acceptable trade-off (CLI help covers CLI reference; full docs are online).