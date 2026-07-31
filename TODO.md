# Blockwebchain

Status via `python _todo/process.py`. Chain 1–34 done. Next code: **35**.

Versioning stays conservative: `v0.x.y` only. `v0.x` marks larger milestones (may include breaking changes), `y` is continuous. `v1.0.0` is reserved for a stable, trustworthy release and is not a near-term target.

## v0.0.17 (current)

- **deployment.1** — Production Deployment (Graceful Shutdown, Cross-Compile, Docker) ✅
- **request-id.1** — Request-ID Middleware (X-Request-ID → Context + Logs) ✅

## v0.0.16

- **form-actions.1** — Form Actions (g-action + auto-validation + redirect) ✅

## v0.0.15

- **api-json.1** — Content-Type Routing (JSON, XML, Custom) ✅

## v0.0.14

- **health-checks.1** — /health + /ready Endpoints ✅
- **security-headers.1** — Security Headers ✅
- **compression.1** — Gzip Compression Middleware ✅

## Available Next

- **dreegotest.1** — Testing Package
- **observability.1** — Metrics + Tracing (Plugin: Prometheus, OpenTelemetry)
- **documentation.1** — docs.dreego.dev + Tutorial + Examples

## Rejected

- **hot-reload.1** — Replaced by Air (`_docs/hot-reload.md`)
- **live-reload.1** — Replaced by Air
- **smart-recompile.1** — Replaced by Air

## Quality Backlog (from code review)

- **codegen-errors.1** — Replace silent fails in CodeGen with explicit errors.
- **servemux-cache.1** — Build and cache the middleware/routing stack once at startup.
- **security-cookie.1** — Harden session and CSRF cookie flags.
- **security-csp.1** — Add Content-Security-Policy header.
- **security-session.1** — Document or encrypt sensitive session payload.
- **golden-tests-core.1** — Add golden-code tests for generated Go output.
- **docs-extensibility.1** — Adapt `dreego docs` to consume docs from plugins/external repos.

## Framework Roadmap (new blocks)

- **dev-server.1** — `dreego dev` with file watcher and auto-regenerate.
- **frontmatter.1** — Parse frontmatter in `.dreego` and expose typed metadata.
- **typed-forms.1** — Extend form binding/validation to int, bool, slices, custom validators, and improve `email` validation.
- **client-reactivity.1** — Research client-side interactivity (Alpine/islands/custom runtime).

See `_todo/index.md` for the full chain and dependency graph.
