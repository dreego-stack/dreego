
## v0.0.26 (2026-08-08) — Post-Release Fixes + SSE Plugin

### Added

- **SSE example plugin:** `plugins/sse` is a full v1 `core.Plugin` implementation (SSE hub, `/sse` route with `text/event-stream`, heartbeat, broadcast, embedded assets, docs with `dreego-sitemap.xml` seed).
- **`--version` / `-v` flags:** the CLI now accepts `--version` and `-v` alongside the `version` subcommand, matching the help/--help/-h trio. All three forms print the identical version (`_tests/core/CLI/version-flag/`).

### Fixed

- **Control flow in attribute values:** `{#if}`/`{#each}` inside quoted HTML attribute values produced silent corrupt codegen (route path: cond never referenced; component path: invalid Go). Detect and reject at parse time with a clear error suggesting wrapping the whole tag in `{#if}` instead (`_tests/core/Bugs/attr-if-in-attribute/`).
- **`<...>` dropped from go-section strings:** `parseGoSection` silently dropped tag tokens, so Go strings like `"TO: <HASH>"` or backtick SVGs lost their `<...>` content (silent data corruption). Tags are now reconstructed and a `SelfClose` flag keeps self-closing tags' trailing slash (`_tests/core/Bugs/go-string-lt/`).
- **Sections after leading template text:** a file starting with non-section text (e.g. `<!doctype html>`) swallowed following `<go>`/`<head>` blocks as template text. `parsePlainTemplate` now stops at section tags so sections work after leading text (`_tests/core/Bugs/text-before-section/`).
- **Scope div before doctype in error pages:** `GenerateErrorHandler` emitted `<div data-scope>` unconditionally, so a 404/500 template starting with `<!doctype html>` rendered the scope div first and fell into quirks mode. The scope wrapper is now suppressed when the first template node starts with `<!` (`_tests/core/Bugs/error-page-doctype/`, `_tests/core/Bugs/error-page-layout/`).
- **Title/meta-description dedupe in `{#head}` merge:** the merge appended layout head and route head, so pages rendered two `<title>` tags and the browser used the landing title. When the route head defines a `<title>` or `<meta name="description">`, the layout's copy is dropped from the merge (route wins); non-description meta/link tags are preserved and the layout title is kept when the route defines none (`_tests/core/Bugs/head-title-dedupe/`).
- **Component prop defaults:** `parseProps` read `= default` but `GenerateComponent` ignored it, so `<@Card/>` without a prop failed to compile. A variadic wrapper + empty-string fallback is generated for string props with defaults; bool/int defaults remain unsupported so explicit `false`/`0` is never overwritten (`_tests/core/Bugs/component-prop-default/`).
- **JSON/XML status committed before encode:** `WriteHeader(status)` ran before `Encode`, so an encode error returned a pre-committed 200. The encoded payload is now buffered and `WriteHeader` runs only after a successful encode; an encode error returns 500 (`core/response.go`).
- **Streaming through middleware:** `responseWriter` and `gzipResponseWriter` now implement `http.Flusher` so streaming handlers work behind `RequestLogging` and `Compress` (browsers send `Accept-Encoding: gzip`).
- **`dreego init` gen import:** `dreego init` generated `import _ "gen"` but `generate` places the package in `dreego/gen/`, so a fresh project failed to build with "package gen is not in std". The placeholder is now replaced with the module name read from the target `go.mod` (`_tests/core/CLI/init-import/`).
- **`plugins/sample` core require drift:** the sample plugin required `core v0.0.23` while `VERSION` was `v0.0.25`, so codegen behavior drifted from the docs. The require is bumped to the current `VERSION` at each release.

### Changed

- **Import style:** core is imported as `dreego` everywhere (`github.com/dreego-stack/dreego/core` → `dreego` alias), for consistency across codegen, tests, and docs.

- Full suite: 185 passed, 0 failed

## v0.0.25 (2026-08-06) — Plugin Interface v1

