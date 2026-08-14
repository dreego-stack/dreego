
## v0.0.42 - 2026-08-14

- Breaking: replace all package-global runtime state with an explicit App instance (New, Register, Build, Handler, Listen)
- Breaking: remove the central Plugin interface in favor of package-owned Register(app, typedOptions) error functions
- Breaking: generated code emits gen.Register(app) instead of init() with dreego.Register
- Feat: multiple App instances are isolated — two apps can run concurrently with different routes, middleware, and sessions
- Bug: reject configuration after Build and reject duplicate or reserved routes
- Bug: generate compilable registration code for projects without routes
- Breaking: make dreegotest request clients own an explicit App handler
- Docs: migrate public examples from removed package-global APIs to App methods

## v0.0.41 - 2026-08-14

- Breaking: remove EventBus, Queue, KVStore, Storage contracts and implementations from core
- Docs: update plugin-interfaces.md; session.Store remains the only core infrastructure contract

## [Unreleased]

- Breaking: remove EventBus, Queue, KVStore, Storage from core (plugins will own these APIs)

## v0.0.40 - 2026-08-14

- Bug: catch-all directories generate invalid Go patterns
- Bug: double-bracket optional segments silently become required
- Bug: duplicate route definitions silently override instead of failing
- Bug: method attributes on go sections are ignored

## v0.0.39 - 2026-08-14

- Breaking: remove dead Context.Render(name string, data any) stub and ErrRender from core
- Docs: replace c.Render() references in forms and form-actions docs

## v0.0.38 - 2026-08-14

- Bug: lexer treats < in Go blocks, > in quoted attributes, and < > in script/style as tags
- Bug: || in template expressions is misparsed as filter pipeline
- Bug: unknown template filters are silently ignored instead of failing

- Chore: block minor/major version bumps in the v0.0.x phase (AGENTS.md, release-prep.py, CI, pr.md.example)

## v0.0.37 - 2026-08-14

- Docs: define the SSR-first roadmap and migrate open work to conflict-resistant todo item files
- Feat: require explicit root sections in `.dreego` files
- Feat: use `{{ expression }}` for escaped template output while preserving typed component props
- Breaking: remove the unused public frontmatter parser and reject frontmatter in `.dreego` files
- Docs: define flat and `+page.dreego` routes with method-specific template sections
- Docs: require named component props, lexical slots, and generation-time validation
- Docs: define strong generated types with explicit dynamic HTTP boundaries
- Docs: choose explicit App-bound plugin registration functions before v1
- Security: clarify that escaped URL attributes still require scheme validation
- Docs: distinguish released behavior from accepted v0.1 architecture

## v0.0.36 - 2026-08-12

- Feat: add core Storage interface (storage-interface.1) — `Storage` interface (Put/Get/Delete/List/URL) like `database/sql`, interface only, plugins implement (S3/R2/Local)

## v0.0.35 - 2026-08-12

- Feat: add core Queue interface (queue-interface.1) — `Queue` interface (Dispatch/DispatchAfter/DispatchBatch/Worker/Use) + `Job`/`JobHandler`/`JobMiddleware`, like `database/sql`, interface only, plugins implement (Redis/NATS/In-Memory)

## v0.0.34 - 2026-08-12

- Feat: add core KVStore interface (kv-store.1) — `KVStore` interface (Get/Set/Delete/Expire with TTL), like `database/sql`, interface only, plugins implement (Redis/Ristretto/In-Memory), distinct from Storage (blobs)
- Chore: TODO reorg — observability.1 + api-swagger.1 moved to new "Plugins (external repos)" section

## v0.0.33 - 2026-08-12

- Bug: fix flaky `TestCookieStoreEncryptValueNotPlaintext` — the byte-level substring check for `user_id`/`42` in the decoded cookie value was statistically unsound: the random AES-GCM ciphertext (base64-encoded) can coincidentally contain those bytes (~1.5% per run), causing intermittent `make test` failures. The test now asserts the encryption marker and checks the decrypted payload instead, making it deterministic (verified 300/300 green).

## v0.0.32 - 2026-08-12

- Feat: add typed in-memory event bus (event-bus.1) — generic `EventBus[T]` interface (Publish/Subscribe/Unsubscribe) + `NewInMemoryBus[T]()`, concurrency-safe, panic recovery, ctx-cancellation
- Chore: CI standard-header check — `_tests/go/standard_header_test.go` validates every `test.sh` under `_tests/core/` carries the standard header (`# Using standard: _tests/how-to-test-sh.md` + `# What:`), enforced via `make test` in CI
- Chore: AGENTS.md test rules — Feature Workflow + Bug workflow updated to `_tests/go/*_test.go` + `dreegotest`, `test.sh` reference dropped
- Chore: `dreegotest` parallel-safe — codegen runs via cached CLI subprocess instead of global `os.Chdir`, enabling `t.Parallel()` across `_tests/go`; shutdown tests poll port readiness instead of fixed sleeps; `test.sh` prints test output on failure; cross-compile sets `CGO_ENABLED=0`
- Chore: remove `dreego run -t <seconds>` timer flag — the last shell test (`_tests/core/CLI/run-timer`) and the flaky `bug_run_timer_sigterm` integration test are removed; graceful shutdown (B20) stays covered by `TestDeploymentGracefulShutdown`; `standard_header_test` tolerates an empty/missing `_tests/core`

## v0.0.31 - 2026-08-10

