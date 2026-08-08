# v0.0.25 — Finalization

Task: finalize v0.0.25 after all 7 feature increments committed.

## Changed files

- `VERSION` — bumped `v0.0.24` → `v0.0.25`
- `_todo/plan.md` — v0.0.25 moved into `## Done` (✅); `## v0.0.26` + plugin-ecosystem `unlock` untouched; notes updated (docs extensibility decided → embedded, plugin docs deferred)
- `TODO.md` — added `## v0.0.25 (done)` with all 7 blocks; removed the now-implemented deferred list; noted `documentation.1`/`api-swagger.1`/`observability.1` for v0.0.26; kept quality backlog + roadmap
- `CHANGELOG.md` — extended v0.0.25 entry with `middleware-hooks.1`, `route-hooks.1`, `docs-extensibility.1`, `docs-embed.1`, `frontmatter.1`, `dev-server.1`; suite count → **164 passed, 0 failed**
- `_docs/plugins.md` — added "Middleware hooks (FIFO)" + "Route hooks (programmatic routes)" sections
- `_docs/cli.md` — documented `dreego dev`; noted offline embedded docs + local plugin-doc priority
- `_docs/index.md` — added Frontmatter, Plugins, Dev Server links
- `_docs/frontmatter.md` — **new**: YAML frontmatter syntax + `ParseFrontmatter` API
- `_docs/dev-server.md` — **new**: `dreego dev` usage, polling, restart semantics
- `cmd/dreego/embedded/` — regenerated via `_scripts/sync-embedded-docs.sh`

## Verification

- `go test ./core/... ./cmd/dreego/... ./dreegotest/...` → ok (all 3 modules)
- `sh _tests/test.sh` → **164 passed, 0 failed**

No commit made (per instruction).