- **plugin-interface.1:** The frozen v1 `core.Plugin` contract shipped (`Name`, `RegisterRoutes`, `Middlewares`, `Assets`, `OnStart`, `OnShutdown`). Plugins import core and register via `core.UsePlugin(p)`; core never imports a plugin. Lifecycle: `StartPlugins(ctx)` / `ShutdownPlugins(ctx)` call `OnStart`/`OnShutdown` on every registered plugin and propagate the first error. Compile-time interface-satisfaction check plus route/middleware/lifecycle tests in `core/plugin_test.go`.
- **BREAKING (behavior change):** `core.Register(method, pattern, handler)` is now **idempotent** — registering the same `method`+`pattern` again **replaces** the existing handler instead of appending a duplicate route. This lets a later-registered plugin (or a reload) override a route deterministically. Downstream callers that relied on duplicate-registration behavior must update. See `core/runtime.go` `Register`.
- **middleware-hooks.1:** A plugin's `Middlewares()` are appended to the runtime middleware chain in FIFO order — the first registered plugin runs first on request entry, then the next, then the handler. The stack is fixated on the first `Build()`; registering a plugin afterwards does not reorder it. Order/overlap coverage in `core/plugin_test.go` (`TestPluginMiddlewareFIFOOrder`, `TestPluginMiddlewareOrderFixatedOnFirstBuild`).
- **route-hooks.1:** Programmatic plugin route registration — a plugin calls `core.Register(...)` inside its `RegisterRoutes()` and all routes are reachable through `core.ServeMux()` alongside file-based routes. Multi-route and overlap-plugin tests in `core/route_hooks_test.go` (+ `core/route_hooks_test_helpers.go`).
- **docs-extensibility.1:** `dreego docs` now resolves plugin docs from the local `plugins/<name>/_docs/` directory with priority over the embedded/remote fallback (no HTTP call for local plugin docs). Coverage in `cmd/dreego/docs_test.go`.
- **docs-embed.1:** Offline embedded docs. `dreego docs` reads `_docs/`, `README.md`, and `CHANGELOG.md` from a `//go:embed` copy (`cmd/dreego/embedded/`) so it works without a network. `_scripts/sync-embedded-docs.sh` mirrors the repo docs into the embedded dir after doc changes. Coverage in `cmd/dreego/docs_embed_test.go`.
- **frontmatter.1:** `core.ParseFrontmatter(src) (map[string]string, body)` splits a leading YAML-like `---` delimited frontmatter block off a `.dreego` source and exposes its `key: value` pairs as typed metadata; `:` in values is preserved and list values (`tags: [a, b]`) are normalized to a comma-joined string. `core/frontmatter.go` + `core/frontmatter_test.go`.
- **dev-server.1:** New `dreego dev` command watches `.dreego` files (500 ms poll), regenerates + rebuilds on change, and gracefully restarts the server (SIGTERM + reap). Build errors do not kill the watcher. `cmd/dreego/dev.go` (+ `cmdDev` dispatch in `cmd/dreego/main.go`).
- **head-dedupe.1:** When a route head defines a `<title>` or `<meta name="description">`, the layout's corresponding tag is dropped from the merged `{#head}` output — the route wins. Non-overridden layout head content (e.g. `<meta charset>`, `<link>`) is preserved. `core/codegen_head_dedupe.go` + `core/codegen_head_dedupe_test.go`; integration test `_tests/core/Bugs/head-title-dedupe/`.
- Full suite: 164 passed, 0 failed

## v0.0.24 (2026-08-05)

