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

## Open work

### Phase 5 — Docs Simplification

- Root `CLI.md` — duplicate of `_docs/cli.md`, decide: keep or delete

### Phase 6 — Todo System Simplification ✅ DONE

- `_todo/` (Blockwebchain, process.py, blocks/) deleted
- Replaced by `TODO.md` (planned code work, checkbox entries + how-to header) and `TODO-Future.md` (SSG, Wails, plugin ideas)
- Plugin blocks, SSG/Wails/client-reactivity, and interface ideas reviewed with the user and routed to TODO-Future.md or dropped
- References updated: AGENTS.md, .agents/index.md, .agents/guides/open-knowledge-format.md, .gitattributes

### Phase 7 — Test Migration (shell → _test.go)

- New behavior tests: `_test.go` using `dreegotest` + `httptest` + golden files
- Shell tests: keep only for true CLI/subprocess E2E (`dreego new`, `dreego init`, `dreego dev`, `dreego fmt`, deployment)
- Migrate existing shell tests incrementally (not all at once)
- Inline test.sh rules into AGENTS.md + CI check for standard header

### Phase 8 — Cleanup Chores

- `_tests/TODO.nd` — typo file, delete
- Root `CLI.md` — delete if Phase 5 decides to remove
- `.gitignore`: remove `go.work`/`go.work.sum` entries (no longer relevant)
- `.review/` — confirm ignored
- `demo/` — keep for now; consider moving to `_docs/examples/` or own repo later

## Execution Order Summary

- Phase 5 (docs)        partial — only `CLI.md` decision remains
- Phase 6 (todo)        done
- Phase 7 (tests)       open (incremental, spans versions)
- Phase 8 (cleanup)     open

Phase 7 is ongoing across multiple versions. Phases 5, 6 and 8 can land in v0.0.27 or be bumped to v0.0.28.

## Risks

- **Phase 5:** Embedded docs stay — no offline regression. Root `CLI.md` duplication remains until Phase 8.
- **Phase 6:** done — no residual risk.
- **Phase 7:** Large migration; doing it incrementally avoids a big-bang rewrite.
- **Phase 8:** `demo/` move is deferred; do not touch unless time permits.
