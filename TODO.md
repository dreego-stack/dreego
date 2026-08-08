# Blockwebchain

Status via `python _todo/process.py`. Chain 1–49 done. Next code: **50**.

Versioning stays conservative: `v0.x.y` only. `v0.x` marks larger milestones (may include breaking changes), `y` is continuous. `v1.0.0` is reserved for a stable, trustworthy release and is not a near-term target.

## v0.0.25 (done)

- **plugin-interface.1** — frozen v1 `core.Plugin` contract + `core.UsePlugin` ✅
- **middleware-hooks.1** — plugin middleware via `app.Use` (FIFO chain) ✅
- **route-hooks.1** — programmatic plugin route registration ✅
- **docs-extensibility.1** — `dreego docs` design for plugin docs ✅
- **docs-embed.1** — offline embedded docs (`cmd/dreego/embedded/`) ✅
- **frontmatter.1** — YAML frontmatter parsing + typed metadata ✅
- **dev-server.1** — `dreego dev` with file watcher + auto-regenerate ✅

## v0.0.24 (done)

- **scaffold-fix.1** — `dreego new` go.sum + `.gitignore` only `dreego/gen/` ✅
- **layout-head.1** — layouts apply, route `<head>` with/without layout ✅
- **scoped-css.2** — `scopeCSS` preserves `{}` declarations, `@media`, `@keyframes` ✅
- **component-attr-props.1** — `{prop}` in HTML attributes substituted + escaped ✅
- **typed-forms.1** — int/bool/slice binding + `RegisterRule` custom validators ✅
- **dreegotest.1** — exported `dreegotest` test helper package ✅
- **golden-tests-core.1** — golden-file assertions for generated Go ✅
- **port-schema / test stability** — deterministic runner ports, DREEGO_BIN fallbacks ✅

## v0.0.23 (done)

- **feedback-intake A/B** — nested `{#if}` in `{#else}` + `<head>` expression resolution ✅
- **security-session.2/.3**, **core.Reset()** ✅

## v0.0.22 (done)

- **servemux-cache.1** — cached built middleware/router stack ✅
- **codegen-errors.1** — all CodeGen returns `(string, error)` ✅
- **security-session.1** — optional AES-256-GCM session encryption ✅

## v0.0.21 (done)

- **monorepo-plugin-layout** — Official plugins moved into `plugins/` in this repo (chain via v0.0.21 commit)

## v0.0.20

- **security-cookie.1** — Harden session and CSRF cookie flags ✅ (chain 35)
- **security-csp.1** — Add Content-Security-Policy header ✅ (chain 36)

## v0.0.17

- **deployment.1** — Production Deployment (Graceful Shutdown, Cross-Compile, Docker) ✅
- **request-id.1** — Request-ID Middleware (X-Request-ID → Context + Logs) ✅

## v0.0.16

- **form-actions.1** — Form Actions (g-action + auto-validation + redirect) ✅

## Available Next

- **observability.1** — Metrics + Tracing (Plugin: Prometheus, OpenTelemetry)
- **documentation.1** — docs.dreego.dev + Tutorial + Examples

Planned for **v0.0.26**: `documentation.1`, `api-swagger.1`, `observability.1`.

## Rejected

- **hot-reload.1** — Replaced by Air (`_docs/hot-reload.md`)
- **live-reload.1** — Replaced by Air
- **smart-recompile.1** — Replaced by Air

## Quality Backlog (from code review)

- **codegen-errors.2** — ✅ already fixed in v0.0.22 (codegen-errors.1, commit 11d33d2); backlog entry obsolete, covered by `TestCompGenIfElseMixedChildren`

## Framework Roadmap (new blocks)

- **client-reactivity.1** — Research client-side interactivity (Alpine/islands/custom runtime).

See `_todo/index.md` for the full chain and dependency graph.