- **scaffold-fix.1:** `dreego new` scaffold now includes `go.sum`, and the generated `.gitignore` ignores only `dreego/gen/` (not the source dirs). `cmd/dreego/version.go` reads the repo `VERSION` file as a fallback so local dev builds work with `dreego new` and `go mod tidy`.
- **layout-head.1:** Layouts apply (`{#slot}`/`{#head}`) and route `<head>` works with or without a layout — covered by `Bugs/layout-not-applied`, `Bugs/route-head-without-layout`, `Bugs/layout-route-head-merge`, `Layout/no-layout` + unit tests (`core/codegen_layout_test.go`, `core/codegen_head_test.go`).
- **scoped-css.2:** `scopeCSS` in `core/codegen_helpers.go` rewritten as a recursive brace-tracked parser so declarations between `{}` are preserved verbatim (`radial-gradient`, `calc()`, `rgb()`), `@media` inner selectors stay scoped, and `@keyframes` bodies survive unscoped. Regression tests: `Bugs/scoped-style-declarations-lost`, `Bugs/scoped-style-comma-parens`, `Bugs/scoped-style-keyframes`.
- **component-attr-props.1:** `{prop}` inside HTML attributes is substituted and escaped, both in component calls (`<@Link url="...">`) and in component bodies (`<a href="{url}">`). `core/codegen_component.go` (`compTextWithAttrs`, `genComponentCall` via `extractAttrValues`) + `core/codegen_helpers.go` (`attrVal` resolves `{expr}` inside quoted values). Tests: `Bugs/component-attr-prop-substitution`, `Components/prop-expression`, `Components/multi-props`, `Components/empty-props` + unit tests.
- **typed-forms.1:** `BindForm` now binds `int`, `bool`, and `[]string` fields (beyond `reflect.String`); new `core.RegisterRule(name, fn)` custom-validator API; `ValidateForm` uses `fmt.Sprint` so `min`/`max`/`required` work on bound numeric values. Tests: `core/validate_typed_test.go` + `FormActions/form-int-binding`, `FormActions/form-bool-binding`.
- **dreegotest.1:** New exported `dreegotest` helper package (`dreegotest/`, own Go module) with `Get(t, path)`, `PostForm(t, path, form)`, and `RenderComponent(t, fn, props...)` for route/component unit tests against `core.ServeMux()`.
- **golden-tests-core.1:** Golden-file assertions for generated `gen/routes.go`/`gen/components.go` output (`core/codegen_golden_test.go` + `core/testdata/golden/*.golden`), run with `-update` to refresh fixtures.
- **port-schema / test stability:** `_tests/test.sh` runner now assigns deterministic sequential ports from `DREEGO_PORT_BASE` (default 20000) and installs `curl` once before the test loop; all ~28 server-based tests read `${DREEGO_PORT:-...}` and write the port directly into `main.go` (no more `sed -i`, no random-port collisions, no apk database-lock races). `run-timer-sigterm/test.sh` gained a `DREEGO_BIN` fallback so it runs standalone.
- Full suite: 161 passed, 1 expected failure (`core/CLI/new-go-sum`) — the `core/v0.0.24` git tag does not exist yet; it is created by `_scripts/release.sh` at release time. Not a regression.

## v0.0.23 (2026-08-03) — Unreleased — Late v0.0.22 Fixes

- **runtime:** New exported `core.Reset()` helper clears the cached middleware/router stack (`builtHandler`) so tests and reload paths start from a clean runtime state.
- **security-session.2:** `core/session_keys.go` now derives signing and encryption keys from the store secret via HMAC-SHA256(secret, label) instead of raw SHA-256 concatenation.
- **security-session.3:** `core/session.go` propagates JSON marshaling and AES-GCM/encryption errors from `CookieStore.Set` instead of silently returning an empty cookie. `core/session_crypto.go` `encryptPayload` now returns `(ciphertext, error)` and accepts an `io.Reader` for nonce generation, enabling testable error paths.
- **coding-standards:** Maximum file line limit raised from 120 to 300 lines.
- Tests: unit tests in `core/runtime_test.go` and `core/session_encrypt_test.go`; integration test `_tests/core/Middleware/session-encrypt/`.
- Full suite: 147 passed, 0 failed

## v0.0.22 (2026-08-03) — Released — ServeMux Cache + CodeGen Error Propagation + Session Encryption