- Chore: migrate remaining ~120 shell tests (`_tests/core/{Components,Template,Imports,Static,Session,Routing,Layout,ContentType,FormActions,Middleware,Config,Deployment,CLI,Bugs}`) to Go in `_tests/go/` via `dreegotest`
- Feat: extend `dreegotest` with CLI/project helpers (`CLIBin`, `ProjectDir`, `RunCLI`, `BuildInDirOK`, `LatestTag`) and a cookie-jar HTTP client (`ServeSetup`, `Client.Request/Cookie`)

- Chore: replace _todo/Blockwebchain system with linear TODO.md + TODO-Future.md
- Chore: lock "plugin" naming (addon → plugin), rename concept docs

## v0.0.30 - 2026-08-10

- Feat: `dreego docs` reads local module docs via `go.mod` (module itself → `vendor/` → module cache) instead of an embedded copy or HTTP — no network, single source of truth per module
- Feat: `dreego docs -p <plugin>` reads a plugin's docs from `github.com/dreego-stack/<plugin>` (requires it in `go.mod`)
- Feat: `dreego docs --list` lists every core + plugin page from each module's `_docs/sitemap.json`
- Feat: per-module `_docs/sitemap.json` defines each module's doc pages (replaces the embedded mirror)
- Fix: remove `cli/dreego/embedded/`, `cli/dreego/embed.go`, `_scripts/sync-embedded-docs.sh` and the duplicated docs copy
- Fix: `run-timer-sigterm` flake — timer logic extracted into `scheduleStop` with deterministic unit tests (`TestScheduleStopSendsSIGTERM`, `TestScheduleStopFallsBackToKill`); integration test uses prebuilt CLI + 1s timer instead of double `go run` + 3s wait

- Bug: guard `process` job in `pull-request-merging` to only run on merged PRs (or pushes to main), preventing failures when a PR is closed without merging

## v0.0.29 - 2026-08-09

- Fix: move CLI from `cmd/dreego/` to `cli/dreego/` — `go install github.com/dreego-stack/dreego/cli/dreego@latest` now works and tracks the single root tag
- Fix: `go install github.com/dreego-stack/dreego/cmd/dreego@latest` was broken because the old `cmd/dreego` path was cached by the Go module proxy as a separate submodule (its v0.0.22–v0.0.26 tags declared a stale `codeberg.org/dreego/dreego` module path). The proxy cache is immutable, so that path can never resolve to the root module again — renaming to `cli/dreego/` gives a fresh, never-seen path that resolves via the single root tag.
- Fix: update all build scripts, Makefile, Dockerfiles, tests, and docs that referenced the old `cmd/dreego` path.

## v0.0.28 - 2026-08-09

- Feat: remove VERSION file — the latest git tag is now the single source of truth for the CLI version
- Feat: CLI version derives from git tag at build time (`-ldflags -X main.version=$(git describe --tags --abbrev=0)`) or from build info (`go install pkg@tag`)
- Fix: merging workflow creates the tag from the version computed by release-prep.py (no more VERSION file read)
- Fix: test environment injects the version via build arg (Dockerfile + Makefile test + test.sh), version-drift test compares against DREEGO_VERSION
- Chore: pull-request-check.yml fetches tags (fetch-depth: 0) so the version-drift test is meaningful

## v0.0.27 - 2026-08-09

- Chore: rewrite CHANGELOG.md to the new line-based standard (Feat/Fix/Chore prefix per line, no section titles)
- Fix: run-timer-sigterm test — 10x faster (timer 30s → 3s), uses go run per how-to standard, no more flaky port race

## v0.0.26 (2026-08-08) — Post-Release Fixes + SSE Plugin

- Feat: SSE example plugin — `plugins/sse` is a full v1 `core.Plugin` implementation (SSE hub, `/sse` route with `text/event-stream`, heartbeat, broadcast, embedded assets, docs with `dreego-sitemap.xml` seed).
- Feat: `--version` / `-v` flags — the CLI now accepts `--version` and `-v` alongside the `version` subcommand, matching the help/--help/-h trio. All three forms print the identical version (`_tests/core/CLI/version-flag/`).
- Fix: Control flow in attribute values — `{#if}`/`{#each}` inside quoted HTML attribute values produced silent corrupt codegen (route path: cond never referenced; component path: invalid Go). Detect and reject at parse time with a clear error suggesting wrapping the whole tag in `{#if}` instead (`_tests/core/Bugs/attr-if-in-attribute/`).
- Fix: `<...>` dropped from go-section strings — `parseGoSection` silently dropped tag tokens, so Go strings like `"TO: <HASH>"` or backtick SVGs lost their `<...>` content (silent data corruption). Tags are now reconstructed and a `SelfClose` flag keeps self-closing tags' trailing slash (`_tests/core/Bugs/go-string-lt/`).
- Fix: Sections after leading template text — a file starting with non-section text (e.g. `<!doctype html>`) swallowed following `<go>`/`<head>` blocks as template text. `parsePlainTemplate` now stops at section tags so sections work after leading text (`_tests/core/Bugs/text-before-section/`).
- Fix: Scope div before doctype in error pages — `GenerateErrorHandler` emitted `<div data-scope>` unconditionally, so a 404/500 template starting with `<!doctype html>` rendered the scope div first and fell into quirks mode. The scope wrapper is now suppressed when the first template node starts with `<!` (`_tests/core/Bugs/error-page-doctype/`, `_tests/core/Bugs/error-page-layout/`).
- Fix: Title/meta-description dedupe in `{#head}` merge — the merge appended layout head and route head, so pages rendered two `<title>` tags and the browser used the landing title. When the route head defines a `<title>` or `<meta name="description">`, the layout's copy is dropped from the merge (route wins); non-description meta/link tags are preserved and the layout title is kept when the route defines none (`_tests/core/Bugs/head-title-dedupe/`).
- Fix: Component prop defaults — `parseProps` read `= default` but `GenerateComponent` ignored it, so `<@Card/>` without a prop failed to compile. A variadic wrapper + empty-string fallback is generated for string props with defaults; bool/int defaults remain unsupported so explicit `false`/`0` is never overwritten (`_tests/core/Bugs/component-prop-default/`).
- Fix: JSON/XML status committed before encode — `WriteHeader(status)` ran before `Encode`, so an encode error returned a pre-committed 200. The encoded payload is now buffered and `WriteHeader` runs only after a successful encode; an encode error returns 500 (`core/response.go`).
- Fix: Streaming through middleware — `responseWriter` and `gzipResponseWriter` now implement `http.Flusher` so streaming handlers work behind `RequestLogging` and `Compress` (browsers send `Accept-Encoding: gzip`).
- Fix: `dreego init` gen import — `dreego init` generated `import _ "gen"` but `generate` places the package in `dreego/gen/`, so a fresh project failed to build with "package gen is not in std". The placeholder is now replaced with the module name read from the target `go.mod` (`_tests/core/CLI/init-import/`).
- Fix: `plugins/sample` core require drift — the sample plugin required `core v0.0.23` while `VERSION` was `v0.0.25`, so codegen behavior drifted from the docs. The require is bumped to the current `VERSION` at each release.
- Chore: Import style — core is imported as `dreego` everywhere (`github.com/dreego-stack/dreego/core` → `dreego` alias), for consistency across codegen, tests, and docs.
- Chore: Full suite: 185 passed, 0 failed

