# Changelog

## v0.0.19 (2026-07-30) — Bug Fixes

- Fix B1: `{#if}` and `{#each}` now transpile correctly inside components **and** route templates
- Fix B2: `BindForm` returns an error instead of panicking on non-string fields
- Lexer: `{` treated as template control-flow everywhere except inside `<go>`, `<head>`, `<script>`, and `<style>` sections
- Lexer: arbitrary HTML tags (e.g. `<ul>`, `<input>`) tokenize without mandatory balancing
- Codegen: `NodeIf`/`NodeEach` cases added to component template generation path
- Parser: unknown root tags and arbitrary open/close tags allowed in template contexts
- Test: `_tests/Bugs/component-if-each` covers `{#if}`/`{#each}` in components
- Test: `_tests/Bugs/bindform-non-string` covers B2
- Test: `_tests/Template/each-loop` converted from expected-failure to expected-success for route-level `{#each}`

## v0.0.18 (2026-07-29) — Package Restructuring

- **BREAKING**: `dreego-core/` → `core/` — import path changes from `codeberg.org/dreego/dreego/dreego-core` to `codeberg.org/dreego/dreego/core`
- **BREAKING**: `dreego-plugin/` removed — plugins live in separate repos under `codeberg.org/dreego/<name>`
- `_docs/ plugins.md`: plugin architecture overview, planned plugins, interface contracts
- AGENTS.md updated to reflect new directory structure

## v0.0.17 (2026-07-29) — Production Deployment + Request-ID

- **Graceful Shutdown**: `core.Listen()` uses `http.Server` with SIGINT/SIGTERM handling, 10s drain timeout
- **Cross-Compile**: `dreego build --target linux/amd64` sets GOOS/GOARCH for target platform
- **Request-ID Middleware**: `X-Request-ID` header on every request — client-supplied or auto-generated (16-char hex), injected into context + JSONL logs (`rid` field), accessible via `c.RequestID()`
- **Production Dockerfile**: `FROM scratch` — 3-stage build, `CGO_ENABLED=0`, static binary
- **Hot Reload**: `_docs/hot-reload.md` — Air config with `.air.toml` + `entr` alternative
- **Rejected**: hot-reload.1, live-reload.1, smart-recompile.1 — replaced by Air documentation
- Block: request-id.1 completed (chain 34)

## v0.0.16 (2026-07-29) — Form Actions

- `<form g-action="Login">` — declarative server-side form handling with auto-generated pipeline
- Generated POST handler: `r.ParseForm()` → struct mapping via `form:"email"` tags → validation via `validate:"required,email"` tags → handler call → redirect
- `c.Redirect(url, code)` — PRG pattern (Post-Redirect-Get) with `ErrRedirect` sentinel
- `c.Errors(field)` / `c.Old(field)` — validation error and old value access in templates
- `BindForm(r, target)` — maps form values to struct fields (explicit `form:` tag or lowercase field name)
- `ValidateForm(form)` — validates struct via `validate:` tags (required, email, min, max) — no external deps
- `SaveErrors(c, errs)` / `SaveOld(c, form)` — automatically stores validation state for template re-render
- Codegen: `splitGoSections` separates type/func declarations from inline code — form structs and handlers at package level
- Context interface extended: `SessionVal`, `SetSessionVal`, `DelSessionVal`, `CSRFToken`, `Redirect`
- `scanFormActions` detects `<form g-action>` in templates, wires matching handlers in POST dispatch
- 15 new tests: 11 parser/codegen + 4 runtime HTTP (valid submit, invalid re-render, CSRF, plain form)
- 112 integration tests total

## v0.0.15 (2026-07-28) — Content-Type Routing

- `<go type="json">` — JSON endpoints with `c.JSON()`, `c.Bind()`, auto Content-Type
- `<go type="xml">` — XML endpoints with `c.XML()`, auto Content-Type
- `<go type="custom">` — developer manages Content-Type + response via `c.Write()`
- Multiple `<go>` blocks merged: shared `<go>` runs always, typed blocks check `Accept` header
- Pure JSON/XML routes (no `<div>`) skip template rendering entirely
- `c.Write(status, contentType, body)` for arbitrary formats (FlatBuffers, Protobuf, etc.)
- `c.Wants(mime)` for manual content negotiation
- `core/response.go`: `JSON()`, `XML()`, `Bind()`, `Write()`, `Wants()`
- 6 new tests (87 total)

## v0.0.14 (2026-07-28) — Production Middleware