- **servemux-cache.1:** `core/runtime.go` now caches the built middleware/router stack. `core.Build()` and `core.Listen()` reuse `builtHandler` once constructed, avoiding repeated `http.NewServeMux` and middleware wrapping.
- **codegen-errors.1:** All `core/codegen*.go` template generators return `(string, error)` and propagate failures instead of silently returning empty strings. New `core/codegen_component.go` contains component template generation. Fixed the nested `{#if}` in `{#else}` branch bug for component templates: `genTemplateNodeComp` now detects an else-if chain vs. a true else branch and emits nested blocks correctly.
- **security-session.1:** Optional AES-256-GCM session encryption in `core/session.go`. Passing `&core.Options{Encrypt: true}` to `store.Set` encrypts the JSON payload before the HMAC signature (encrypt-then-MAC). `core/session_crypto.go` provides `encryptPayload`/`decryptPayload`; tampered or key-rotated cookies are rejected.
- Tests: unit tests in `core/codegen_template_test.go`, `core/runtime_test.go`, and `core/session_encrypt_test.go`; integration tests `_tests/core/Template/component-nested-if-else/` and `_tests/core/Middleware/session-encrypt/`.
- Full suite: 147 passed, 0 failed
- Released. Tags pushed: core/v0.0.22, cmd/dreego/v0.0.22, plugins/sample/v0.0.22.

## v0.0.23 (2026-08-03) — Nested Control Flow + Head Expression Resolution

- **Fix (feedback-intake A):** Nested `{#if}` blocks inside the `{#else}` branch of a route template are no longer silently dropped. `core/codegen_template.go` `NodeIf` codegen now distinguishes an else-if chain from a true else branch and emits the nested blocks instead of returning an empty string — previously `dreego generate` succeeded but produced an empty template (with follow-up `go build` errors like `declared and not used`).
- **Fix (feedback-intake B):** Expressions in the `<head>` section of a route (e.g. `<title>{doc.Title}</title>`) are now resolved instead of being emitted raw. New `core/codegen_head.go` (`genHead`) splits head markup into literal and expression segments, applies escaping and the `raw`/`upper` filters; the four head emission sites in `core/codegen.go` (lines 137, 173, 187, 388) use it.
- Tests: unit test `TestGenTemplateNodeNestedIfInElseNotDropped` (`core/codegen_template_test.go`) + regression tests `_tests/core/Bugs/nested-if-in-else/` and `_tests/core/Bugs/head-expression-raw/`; existing `_tests/core/Bugs/head-expression/` extended. The component-path variant is covered by `_tests/core/Template/component-nested-if-else/` (v0.0.22).
- Full suite: 144 passed, 0 failed

## v0.0.21 (2026-08-03) — Single-Source Versioning + go install Fix

- **Fix:** `go install github.com/dreego-stack/dreego/cmd/dreego@latest` now works. Removed the relative `replace` directive from `cmd/dreego/go.mod` and `plugins/sample/go.mod` (relative replaces are invalid for non-main modules), replaced with a real published `require github.com/dreego-stack/dreego/core v0.0.22`. Local development still resolves `core` via `go.work` (`use ./core`).
- **Versioning:** New single source of truth `VERSION` file at repo root (`v0.0.22`). The CLI version derives from it at build time (`-ldflags -X main.version`) or, when installed via `go install pkg@tag`, from the module build info.
- New `dreego version` command prints the CLI version.
- `dreego new` now requires the CLI's own `core` version instead of a hardcoded one.
- New `_scripts/release.sh` creates `core/<V>`, `cmd/dreego/<V>`, `plugins/sample/<V>` directory-prefix tags from the `VERSION` file.
- Full suite: expected pass.

## v0.0.20 (2026-07-31) — Security Hardening

- Official plugins moved from separate repos into `plugins/` in this repository (one repo, many modules)
- New `plugins/sample/` minimal example plugin with its own `go.mod` importing `github.com/dreego-stack/dreego/core`
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

- **BREAKING**: `dreego-core/` → `core/` — import path changes from `github.com/dreego-stack/dreego/dreego-core` to `github.com/dreego-stack/dreego/core`
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
- **Breaking:** `pkg/` → `core/` (single package), single import `import "github.com/dreego-stack/dreego/core"`
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
