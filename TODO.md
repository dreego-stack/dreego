# TODO

Concrete, planned code work. Ideas without a near-term plan live in [TODO-Future.md](TODO-Future.md).

## How to use this file

- One line per entry, checkbox for progress:
  `- [ ] **name** — short description`
- Status via the checkbox: `[ ]` planned / `[x]` done
- Group entries under a version heading (`## v0.0.x`) or a topic heading
- When an entry ships, tick the box and add a changelog line via `pr.md`
- Versioning stays conservative: `v0.x.y` only. `v1.0.0` is reserved for a stable release.

## Planned

### Core

- [x] **event-bus.1** — Core Pub/Sub Event Bus interface (abstracts Redis/NATS/In-Memory), typed via generics: Publish, Subscribe, Unsubscribe
- [ ] **queue-interface.1** — Core Background Job Queue interface (abstracts Redis/NATS/In-Memory): job middleware, batching, chaining, delayed dispatch
- [ ] **storage-interface.1** — Core File Storage interface (S3/R2/Local): Put, Get, Delete, List, URL — interface only, like `database/sql`
- [x] **kv-store.1** — Core Key-Value Store interface (abstracts Redis/Ristretto/In-Memory): Get, Set, Delete, Expire — interface only, like database/sql, distinct from Storage (blobs), small values with TTL

### Plugins (external repos)

- [ ] **observability.1** — Metrics + Tracing: Prometheus `/metrics` + OpenTelemetry spans as plugins (separate repos, own go.mod); plugin-interface.1 (v0.0.25) is the foundation
- [ ] **api-swagger.1** — Auto-generated OpenAPI 3.0 spec as plugin: `api:"..."`/`validate:"..."` struct tags on routes, `/openapi.json` endpoint, optionally embedded Swagger UI

### Decision needed

- [x] **addon vs plugin naming** — decided: "plugin" (v0.0.27) — `TODO-Future.md` uses "plugin" consistently, concept docs renamed accordingly

## v0.0.27 (in progress — chore batch)

- Platform migration codeberg → github.com/dreego-stack ✅
- Single root module (one tag per release) ✅
- Plugin extraction to separate repos ✅
- PR-based release workflow (pr.md + CI tag) ✅
- AGENTS.md + docs updated for new workflow ✅
- Flaky test fix (run-timer-sigterm stdout flush) ✅
- Todo system simplification (_todo/ → TODO.md + TODO-Future.md) ✅
- Test migration strategy (shell → _test.go) ✅
- [x] AGENTS.md test rules — update Feature Workflow to `_test.go` + `dreegotest`, drop `test.sh` reference
- [x] CI standard-header check — add a CI check for the standard test header in workflows

## v0.0.26 (done)

- Post-v0.0.25 bugfix batch — `{#if}/{#each}` in attributes, `<...>` in go strings, sections after leading text, doctype/scope in error pages, head dedupe, prop defaults, JSON/XML status-before-encode, `dreego init` gen import, `plugins/sample` core-require drift
- **SSE example plugin** — full v1 `core.Plugin` (`plugins/sse`), streaming fixed through `RequestLogging`/`Compress` via `http.Flusher`
- **`--version` / `-v` flags** — match the help trio

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

## v0.0.20 (done)

- **security-cookie.1** — Harden session and CSRF cookie flags ✅ (chain 35)
- **security-csp.1** — Add Content-Security-Policy header ✅ (chain 36)

## v0.0.17 (done)

- **deployment.1** — Production Deployment (Graceful Shutdown, Cross-Compile, Docker) ✅
- **request-id.1** — Request-ID Middleware (X-Request-ID → Context + Logs) ✅

## v0.0.16 (done)

- **form-actions.1** — Form Actions (g-action + auto-validation + redirect) ✅

## Rejected / superseded

- **hot-reload.1**, **live-reload.1**, **smart-recompile.1** — replaced by Air (`_docs/hot-reload.md`)
- **documentation.1** (docs.dreego.dev + Tutorial) — replaced by the decentralized `dreego docs` system
- **golden-tests.1** — already covered by `golden-tests-core.1` (v0.0.24)