## v0.0.25 (2026-08-06) — Plugin Interface v1

- Feat: plugin-interface.1 — The frozen v1 `core.Plugin` contract shipped (`Name`, `RegisterRoutes`, `Middlewares`, `Assets`, `OnStart`, `OnShutdown`). Plugins import core and register via `core.UsePlugin(p)`; core never imports a plugin. Lifecycle: `StartPlugins(ctx)` / `ShutdownPlugins(ctx)` call `OnStart`/`OnShutdown` on every registered plugin and propagate the first error. Compile-time interface-satisfaction check plus route/middleware/lifecycle tests in `core/plugin_test.go`.
- Feat: BREAKING (behavior change) — `core.Register(method, pattern, handler)` is now **idempotent** — registering the same `method`+`pattern` again **replaces** the existing handler instead of appending a duplicate route. This lets a later-registered plugin (or a reload) override a route deterministically. Downstream callers that relied on duplicate-registration behavior must update. See `core/runtime.go` `Register`.
- Feat: middleware-hooks.1 — A plugin's `Middlewares()` are appended to the runtime middleware chain in FIFO order — the first registered plugin runs first on request entry, then the next, then the handler. The stack is fixated on the first `Build()`; registering a plugin afterwards does not reorder it. Order/overlap coverage in `core/plugin_test.go` (`TestPluginMiddlewareFIFOOrder`, `TestPluginMiddlewareOrderFixatedOnFirstBuild`).
- Feat: route-hooks.1 — Programmatic plugin route registration — a plugin calls `core.Register(...)` inside its `RegisterRoutes()` and all routes are reachable through `core.ServeMux()` alongside file-based routes. Multi-route and overlap-plugin tests in `core/route_hooks_test.go` (+ `core/route_hooks_test_helpers.go`).
- Feat: docs-extensibility.1 — `dreego docs` now resolves plugin docs from the local `plugins/<name>/_docs/` directory with priority over the embedded/remote fallback (no HTTP call for local plugin docs). Coverage in `cmd/dreego/docs_test.go`.
- Feat: docs-embed.1 — Offline embedded docs. `dreego docs` reads `_docs/`, `README.md`, and `CHANGELOG.md` from a `//go:embed` copy (`cmd/dreego/embedded/`) so it works without a network. `_scripts/sync-embedded-docs.sh` mirrors the repo docs into the embedded dir after doc changes. Coverage in `cmd/dreego/docs_embed_test.go`.
- Feat: frontmatter.1 — `core.ParseFrontmatter(src) (map[string]string, body)` splits a leading YAML-like `---` delimited frontmatter block off a `.dreego` source and exposes its `key: value` pairs as typed metadata; `:` in values is preserved and list values (`tags: [a, b]`) are normalized to a comma-joined string. `core/frontmatter.go` + `core/frontmatter_test.go`.
- Feat: dev-server.1 — New `dreego dev` command watches `.dreego` files (500 ms poll), regenerates + rebuilds on change, and gracefully restarts the server (SIGTERM + reap). Build errors do not kill the watcher. `cmd/dreego/dev.go` (+ `cmdDev` dispatch in `cmd/dreego/main.go`).
- Feat: head-dedupe.1 — When a route head defines a `<title>` or `<meta name="description">`, the layout's corresponding tag is dropped from the merged `{#head}` output — the route wins. Non-overridden layout head content (e.g. `<meta charset>`, `<link>`) is preserved. `core/codegen_head_dedupe.go` + `core/codegen_head_dedupe_test.go`; integration test `_tests/core/Bugs/head-title-dedupe/`.
- Chore: Full suite: 164 passed, 0 failed

## v0.0.24 (2026-08-05)

