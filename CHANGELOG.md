# Changelog

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
- **Breaking:** `pkg/` → `dreego-core/` (single package), single import `import core "codeberg.org/dreego/dreego/dreego-core"`
- `dreego-plugin/` for future plugins (Auth, Redis, DB, etc.)

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