- `GET /health` → 200 `ok` — process liveness probe, always available
- `GET /ready` → 200 `ready` / 503 `not ready` — traffic readiness via `core.SetReady(bool)`
- Health endpoints registered before user routes (cannot be overridden)
- Security Headers middleware: `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` — core-fixed
- Gzip compression middleware: compresses responses when client accepts gzip (`Accept-Encoding`) — core-fixed
- Middleware chain: Recovery → SecurityHeaders → Compression → RequestLogging → Session → CSRF → Redirect/Rewrite → Router
- 3 new tests: health-checks, security-headers, compression
- Fix: `layout.Head.Content` now rendered in codegen — `{#head}` inside `<head>` section processed correctly
- Fix: `genLayoutNode` handles both `{#head}` and `{#slot}` in same NodeText (was dropping `{#slot}` after `{#head}`)
- 3 new regression tests: layout-head-lost, layout-head-ok, head-dropped-by-layout (content verification)
- 81 integration tests total

## v0.0.13 (2026-07-28) — Scaffolding + Split-Gen

- `dreego new <name>` — creates landing page project with `go mod init` auto-setup
- Landing blueprint: Hero + FeatureCard components, layout with `{#head}` + `{#slot}`, Tailwind CDN
- Split-Gen: `gen/routes.go` + `gen/components.go` + `gen/dree.go` (config + static) — file-level caching
- `isUpToDate()` — files only written when content changes
- `.gitignore` fix: `/dreego` → only root binary; `dreego/gen/` covers all generated files
- `dreego generate` timing: ns → ms display
- 5 new tests: head-without-layout, head-dropped-by-layout, name-clash, nested-routes, fmt
- 74 integration tests total
- Tags v0.0.1–v0.0.13 on Codeberg

## v0.0.12 (2026-07-28) — dreego fmt

- `dreego fmt` — formats `.dreego` files in-place (section ordering, template expressions, control flow)
- `dreego fmt --check` — CI mode: exit non-zero if files need formatting
- `dreego fmt --stdout` — print formatted output to stdout
- `core.Format()` — reusable formatting function
- Formats: `{ var }` → `{var}`, `{#if  cond}` → `{#if cond}`, Component header normalization
- Tagged all releases v0.0.1–v0.0.12

## v0.0.11 (2026-07-28) — English + KB Sync + CLI Docs

- Entire repository translated to English: 130+ files (docs, agents, blocks, demo, guides)
- AGENTS.md: documented language rule (chat=DE, repo=EN)
- `dreego docs [--web] [--json] [--dump] [path]` — CLI doc browser
  - `--web`: open in browser with Codeberg-rendered markdown
  - `--json`: structured JSON output (headings, code blocks, links) for AI agents
  - `--dump`: concatenate all docs for LLM context
  - URL filtering: full Codeberg URLs in source → stripped in CLI, resolve correctly in browser
- `dreego feedback` — opens Codeberg issues page
- New docs: `_docs/runtime.md` (full API reference), expanded `_docs/getting-started.md` (components, layouts, dynamic routes)
- `_docs/` restructured with cross-links using full Codeberg URLs (browser-compatible)
- Knowledge base fully synced with v0.0.10 state
- Block created for future `dreego feedback` POST endpoint

## v0.0.10 (2026-07-28) — Static Assets

- `dreego/static/` folder: files are read during generate and registered inline
- MIME type via extension (.css, .js, .svg, .png, .ico, .html, .json, .woff2)
- Collision check: when static path collides with route → `dreego generate` error
- 3 static tests: basic, subdir, collision
- 71 integration tests total
- VS Code Extension v0.0.4: `<@Component>`-Tags, `import`, Filter `{var|raw}`, `$loop` highlighting

## v0.0.9 (2026-07-28) — Template Primitives

- `{#verbatim}` Block: raw output for JS templates
- `{#each}` with `$loop` variable: `$loop.Index`, `$loop.First`, `$loop.Last`, `$loop.Even`, `$loop.Odd`
- Template filters: `{var|raw}` (no escaping), `{var|upper}` (uppercase). Pipe syntax.
- `{#else}` in `{#if}` block: `{#if cond}...{#else}...{/if}`
- `{#each else}`: `{#each items as item}...{#each else}...{/each}` — empty list fallback
- Fix: `<header>`, `<main>`, `<footer>` prefix-match bug in scanTag
- Test system: `PASS/FAIL <path>`, `DREEGO_FILTER=<pattern>`, Docker build logs suppressed

## v0.0.8 (2026-07-28) — Named Slots

- Named Slots: `{#slot header}...{/slot}` block syntax in Components + Routes
- Component: `{#slot header}{/slot}` — placeholder renders `c.Get("slot_header")`
- Route: `{#slot header}<content>{/slot}` — defines content for named slot
- Default-Slot `{#slot}` stays without `{/slot}` (no change)
- 4 positive tests + 2 negative tests

## v0.0.7 (2026-07-28) — Test Coverage