- Feat: scaffold-fix.1 — `dreego new` scaffold now includes `go.sum`, and the generated `.gitignore` ignores only `dreego/gen/` (not the source dirs). `cmd/dreego/version.go` reads the repo `VERSION` file as a fallback so local dev builds work with `dreego new` and `go mod tidy`.
- Feat: layout-head.1 — Layouts apply (`{#slot}`/`{#head}`) and route `<head>` works with or without a layout — covered by `Bugs/layout-not-applied`, `Bugs/route-head-without-layout`, `Bugs/layout-route-head-merge`, `Layout/no-layout` + unit tests (`core/codegen_layout_test.go`, `core/codegen_head_test.go`).
- Feat: scoped-css.2 — `scopeCSS` in `core/codegen_helpers.go` rewritten as a recursive brace-tracked parser so declarations between `{}` are preserved verbatim (`radial-gradient`, `calc()`, `rgb()`), `@media` inner selectors stay scoped, and `@keyframes` bodies survive unscoped. Regression tests: `Bugs/scoped-style-declarations-lost`, `Bugs/scoped-style-comma-parens`, `Bugs/scoped-style-keyframes`.
- Feat: component-attr-props.1 — `{prop}` inside HTML attributes is substituted and escaped, both in component calls (`<@Link url="...">`) and in component bodies (`<a href="{url}">`). `core/codegen_component.go` (`compTextWithAttrs`, `genComponentCall` via `extractAttrValues`) + `core/codegen_helpers.go` (`attrVal` resolves `{expr}` inside quoted values). Tests: `Bugs/component-attr-prop-substitution`, `Components/prop-expression`, `Components/multi-props`, `Components/empty-props` + unit tests.
- Feat: typed-forms.1 — `BindForm` now binds `int`, `bool`, and `[]string` fields (beyond `reflect.String`); new `core.RegisterRule(name, fn)` custom-validator API; `ValidateForm` uses `fmt.Sprint` so `min`/`max`/`required` work on bound numeric values. Tests: `core/validate_typed_test.go` + `FormActions/form-int-binding`, `FormActions/form-bool-binding`.
- Feat: dreegotest.1 — New exported `dreegotest` helper package (`dreegotest/`, own Go module) with `Get(t, path)`, `PostForm(t, path, form)`, and `RenderComponent(t, fn, props...)` for route/component unit tests against `core.ServeMux()`.
- Feat: golden-tests-core.1 — Golden-file assertions for generated `gen/routes.go`/`gen/components.go` output (`core/codegen_golden_test.go` + `core/testdata/golden/*.golden`), run with `-update` to refresh fixtures.
- Feat: port-schema / test stability — `_tests/test.sh` runner now assigns deterministic sequential ports from `DREEGO_PORT_BASE` (default 20000) and installs `curl` once before the test loop; all ~28 server-based tests read `${DREEGO_PORT:-...}` and write the port directly into `main.go` (no more `sed -i`, no random-port collisions, no apk database-lock races). `run-timer-sigterm/test.sh` gained a `DREEGO_BIN` fallback so it runs standalone.
- Chore: Full suite: 161 passed, 1 expected failure (`core/CLI/new-go-sum`) — the `core/v0.0.24` git tag does not exist yet; it is created by `_scripts/release.sh` at release time. Not a regression.

## v0.0.23 (2026-08-03) — Unreleased — Late v0.0.22 Fixes

- Feat: runtime — New exported `core.Reset()` helper clears the cached middleware/router stack (`builtHandler`) so tests and reload paths start from a clean runtime state.
- Fix: security-session.2 — `core/session_keys.go` now derives signing and encryption keys from the store secret via HMAC-SHA256(secret, label) instead of raw SHA-256 concatenation.
- Fix: security-session.3 — `core/session.go` propagates JSON marshaling and AES-GCM/encryption errors from `CookieStore.Set` instead of silently returning an empty cookie. `core/session_crypto.go` `encryptPayload` now returns `(ciphertext, error)` and accepts an `io.Reader` for nonce generation, enabling testable error paths.
- Chore: coding-standards — Maximum file line limit raised from 120 to 300 lines.
- Chore: Tests — unit tests in `core/runtime_test.go` and `core/session_encrypt_test.go`; integration test `_tests/core/Middleware/session-encrypt/`.
- Chore: Full suite: 147 passed, 0 failed

## v0.0.22 (2026-08-03) — Released — ServeMux Cache + CodeGen Error Propagation + Session Encryption

- Feat: servemux-cache.1 — `core/runtime.go` now caches the built middleware/router stack. `core.Build()` and `core.Listen()` reuse `builtHandler` once constructed, avoiding repeated `http.NewServeMux` and middleware wrapping.
- Feat: codegen-errors.1 — All `core/codegen*.go` template generators return `(string, error)` and propagate failures instead of silently returning empty strings. New `core/codegen_component.go` contains component template generation. Fixed the nested `{#if}` in `{#else}` branch bug for component templates: `genTemplateNodeComp` now detects an else-if chain vs. a true else branch and emits nested blocks correctly.
- Feat: security-session.1 — Optional AES-256-GCM session encryption in `core/session.go`. Passing `&core.Options{Encrypt: true}` to `store.Set` encrypts the JSON payload before the HMAC signature (encrypt-then-MAC). `core/session_crypto.go` provides `encryptPayload`/`decryptPayload`; tampered or key-rotated cookies are rejected.
- Chore: Tests — unit tests in `core/codegen_template_test.go`, `core/runtime_test.go`, and `core/session_encrypt_test.go`; integration tests `_tests/core/Template/component-nested-if-else/` and `_tests/core/Middleware/session-encrypt/`.
- Chore: Full suite: 147 passed, 0 failed
- Chore: Released. Tags pushed: core/v0.0.22, cmd/dreego/v0.0.22, plugins/sample/v0.0.22.

