# Changelog

## v0.0.22 (2026-08-03) — ServeMux Cache + CodeGen Error Propagation + Session Encryption

- **servemux-cache.1:** `core/runtime.go` now caches the built middleware/router stack. `core.Build()` and `core.Listen()` reuse `builtHandler` once constructed, avoiding repeated `http.NewServeMux` and middleware wrapping. A `Reset()` helper clears the cache for tests.
- **codegen-errors.1:** All `core/codegen*.go` template generators return `(string, error)` and propagate failures instead of silently returning empty strings. New `core/codegen_component.go` contains component template generation. Fixed the nested `{#if}` in `{#else}` branch bug for component templates: `genTemplateNodeComp` now detects an else-if chain vs. a true else branch and emits nested blocks correctly.
- **security-session.1:** Optional AES-256-GCM session encryption in `core/session.go`. Passing `&core.Options{Encrypt: true}` to `c.SetSessionVal` (or `store.Set`) encrypts the JSON payload before the HMAC signature. `core/session_crypto.go` provides `encryptPayload`/`decryptPayload`; `core/session_keys.go` derives separate signing and encryption keys from the store secret. Tampered or key-rotated cookies are rejected.
- Tests: unit tests in `core/codegen_template_test.go` and `core/session_encrypt_test.go`; integration tests `_tests/core/Bugs/component-nested-if-else/` and `_tests/core/Middleware/session-encrypt/`.
- Full suite: 144 passed, 0 failed

## v0.0.23 (2026-08-03) — Nested Control Flow + Head Expression Resolution

- **Fix (feedback-intake A):** Nested `{#if}` blocks inside the `{#else}` branch of a route template are no longer silently dropped. `core/codegen_template.go` `NodeIf` codegen now distinguishes an else-if chain from a true else branch and emits the nested blocks instead of returning an empty string — previously `dreego generate` succeeded but produced an empty template (with follow-up `go build` errors like `declared and not used`).
- **Fix (feedback-intake B):** Expressions in the `<head>` section of a route (e.g. `<title>{doc.Title}</title>`) are now resolved instead of being emitted raw. New `core/codegen_head.go` (`genHead`) splits head markup into literal and expression segments, applies escaping and the `raw`/`upper` filters; the four head emission sites in `core/codegen.go` (lines 137, 173, 187, 388) use it.
- Tests: unit test `TestGenTemplateNodeNestedIfInElseNotDropped` (`core/codegen_template_test.go`) + regression tests `_tests/core/Bugs/nested-if-in-else/` and `_tests/core/Bugs/head-expression-raw/`; existing `_tests/core/Bugs/head-expression/` extended.
- Full suite: 144 passed, 0 failed

## v0.0.21 (2026-08-03) — Single-Source Versioning + go install Fix

- **Fix:** `go install codeberg.org/dreego/dreego/cmd/dreego@latest` now works. Removed the relative `replace` directive from `cmd/dreego/go.mod` and `plugins/sample/go.mod` (relative replaces are invalid for non-main modules), replaced with a real published `require codeberg.org/dreego/dreego/core v0.0.22`. Local development still resolves `core` via `go.work` (`use ./core`).
- **Versioning:** New single source of truth `VERSION` file at repo root (`v0.0.22`). The CLI version derives from it at build time (`-ldflags -X main.version`) or, when installed via `go install pkg@tag`, from the module build info.
- New `dreego version` command prints the CLI version.
- `dreego new` now requires the CLI's own `core` version instead of a hardcoded one.
- New `_scripts/release.sh` creates `core/<V>`, `cmd/dreego/<V>`, `plugins/sample/<V>` directory-prefix tags from the `VERSION` file.
- Full suite: expected pass.

## v0.0.20 (2026-07-31) — Security Hardening

- Official plugins moved from separate repos into `plugins/` in this repository (one repo, many modules)
- New `plugins/sample/` minimal example plugin with its own `go.mod` importing `codeberg.org/dreego/dreego/core`
- New `go.work` linking the root module and `plugins/sample` for local development
- Integration tests moved from `_tests/<Category>/` to `_tests/core/<Category>/`; `test.sh` runner now scans `_tests/core` and `_tests/plugins`
- All `test.sh` realrepo depth updated from `../../..` to `../../../..` (4 levels up from `_tests/core/<Group>/<name>/`)
- `_docs/plugins.md`, `_docs/plugin-interfaces.md`, `AGENTS.md` updated to describe the monorepo plugin model; Core never imports a plugin package
- Plugins with external dependencies get their own `go.mod`; dependency-free plugins can be plain packages
- Full suite: 141 passed, 0 failed