- 41+ integration tests (up from 36): edge cases, negative tests, bugs
- `_docs/testing.md`: complete test strategy with 60+ test ideas
- Prop expressions in Components: `<@Card title={expr}/>`
- Nested Components: `<@Outer>` calls `<@Inner>`
- Session: `DelSessionVal`, `DestroySession`, no-store
- CSRF: `SetCSRF(false)` + disable test
- `--check` uses timestamp comparison (not git diff)
- All tests write .dreego files inline (no Docker COPY fixtures anymore)

## v0.0.6 (2026-07-28) — Component Completion

- Children slot passing: `<@Card>content</@Card>` → `{#slot}` in component works
- `dreego generate --check`: CI validation — exit non-zero when generated files are stale
- Named Slots: `{#slot header}` lexer/parser prepared (v0.0.7)

## v0.0.5 (2026-07-27) — Component Model

- Component system: `Component Name (props)` in `dreego/components/`, call via `<@Name>`
- Self-closing (`<@Icon name="star"/>`) and with children (`<@Card>...</@Card>`)
- Default slot via `{#slot}` in component template
- Scoped styles per component (`data-scope`)
- File-based discovery: `dreego/components/Card.dreego` → `<@Card>`
- 6 component integration tests + 2 bug tests
- `import "dreego/components/Name"` in route files (ParseHeader before Lex)
- Multi-file directory import: `import "dreego/components/button"` → `<@Login/>`

## v0.0.4 (2026-07-27) — Blueprints & Tests

- `dreego init <path>` — scaffold new project from embedded blueprint
- Blueprints via `//go:embed` in CLI binary, no external files needed
- Integration tests in `_tests/` via Docker containers (`make test`)

## v0.0.3 (2026-07-27) — Security & Developer Experience

- Session integration: `session.Store` interface + `CookieStore` (HMAC-signed) hooked into runtime
- Session middleware: context-based store injection per request
- SSRContext: `SessionVal`/`SetSessionVal`/`DelSessionVal`/`DestroySession` with secure defaults (`HttpOnly`, `Secure` TLS-aware)
- CSRF protection: double-submit cookie (Core-Conditional, default on) — Token via X-CSRF-Token header or csrf_token form field
- SSRContext: `CSRFToken()` for template rendering (hidden field)
- VS Code Extension: syntax highlighting + raccoon icon for `.dreego` files (`make dx`)
- **Breaking:** `pkg/` → `core/` (single package), single import `import "codeberg.org/dreego/dreego/core"`
- Plugins in separate repos (see `_docs/plugins.md`)

## v0.0.2 (2026-07-25) — Safety & Structure

- Route segments: `[id]` (square brackets) as convention for dynamic segments, compatible with Next.js/SvelteKit/Astro
- Route groups: `(group)/` — directories that do not appear in the URL (layout/middleware grouping)
- Flat gen package: all route handlers in `gen/routes.go` (no more `_ "import"`), solves Go import path problem with special characters
- Context refactoring: `map[string]string` → Interface + Embedding (`Context` interface + `SSRContext` struct)
- Recovery middleware: Panic → 500 with stack trace logging via slog
- XSS protection: auto-escaping of all `{variable}` expressions via `html.EscapeString`
- Custom error pages: `404.dreego` + `500.dreego`

## v0.0.1 (2026-07-25) — The Prototype

First prototype. Transpiler, Routing, Layout, Middleware, CLI.

### Features

- Formal transpiler pipeline: Lexer → Parser → AST → CodeGen
- All 5 sections: `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- Template logic: `{var}`, `{#if}`, `{#each}`, `{#slot}`, `{#head}`
- File-based routing: `dreego/routes/*.dreego`
- Dynamic segments: `[id]`, `[...catchall]`, `[[optional]]`, `(group)/`
- Layout system: `dreego/layouts/default.dreego` with `{#slot}` + `{#head}`
- CSS scoping: `data-scope` via source hash (12 characters)
- Central `dreego/gen/dree.go` for route imports
- `dreego/config.json`: redirects, rewrites, logging config
- RequestLogging middleware (Core-Conditional, JSONL format, IP capture)
- Redirect/Rewrite middleware
- CLI: `dreego generate [--force]`, `dreego build`, `dreego run [-d] [-t N]`
- Working demo server with net/http 1.22+

### Decisions

- [Error Handling](.agents/decisions/error-handling.md) — Typed errors, Recovery, Dev/Prod
- [Routing & Components](.agents/decisions/routing-and-components.md) — Hybrid routing, Plugin routes
- [Middleware System](.agents/decisions/middleware-system.md) — Core-Fixed vs Core-Conditional vs Plugin
- GLM Review: Per-directory `dree.go` output strategy confirmed