## v0.0.23 (2026-08-03) — Nested Control Flow + Head Expression Resolution

- Fix: Nested `{#if}` blocks inside the `{#else}` branch of a route template are no longer silently dropped. `core/codegen_template.go` `NodeIf` codegen now distinguishes an else-if chain from a true else branch and emits the nested blocks instead of returning an empty string — previously `dreego generate` succeeded but produced an empty template (with follow-up `go build` errors like `declared and not used`).
- Fix: Expressions in the `<head>` section of a route (e.g. `<title>{doc.Title}</title>`) are now resolved instead of being emitted raw. New `core/codegen_head.go` (`genHead`) splits head markup into literal and expression segments, applies escaping and the `raw`/`upper` filters; the four head emission sites in `core/codegen.go` (lines 137, 173, 187, 388) use it.
- Chore: Tests — unit test `TestGenTemplateNodeNestedIfInElseNotDropped` (`core/codegen_template_test.go`) + regression tests `_tests/core/Bugs/nested-if-in-else/` and `_tests/core/Bugs/head-expression-raw/`; existing `_tests/core/Bugs/head-expression/` extended. The component-path variant is covered by `_tests/core/Template/component-nested-if-else/` (v0.0.22).
- Chore: Full suite: 144 passed, 0 failed

## v0.0.21 (2026-08-03) — Single-Source Versioning + go install Fix

- Fix: `go install github.com/dreego-stack/dreego/cmd/dreego@latest` now works. Removed the relative `replace` directive from `cmd/dreego/go.mod` and `plugins/sample/go.mod` (relative replaces are invalid for non-main modules), replaced with a real published `require github.com/dreego-stack/dreego/core v0.0.22`. Local development still resolves `core` via `go.work` (`use ./core`).
- Feat: Versioning — New single source of truth `VERSION` file at repo root (`v0.0.22`). The CLI version derives from it at build time (`-ldflags -X main.version`) or, when installed via `go install pkg@tag`, from the module build info.
- Feat: New `dreego version` command prints the CLI version.
- Feat: `dreego new` now requires the CLI's own `core` version instead of a hardcoded one.
- Chore: New `_scripts/release.sh` creates `core/<V>`, `cmd/dreego/<V>`, `plugins/sample/<V>` directory-prefix tags from the `VERSION` file.
- Chore: Full suite: expected pass.

## v0.0.20 (2026-07-31) — Security Hardening

- Chore: Official plugins moved from separate repos into `plugins/` in this repository (one repo, many modules)
- Feat: New `plugins/sample/` minimal example plugin with its own `go.mod` importing `github.com/dreego-stack/dreego/core`
- Chore: New `go.work` linking the root module and `plugins/sample` for local development
- Chore: Integration tests moved from `_tests/<Category>/` to `_tests/core/<Category>/`; `test.sh` runner now scans `_tests/core` and `_tests/plugins`
- Chore: All `test.sh` realrepo depth updated from `../../..` to `../../../..` (4 levels up from `_tests/core/<Group>/<name>/`)
- Chore: `_docs/plugins.md`, `_docs/plugin-interfaces.md`, `AGENTS.md` updated to describe the monorepo plugin model; Core never imports a plugin package
- Chore: Plugins with external dependencies get their own `go.mod`; dependency-free plugins can be plain packages
- Chore: Full suite: 141 passed, 0 failed

## v0.0.20 (2026-07-31) — Security Hardening

- Feat: `Content-Security-Policy` header now set by `SecurityHeaders` middleware with a permissive default (`self` + `unsafe-inline` for scripts/styles, common CDN/font sources) to support HTMX/Alpine.js and scoped CSS
- Feat: `core.SetCSP(value)` — override the policy from `main.go` for stricter/looser setups
- Fix: `csrf_token` readable cookie now sets `Secure` when the request is over TLS; `HttpOnly: false` and `SameSite=Strict` retained so the token stays JS-accessible yet transport-protected
- Chore: Tests — `core/middleware_csrf_test.go`, `core/middleware_security_test.go`, `core/session_secure_test.go`
- Chore: Integration tests — `csp-runtime`, `csp-override`, `csrf-cookie-samesite`
- Chore: Full suite: 141 passed, 0 failed

## v0.0.19 (2026-07-30) — Bug Fixes