## v0.0.20 (2026-07-31) — Security Hardening

- `Content-Security-Policy` header now set by `SecurityHeaders` middleware with a permissive default (`self` + `unsafe-inline` for scripts/styles, common CDN/font sources) to support HTMX/Alpine.js and scoped CSS
- `core.SetCSP(value)` — override the policy from `main.go` for stricter/looser setups
- `csrf_token` readable cookie now sets `Secure` when the request is over TLS; `HttpOnly: false` and `SameSite=Strict` retained so the token stays JS-accessible yet transport-protected
- Tests: `core/middleware_csrf_test.go`, `core/middleware_security_test.go`, `core/session_secure_test.go`
- Integration tests: `csp-runtime`, `csp-override`, `csrf-cookie-samesite`
- Full suite: 141 passed, 0 failed

## v0.0.19 (2026-07-30) — Bug Fixes

- Fix B1: `{#if}` and `{#each}` now transpile correctly inside components **and** route templates
- Fix B2: `BindForm` returns an error instead of panicking on non-string fields
- Fix B3: `scopeCSS` preserves nested CSS blocks (e.g. `@media` queries) and scopes their inner selectors
- Fix B4: `hasValidateTag`/`hasFormTag` now only match tags inside the target struct body
- Fix B5: `splitGoSections` skips leading comments before deciding if a <go> block is a declaration
- Fix B6: `findMain` now matches `cmd/main.go` in addition to `demo/main.go`
- Fix B8: landing blueprint `config.json` uses `{"logging": {"enabled": true}}` instead of boolean
- Fix B11: `SetReady`/`readyHandler` use `atomic.Bool` to eliminate data race
- Fix B12: `generateCSRFToken` panics on `crypto/rand.Read` errors instead of ignoring them
- Fix B13: `findLayout` and `scanComponents` now propagate read/lex/parse/generate errors instead of swallowing them
- Fix B14: `GenerateComponent` now emits all `<go>` sections instead of only `Go[0]`
- Fix B15: `cleanSegment` and `patternSegment` now strip all wrapping bracket/underscore pairs from optional/dynamic segments
- Fix B16: `extractAttrValues` no longer splits on spaces inside brace expressions
- Fix B17: `atoi` now returns an error for non-digit input; `min=`/`max=` validation reports invalid numbers
- Fix B18: `unindent` handles both tabs and spaces; `splitGoSections` passes raw go code to unindent
- Fix B19: `findFormStruct` regex now handles complex parameter types and named returns
- Fix B20: `dreego run -t` sends SIGTERM instead of SIGKILL for graceful shutdown
- Fix B21: `RequestID` panics on `crypto/rand.Read` errors instead of ignoring them
- Lexer: `{` treated as template control-flow everywhere except inside `<go>`, `<head>`, `<script>`, and `<style>` sections
- Lexer: arbitrary HTML tags (e.g. `<ul>`, `<input>`) tokenize without mandatory balancing
- Codegen: `NodeIf`/`NodeEach` cases added to component template generation path
- Parser: unknown root tags and arbitrary open/close tags allowed in template contexts
- Test: `_tests/Bugs/component-if-each` covers `{#if}`/`{#each}` in components
- Test: `_tests/Bugs/bindform-non-string` covers B2
- Test: `_tests/Bugs/scoped-css-media` covers B3
- Test: `_tests/Bugs/form-tag-struct-name` covers B4
- Test: `_tests/Bugs/splitgo-comment-prefix` covers B5
- Test: `_tests/Bugs/findmain-cmd-dir` covers B6
- Test: `_tests/Bugs/landing-config-type` covers B8
- Test: `_tests/Bugs/component-multi-go` covers B14
- Test: `_tests/Bugs/clean-segment-optional` covers B15
- Test: `_tests/Bugs/component-attr-space` covers B16
- Test: `_tests/Bugs/validate-atoi-non-digit` covers B17
- Test: `_tests/Bugs/unindent-spaces` covers B18
- Test: `_tests/Bugs/form-handler-named-return` covers B19
- Test: `_tests/Bugs/run-timer-sigterm` covers B20
- Test: `_tests/Template/each-loop` converted from expected-failure to expected-success for route-level `{#each}`

## v0.0.18 (2026-07-29) — Package Restructuring

- **BREAKING**: `dreego-core/` → `core/` — import path changes from `codeberg.org/dreego/dreego/dreego-core` to `codeberg.org/dreego/dreego/core`
- **BREAKING**: `dreego-plugin/` removed — plugins live in separate repos under `codeberg.org/dreego/<name>`
- `_docs/plugins.md`: plugin architecture overview, planned plugins, interface contracts
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