- Fix: B1 — `{#if}` and `{#each}` now transpile correctly inside components **and** route templates
- Fix: B2 — `BindForm` returns an error instead of panicking on non-string fields
- Fix: B3 — `scopeCSS` preserves nested CSS blocks (e.g. `@media` queries) and scopes their inner selectors
- Fix: B4 — `hasValidateTag`/`hasFormTag` now only match tags inside the target struct body
- Fix: B5 — `splitGoSections` skips leading comments before deciding if a <go> block is a declaration
- Fix: B6 — `findMain` now matches `cmd/main.go` in addition to `demo/main.go`
- Fix: B8 — landing blueprint `config.json` uses `{"logging": {"enabled": true}}` instead of boolean
- Fix: B11 — `SetReady`/`readyHandler` use `atomic.Bool` to eliminate data race
- Fix: B12 — `generateCSRFToken` panics on `crypto/rand.Read` errors instead of ignoring them
- Fix: B13 — `findLayout` and `scanComponents` now propagate read/lex/parse/generate errors instead of swallowing them
- Fix: B14 — `GenerateComponent` now emits all `<go>` sections instead of only `Go[0]`
- Fix: B15 — `cleanSegment` and `patternSegment` now strip all wrapping bracket/underscore pairs from optional/dynamic segments
- Fix: B16 — `extractAttrValues` no longer splits on spaces inside brace expressions
- Fix: B17 — `atoi` now returns an error for non-digit input; `min=`/`max=` validation reports invalid numbers
- Fix: B18 — `unindent` handles both tabs and spaces; `splitGoSections` passes raw go code to unindent
- Fix: B19 — `findFormStruct` regex now handles complex parameter types and named returns
- Fix: B20 — `dreego run -t` sends SIGTERM instead of SIGKILL for graceful shutdown
- Fix: B21 — `RequestID` panics on `crypto/rand.Read` errors instead of ignoring them
- Chore: Lexer — `{` treated as template control-flow everywhere except inside `<go>`, `<head>`, `<script>`, and `<style>` sections
- Chore: Lexer — arbitrary HTML tags (e.g. `<ul>`, `<input>`) tokenize without mandatory balancing
- Chore: Codegen — `NodeIf`/`NodeEach` cases added to component template generation path
- Chore: Parser — unknown root tags and arbitrary open/close tags allowed in template contexts
- Chore: Test — `_tests/Bugs/component-if-each` covers `{#if}`/`{#each}` in components
- Chore: Test — `_tests/Bugs/bindform-non-string` covers B2
- Chore: Test — `_tests/Bugs/scoped-css-media` covers B3
- Chore: Test — `_tests/Bugs/form-tag-struct-name` covers B4
- Chore: Test — `_tests/Bugs/splitgo-comment-prefix` covers B5
- Chore: Test — `_tests/Bugs/findmain-cmd-dir` covers B6
- Chore: Test — `_tests/Bugs/landing-config-type` covers B8
- Chore: Test — `_tests/Bugs/component-multi-go` covers B14
- Chore: Test — `_tests/Bugs/clean-segment-optional` covers B15
- Chore: Test — `_tests/Bugs/component-attr-space` covers B16
- Chore: Test — `_tests/Bugs/validate-atoi-non-digit` covers B17
- Chore: Test — `_tests/Bugs/unindent-spaces` covers B18
- Chore: Test — `_tests/Bugs/form-handler-named-return` covers B19
- Chore: Test — `_tests/Bugs/run-timer-sigterm` covers B20
- Chore: Test — `_tests/Template/each-loop` converted from expected-failure to expected-success for route-level `{#each}`

## v0.0.18 (2026-07-29) — Package Restructuring

- Chore: BREAKING — `dreego-core/` → `core/` — import path changes from `github.com/dreego-stack/dreego/dreego-core` to `github.com/dreego-stack/dreego/core`
- Chore: BREAKING — `dreego-plugin/` removed — plugins live in separate repos under `codeberg.org/dreego/<name>`
- Chore: `_docs/plugins.md`: plugin architecture overview, planned plugins, interface contracts
- Chore: AGENTS.md updated to reflect new directory structure

## v0.0.17 (2026-07-29) — Production Deployment + Request-ID

- Feat: Graceful Shutdown — `core.Listen()` uses `http.Server` with SIGINT/SIGTERM handling, 10s drain timeout
- Feat: Cross-Compile — `dreego build --target linux/amd64` sets GOOS/GOARCH for target platform
- Feat: Request-ID Middleware — `X-Request-ID` header on every request — client-supplied or auto-generated (16-char hex), injected into context + JSONL logs (`rid` field), accessible via `c.RequestID()`
- Feat: Production Dockerfile — `FROM scratch` — 3-stage build, `CGO_ENABLED=0`, static binary
- Chore: Hot Reload — `_docs/hot-reload.md` — Air config with `.air.toml` + `entr` alternative
- Chore: Rejected — hot-reload.1, live-reload.1, smart-recompile.1 — replaced by Air documentation
- Chore: Block — request-id.1 completed (chain 34)

## v0.0.16 (2026-07-29) — Form Actions

- Feat: `<form g-action="Login">` — declarative server-side form handling with auto-generated pipeline
- Feat: Generated POST handler — `r.ParseForm()` → struct mapping via `form:"email"` tags → validation via `validate:"required,email"` tags → handler call → redirect
- Feat: `c.Redirect(url, code)` — PRG pattern (Post-Redirect-Get) with `ErrRedirect` sentinel
- Feat: `c.Errors(field)` / `c.Old(field)` — validation error and old value access in templates
- Feat: `BindForm(r, target)` — maps form values to struct fields (explicit `form:` tag or lowercase field name)
- Feat: `ValidateForm(form)` — validates struct via `validate:` tags (required, email, min, max) — no external deps
- Feat: `SaveErrors(c, errs)` / `SaveOld(c, form)` — automatically stores validation state for template re-render
- Feat: Codegen — `splitGoSections` separates type/func declarations from inline code — form structs and handlers at package level
- Feat: Context interface extended — `SessionVal`, `SetSessionVal`, `DelSessionVal`, `CSRFToken`, `Redirect`
- Feat: `scanFormActions` detects `<form g-action>` in templates, wires matching handlers in POST dispatch
- Chore: 15 new tests — 11 parser/codegen + 4 runtime HTTP (valid submit, invalid re-render, CSRF, plain form)
- Chore: 112 integration tests total

## v0.0.15 (2026-07-28) — Content-Type Routing

- Feat: `<go type="json">` — JSON endpoints with `c.JSON()`, `c.Bind()`, auto Content-Type
- Feat: `<go type="xml">` — XML endpoints with `c.XML()`, auto Content-Type
- Feat: `<go type="custom">` — developer manages Content-Type + response via `c.Write()`
- Feat: Multiple `<go>` blocks merged — shared `<go>` runs always, typed blocks check `Accept` header
- Feat: Pure JSON/XML routes (no `<div>`) skip template rendering entirely
- Feat: `c.Write(status, contentType, body)` for arbitrary formats (FlatBuffers, Protobuf, etc.)
- Feat: `c.Wants(mime)` for manual content negotiation
- Feat: `core/response.go` — `JSON()`, `XML()`, `Bind()`, `Write()`, `Wants()`
- Chore: 6 new tests (87 total)

## v0.0.14 (2026-07-28) — Production Middleware

- Feat: `GET /health` → 200 `ok` — process liveness probe, always available
- Feat: `GET /ready` → 200 `ready` / 503 `not ready` — traffic readiness via `core.SetReady(bool)`
- Feat: Health endpoints registered before user routes (cannot be overridden)
- Feat: Security Headers middleware — `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` — core-fixed
- Feat: Gzip compression middleware — compresses responses when client accepts gzip (`Accept-Encoding`) — core-fixed
- Feat: Middleware chain — Recovery → SecurityHeaders → Compression → RequestLogging → Session → CSRF → Redirect/Rewrite → Router
- Chore: 3 new tests — health-checks, security-headers, compression
- Fix: `layout.Head.Content` now rendered in codegen — `{#head}` inside `<head>` section processed correctly
- Fix: `genLayoutNode` handles both `{#head}` and `{#slot}` in same NodeText (was dropping `{#slot}` after `{#head}`)
- Chore: 3 new regression tests — layout-head-lost, layout-head-ok, head-dropped-by-layout (content verification)
- Chore: 81 integration tests total

## v0.0.13 (2026-07-28) — Scaffolding + Split-Gen

- Feat: `dreego new <name>` — creates landing page project with `go mod init` auto-setup
- Feat: Landing blueprint — Hero + FeatureCard components, layout with `{#head}` + `{#slot}`, Tailwind CDN
- Feat: Split-Gen — `gen/routes.go` + `gen/components.go` + `gen/dree.go` (config + static) — file-level caching
- Feat: `isUpToDate()` — files only written when content changes
- Fix: `.gitignore` fix — `/dreego` → only root binary; `dreego/gen/` covers all generated files
- Feat: `dreego generate` timing — ns → ms display
- Chore: 5 new tests — head-without-layout, head-dropped-by-layout, name-clash, nested-routes, fmt
- Chore: 74 integration tests total
- Chore: Tags v0.0.1–v0.0.13 on Codeberg

## v0.0.12 (2026-07-28) — dreego fmt

- Feat: `dreego fmt` — formats `.dreego` files in-place (section ordering, template expressions, control flow)
- Feat: `dreego fmt --check` — CI mode: exit non-zero if files need formatting
- Feat: `dreego fmt --stdout` — print formatted output to stdout
- Feat: `core.Format()` — reusable formatting function
- Feat: Formats — `{ var }` → `{var}`, `{#if  cond}` → `{#if cond}`, Component header normalization
- Chore: Tagged all releases v0.0.1–v0.0.12

## v0.0.11 (2026-07-28) — English + KB Sync + CLI Docs

- Chore: Entire repository translated to English — 130+ files (docs, agents, blocks, demo, guides)
- Chore: AGENTS.md — documented language rule (chat=DE, repo=EN)
- Feat: `dreego docs [--web] [--json] [--dump] [path]` — CLI doc browser
  - `--web`: open in browser with Codeberg-rendered markdown
  - `--json`: structured JSON output (headings, code blocks, links) for AI agents
  - `--dump`: concatenate all docs for LLM context
  - URL filtering: full Codeberg URLs in source → stripped in CLI, resolve correctly in browser
- Feat: `dreego feedback` — opens Codeberg issues page
- Feat: New docs — `_docs/runtime.md` (full API reference), expanded `_docs/getting-started.md` (components, layouts, dynamic routes)
- Chore: `_docs/` restructured with cross-links using full Codeberg URLs (browser-compatible)
- Chore: Knowledge base fully synced with v0.0.10 state
- Chore: Block created for future `dreego feedback` POST endpoint

## v0.0.10 (2026-07-28) — Static Assets

- Feat: `dreego/static/` folder — files are read during generate and registered inline
- Feat: MIME type via extension (.css, .js, .svg, .png, .ico, .html, .json, .woff2)
- Feat: Collision check — when static path collides with route → `dreego generate` error
- Chore: 3 static tests — basic, subdir, collision
- Chore: 71 integration tests total
- Chore: VS Code Extension v0.0.4 — `<@Component>`-Tags, `import`, Filter `{var|raw}`, `$loop` highlighting

## v0.0.9 (2026-07-28) — Template Primitives

- Feat: `{#verbatim}` Block — raw output for JS templates
- Feat: `{#each}` with `$loop` variable — `$loop.Index`, `$loop.First`, `$loop.Last`, `$loop.Even`, `$loop.Odd`
- Feat: Template filters — `{var|raw}` (no escaping), `{var|upper}` (uppercase). Pipe syntax.
- Feat: `{#else}` in `{#if}` block — `{#if cond}...{#else}...{/if}`
- Feat: `{#each else}` — `{#each items as item}...{#each else}...{/each}` — empty list fallback
- Fix: `<header>`, `<main>`, `<footer>` prefix-match bug in scanTag
- Chore: Test system — `PASS/FAIL <path>`, `DREEGO_FILTER=<pattern>`, Docker build logs suppressed

## v0.0.8 (2026-07-28) — Named Slots

- Feat: Named Slots — `{#slot header}...{/slot}` block syntax in Components + Routes
- Feat: Component — `{#slot header}{/slot}` — placeholder renders `c.Get("slot_header")`
- Feat: Route — `{#slot header}<content>{/slot}` — defines content for named slot
- Feat: Default-Slot `{#slot}` stays without `{/slot}` (no change)
- Chore: 4 positive tests + 2 negative tests

## v0.0.7 (2026-07-28) — Test Coverage

- Chore: 41+ integration tests (up from 36) — edge cases, negative tests, bugs
- Chore: `_docs/testing.md` — complete test strategy with 60+ test ideas
- Feat: Prop expressions in Components — `<@Card title={expr}/>`
- Feat: Nested Components — `<@Outer>` calls `<@Inner>`
- Feat: Session — `DelSessionVal`, `DestroySession`, no-store
- Feat: CSRF — `SetCSRF(false)` + disable test
- Chore: `--check` uses timestamp comparison (not git diff)
- Chore: All tests write .dreego files inline (no Docker COPY fixtures anymore)

## v0.0.6 (2026-07-28) — Component Completion

- Feat: Children slot passing — `<@Card>content</@Card>` → `{#slot}` in component works
- Feat: `dreego generate --check` — CI validation — exit non-zero when generated files are stale
- Chore: Named Slots — `{#slot header}` lexer/parser prepared (v0.0.7)

## v0.0.5 (2026-07-27) — Component Model

- Feat: Component system — `Component Name (props)` in `dreego/components/`, call via `<@Name>`
- Feat: Self-closing (`<@Icon name="star"/>`) and with children (`<@Card>...</@Card>`)
- Feat: Default slot via `{#slot}` in component template
- Feat: Scoped styles per component (`data-scope`)
- Feat: File-based discovery — `dreego/components/Card.dreego` → `<@Card>`
- Chore: 6 component integration tests + 2 bug tests
- Feat: `import "dreego/components/Name"` in route files (ParseHeader before Lex)
- Feat: Multi-file directory import — `import "dreego/components/button"` → `<@Login/>`

## v0.0.4 (2026-07-27) — Blueprints & Tests

- Feat: `dreego init <path>` — scaffold new project from embedded blueprint
- Feat: Blueprints via `//go:embed` in CLI binary, no external files needed
- Chore: Integration tests in `_tests/` via Docker containers (`make test`)

## v0.0.3 (2026-07-27) — Security & Developer Experience

- Feat: Session integration — `session.Store` interface + `CookieStore` (HMAC-signed) hooked into runtime
- Feat: Session middleware — context-based store injection per request
- Feat: SSRContext — `SessionVal`/`SetSessionVal`/`DelSessionVal`/`DestroySession` with secure defaults (`HttpOnly`, `Secure` TLS-aware)
- Feat: CSRF protection — double-submit cookie (Core-Conditional, default on) — Token via X-CSRF-Token header or csrf_token form field
- Feat: SSRContext — `CSRFToken()` for template rendering (hidden field)
- Chore: VS Code Extension — syntax highlighting + raccoon icon for `.dreego` files (`make dx`)
- Chore: Breaking — `pkg/` → `core/` (single package), single import `import "github.com/dreego-stack/dreego/core"`
- Chore: Plugins in separate repos (see `_docs/plugins.md`)

## v0.0.2 (2026-07-25) — Safety & Structure

- Feat: Route segments — `[id]` (square brackets) as convention for dynamic segments, compatible with Next.js/SvelteKit/Astro
- Feat: Route groups — `(group)/` — directories that do not appear in the URL (layout/middleware grouping)
- Feat: Flat gen package — all route handlers in `gen/routes.go` (no more `_ "import"`), solves Go import path problem with special characters
- Feat: Context refactoring — `map[string]string` → Interface + Embedding (`Context` interface + `SSRContext` struct)
- Feat: Recovery middleware — Panic → 500 with stack trace logging via slog
- Feat: XSS protection — auto-escaping of all `{variable}` expressions via `html.EscapeString`
- Feat: Custom error pages — `404.dreego` + `500.dreego`

## v0.0.1 (2026-07-25) — The Prototype

First prototype. Transpiler, Routing, Layout, Middleware, CLI.

- Feat: Formal transpiler pipeline — Lexer → Parser → AST → CodeGen
- Feat: All 5 sections — `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- Feat: Template logic — `{var}`, `{#if}`, `{#each}`, `{#slot}`, `{#head}`
- Feat: File-based routing — `dreego/routes/*.dreego`
- Feat: Dynamic segments — `[id]`, `[...catchall]`, `[[optional]]`, `(group)/`
- Feat: Layout system — `dreego/layouts/default.dreego` with `{#slot}` + `{#head}`
- Feat: CSS scoping — `data-scope` via source hash (12 characters)
- Feat: Central `dreego/gen/dree.go` for route imports
- Feat: `dreego/config.json` — redirects, rewrites, logging config
- Feat: RequestLogging middleware (Core-Conditional, JSONL format, IP capture)
- Feat: Redirect/Rewrite middleware
- Feat: CLI — `dreego generate [--force]`, `dreego build`, `dreego run [-d] [-t N]`
- Feat: Working demo server with net/http 1.22+

- Chore: Decisions — [Error Handling](.agents/decisions/error-handling.md), [Routing & Components](.agents/decisions/routing-and-components.md), [Middleware System](.agents/decisions/middleware-system.md), GLM Review
