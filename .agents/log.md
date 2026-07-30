---
type: Log
title: Knowledge Base Changelog
description: Record of all changes to the .agents/ knowledge bundle
tags: [log]
timestamp: 2026-07-28T21:33:00Z
---

# log

## 2026-07-30 — B1: `{#if}`/`{#each}` transpilation in components and routes

- Fix B1: `core/codegen.go` `genTemplateNodeComp` now handles `NodeIf`/`NodeEach`
- Lexer: `{` treated as template control-flow outside `<go>`/`<head>`/`<script>`/`<style>`; arbitrary HTML tags tokenized without forced balancing
- Parser: arbitrary `TokenTagOpen`/`TokenTagClose` allowed inside `{#if}`/`{#each}`/`{#slot}` nodes; attributes preserved
- Test added: `_tests/Bugs/component-if-each` (B1 regression test)
- Test updated: `_tests/Template/each-loop` now expects route-level `{#each}` to transpile and build successfully
- Full suite: 124 passed, 0 failed

## 2026-07-30 — B2: `BindForm` non-string fields

- Fix B2: `core/validate.go` `BindForm` checks `field.Type.Kind()` before `SetString`
- Returns `fmt.Errorf("unsupported field type %s for field %s", ...)` instead of panicking
- Unit test: `TestBindFormNonStringFieldReturnsError` replaces `TestBindFormNonStringFieldSkipped`
- Integration test added: `_tests/Bugs/bindform-non-string`
- Full suite: 125 passed, 0 failed

## 2026-07-30 — B3: `scopeCSS` preserves `@media` queries

- Fix B3: `core/codegen_template.go` `scopeCSS` rewritten with brace-depth state machine
- Top-level selectors get `[data-scope=hash]` prefix; `@media` blocks keep their prefix on inner selectors
- Integration test added: `_tests/Bugs/scoped-css-media`
- Full suite: 125 passed, 0 failed

## 2026-07-30 — B4: `hasValidateTag`/`hasFormTag` respect `structName`

- Fix B4: `core/forms.go` `hasValidateTag`/`hasFormTag` now locate the target struct body via `type <structName> struct {`
- Only tags inside that struct's fields are matched, avoiding false positives from unrelated structs
- Added helper `hasTagInStruct` with brace-depth extraction
- Integration test added: `_tests/Bugs/form-tag-struct-name`
- Full suite: 127 passed, 0 failed

## 2026-07-30 — B5: `splitGoSections` comment prefix

- Fix B5: `core/codegen.go` `splitGoSections` skips leading `//` comments when determining the first non-empty line
- Declarations with doc comments are now correctly classified as package-level code
- Integration test added: `_tests/Bugs/splitgo-comment-prefix`
- Full suite: 128 passed, 0 failed

## 2026-07-30 — B6: `findMain` `cmd/` directory

- Fix B6: `cmd/dreego/main.go` `findMain` now matches `cmd/main.go` (removed erroneous `d != "cmd"` guard)
- Integration test added: `_tests/Bugs/findmain-cmd-dir`
- Full suite: 129 passed, 0 failed

## 2026-07-30 — B8: landing blueprint `config.json` logging type

- Fix B8: `cmd/dreego/blueprints/landing/dreego/config.json` changed `"logging": true` to `"logging": {"enabled": true}`
- Matches `core.Settings.Logging` struct shape
- Integration test added: `_tests/Bugs/landing-config-type`
- Full suite: 130 passed, 0 failed

## 2026-07-30 — B11: `SetReady` data race

- Fix B11: `core/middleware_health.go` `ready` variable replaced with `sync/atomic.Bool`
- `SetReady` uses `Store`, `readyHandler` uses `Load`
- Unit test added: `core/middleware_health_test.go` `TestSetReadyNoRace`
- Full suite: 130 passed, 0 failed

## 2026-07-30 — B12: CSRF `rand.Read` error ignored

- Fix B12: `core/middleware_csrf.go` `generateCSRFToken` now checks `crypto/rand.Read` error
- Panics with clear message on failure instead of silently using weak entropy
- Full suite: 130 passed, 0 failed

## 2026-07-30 — B13: `findLayout`/`scanComponents` error swallowing

- Fix B13: `core/generate.go` `findLayout` now returns `(*File, error)` and propagates read/lex/parse errors
- Fix B13: `core/generate.go` `scanComponents` now returns `(genDir, sources, error)` and propagates read/lex/parse/generate errors
- Callers in `Run()` updated to return early on component/layout errors
- Full suite: 130 passed, 0 failed

## 2026-07-30 — B14: `GenerateComponent` ignores `<go>` after `Go[0]`

- Fix B14: `core/codegen.go` `GenerateComponent` now iterates over all `file.Go` sections
- Previously only `file.Go[0].Code` was emitted into the component render function
- Integration test added: `_tests/Bugs/component-multi-go`
- Full suite: 131 passed, 0 failed

## 2026-07-30 — B15: `cleanSegment`/`patternSegment` nested brackets

- Fix B15: `core/generate.go` `cleanSegment` now strips all wrapping bracket/underscore pairs
- Fix B15: `core/generate.go` `patternSegment` only wraps segments that were actually bracketed/underscored
- Plain segments like `about` remain literal route parts; `[[opt]]` becomes `/{opt}`
- Integration test added: `_tests/Bugs/clean-segment-optional`
- Full suite: 132 passed, 0 failed

## 2026-07-30 — B16: `extractAttrValues` spaces in braces

- Fix B16: `core/codegen_template.go` `extractAttrValues` now tracks brace depth
- Spaces inside `{...}` expressions are no longer treated as attribute separators
- Integration test added: `_tests/Bugs/component-attr-space`
- Full suite: 133 passed, 0 failed

## 2026-07-30 — B17: `atoi` silently eats non-digits

- Fix B17: `core/validate.go` `atoi` now returns `(int, error)` and rejects empty/non-digit input
- `applyRule` for `min=`/`max=` returns "must be a valid number" on invalid input
- Unit tests `TestAtoi`, `TestAtoiEmpty`, `TestAtoiNonDigits`, `TestApplyRuleMinNonNumeric`, `TestApplyRuleMaxNonNumeric` updated
- Integration test added: `_tests/Bugs/validate-atoi-non-digit`
- Full suite: 134 passed, 0 failed

## 2026-07-30 — B18: `unindent` spaces

- Fix B18: `core/codegen.go` `unindent` now strips both tabs and spaces (`TrimLeft(l, " \t")`)
- Fix B18: `splitGoSections` passes raw `g.Code` (not trimmed) to `unindent` so consistent indentation can be removed
- Integration test added: `_tests/Bugs/unindent-spaces`
- Full suite: 135 passed, 0 failed

## 2026-07-30 — B19: `findFormStruct` regex

- Fix B19: `core/forms.go` `findFormStruct` regex updated to `func\s+Action\s*\(\s*\w+\s+[^,]+,\s*\w+\s+([^,)]+)\s*\)`
- Now matches complex parameter types (pointers, slices) and named return values
- Integration test added: `_tests/Bugs/form-handler-named-return`
- Full suite: 136 passed, 0 failed

## 2026-07-30 — Random ports + test cleanup standard

- All 116 `test.sh` files converted to standardized pattern: `mktemp -d`, `trap`, `go run`
- `_tests/how-to-test-sh.md` created with template and rules
- 17 server tests migrated from `:8080` to random ports (5-digit via `awk`), flaky port conflicts eliminated
- Known issue: `smd` Docker container fails on macOS due to `/var/folders` symlink — run `make test` directly

## 2026-07-29 — v0.0.17 Production Deployment + Request-ID

- Graceful shutdown: `http.Server` + SIGINT/SIGTERM in `core.Listen`, 10s drain timeout
- `dreego build --target linux/amd64` — cross-compile with GOOS/GOARCH
- Request-ID middleware: `X-Request-ID` header (16-char hex), context injection, JSONL log field `rid`
- `c.RequestID()` accessor on Context interface + SSRContext
- Production Dockerfile: `FROM scratch`, multi-stage, CGO_ENABLED=0, static binary
- `_docs/hot-reload.md`: Air config + entr alternative
- Rejected: hot-reload.1, live-reload.1, smart-recompile.1 — replaced by Air docs
- Block: request-id.1 completed (chain 34)

## 2026-07-29 — v0.0.16 Form Actions

- `<form g-action="Login">` — declarative server-side form handling
- Generated POST handler pipeline: BindForm → ValidateForm → Handler → Redirect
- `c.Redirect(url, code)` — PRG pattern with ErrRedirect sentinel
- `c.Errors(field)` / `c.Old(field)` — template accessors for validation state
- `BindForm()`, `ValidateForm()`, `SaveOld()`, `SaveErrors()` — no external deps
- `form:` and `validate:` struct tags — automatic form mapping and validation
- `splitGoSections` separates type/func declarations from inline code in codegen
- Context interface extended with session + redirect methods
- `scanFormActions` detects g-action in template, wires matching handlers
- 15 new tests (11 parser/codegen + 4 runtime HTTP), 112 total
- `_docs/forms.md` created, README expanded

## 2026-07-28 — v0.0.14–v0.0.15 Production Middleware + Content-Type Routing

- Health checks, security headers, gzip compression middleware
- Content-type routing: `<go type="json|xml|custom">` with `c.JSON()`, `c.XML()`, `c.Write()`
- 10 new runtime HTTP tests, curl-based Docker requests
- Fix: layout.Head.Content rendering, c.Wants() empty Accept handling

## 2026-07-25

- Route segments: `[id]` (square brackets) as convention — Demo migrated from `_id_`
- Route groups: `(group)/` — Folders without URL presence, `patternSegment` returns ""
- Flat gen package: all handlers in `gen/routes.go`, `gen/dree.go` without route imports
  - Solves Go import path problem with `[` `(` in directory names
  - `sanitizeImportPath` removed, no more `_ "import"`

- Context refactoring: `map[string]string` → Interface + Embedding
  - New `Context` interface embedding `context.Context` with `Param`, `Data`, `Render`
  - New `SSRContext` struct with `NewSSR()`, `Set`, `Get`, `Query`, `FormValue`
  - Updated codegen to generate `*context.SSRContext` + `context.NewSSR(w, r)`
  - Deleted stale `demo/dreego/layouts/dree.go` (layout is now inlined per route)
- Recovery middleware: Panic → 500 + stack trace logging via slog
  - New `dreego-core/recovery.go` — defer recover() with JSON log
  - Integrated as outermost core-fixed middleware in runtime pipeline
- XSS protection: Auto-escaping of all `{variable}` template expressions via `html.EscapeString`
  - Expression nodes now generate `html.EscapeString(fmt.Sprintf("%v", expr))`
  - `"html"` import only added when expressions exist (conditional in codegen)
- Custom error pages: `404.dreego` + `500.dreego`
  - `GenerateErrorHandler` in codegen: no layout, custom init (catch-all / SetErrorHandler)
  - Per-directory 404: Go Mux pattern precedence selects most specific 404
    - `dreego/routes/users/404.dreego` → `/users/{p...}` (only under `/users/*`)
    - `dreego/routes/404.dreego` → `/{p...}` (global fallback)
    - No 404 present → standard HTTP 404 text
  - 500: `runtime.SetErrorHandler(500, handler)` → Recovery middleware renders on panic
 - Converted entire knowledge base to Open Knowledge Format (OKF) v0.1

## 2026-07-28

### v0.0.15 — Content-Type Routing
- `<go type="json">` / `<go type="xml">` — typed Go blocks with `c.JSON()`, `c.XML()`, `c.Bind()`
- Multiple `<go>` blocks merged: shared runs always, typed checks `Accept` header via `c.Wants()`
- Pure typed routes skip template rendering (no `<div>` needed)
- `c.Write(status, contentType, body)` for arbitrary formats (FlatBuffers, Protobuf, etc.)
- `dreego-core/response.go`: 55 lines, stdlib-only (encoding/json, encoding/xml, net/http, strings)
- 86 integration tests total (5 new)

### v0.0.14 — Production Middleware
- Health checks: `GET /health` (liveness) + `GET /ready` (readiness via `core.SetReady(bool)`), core-fixed, registered before user routes
- Security headers: `X-Content-Type-Options`, `X-Frame-Options`, `Referrer-Policy`, `Permissions-Policy` — core-fixed middleware
- Gzip compression: `compress/gzip` wrapping ResponseWriter, checks `Accept-Encoding`, core-fixed
- Middleware chain: Recovery → SecurityHeaders → Compression → RequestLogging → Session → CSRF → Redirect/Rewrite → Router
- Bugfix: `layout.Head.Content` now rendered in codegen — `{#head}` inside `<head>` section processed correctly
- Bugfix: `genLayoutNode` handles both `{#head}` and `{#slot}` in same NodeText (splitLayoutText function)
- `dreego new`: force `go 1.22` in `go.mod` for Docker compatibility, bump core dep to v0.0.14
- 81 integration tests total (2 new regression tests + 1 updated)

### Recent Fixes (post v0.0.13)
- `.gitattributes`: Added `_todo export-ignore` — exclude `_todo/` from `go get` downloads
- `.gitignore`: Fixed `/dreego` (root binary only, not `cmd/dreego/` directory)
- `.gitignore`: Changed `**/gen/routes.go` → `dreego/gen/` (covers all split-gen files)
- Go embed cross-module fix: removed go.mod from blueprint (generated by `cmdNew`)

### v0.0.13 — Scaffolding + Split-Gen + Landing Blueprint
- `dreego new <name>` — scaffolding from landing blueprint with auto `go mod init`+`go mod edit`
- Landing blueprint: Tailwind CDN, Hero + FeatureCard components, layout with `{#head}`+`{#slot}`
- Blueprint placeholder: `§$name$§` replaced via `filepath.Base()`
- Split-Gen: `gen/routes.go` + `gen/components.go` + `gen/dree.go` (config+static)
- `isUpToDate()` file-level caching: file written only when content changes
- Dockerfile (golang:1.22-alpine → distroless nonroot)
- `.gitignore` fix: `/dreego` (root only), `dreego/gen/` (covers all generated files)
- Blueprint no go.mod: generated via `go mod init` in `cmdNew`
- 74 integration tests total

### v0.0.12 — Formatter
- `dreego fmt` — formats `.dreego` files in-place using `core.Format()`
- `--check` CI mode, `--stdout` mode
- `dreego-core/fmt.go`: component headers, expressions, control flow, section ordering
- Idempotency test

### v0.0.11 — English Translation + CLI Docs
- Full repo English translation: 130+ files (agents, docs, blocks, demo, guides)
- `dreego docs [--web] [--json] [--dump] [path]` — fetches from Codeberg
- URL filter strips Codeberg base for clean CLI output
- `dreego feedback` opens issues/new
- VS Code extension v0.0.4: component tags, imports, filters, `$loop` syntax highlighting

### v0.0.10 — Static Assets
- `dreego/static/` folder: Files are read during generate and registered inline as `[]byte`
- MIME type via extension (.css, .js, .svg, .png, .ico, .html, .json, .woff2)
- Collision check: static path vs route pattern → Error
- `core.RegisterStatic(path, mime, content)` in runtime
- 3 static tests (basic, subdir, collision), 71 integration tests total

### v0.0.9 — Template Primitives
- `{#verbatim}` block: Raw output for JS templates, content output verbatim
- `{#each}` with `$loop` variable: `$loop.Index`, `.First`, `.Last`, `.Even`, `.Odd`
- Template filters: `{var|raw}` (no escaping), `{var|upper}` (uppercase). Pipe syntax.
- `{#else}` in `{#if}` block
- `{#each else}`: Empty list fallback — `if len(items) > 0 { for } else { }`
- Fix: `<header>`, `<main>`, `<footer>` prefix-match bug in scanTag (tag terminator check)
- Test system: `PASS/FAIL <path>` format, `DREEGO_FILTER=<pattern>`, Docker build logs suppressed

### v0.0.8 — Named Slots
- Named slots: `{#slot header}...{/slot}` block syntax
- Component: `{#slot header}{/slot}` — placeholder, Route: `{#slot header}<content>{/slot}` — definition
- `c.Set("slot_name", ...)` / `c.Get("slot_name")` in codegen
- Default slot `{#slot}` remains without `{/slot}`
- 4 positive tests + 2 negative tests

### v0.0.7 — Test Coverage
- 41+ integration tests (up from 36): edge cases, negative tests
- `_docs/testing.md`: 60+ test ideas
- Bug fixes: `extractAttrValues` for `{expr}` props, `scanComponents` path filter

### v0.0.6 — Component Children
- Children slot passing: `<@Card>content</@Card>` → `{#slot}`
- `dreego generate --check`: timestamp comparison, exit non-zero if stale
- 36 tests

### v0.0.5 — Component System
- `Component Name (props)` declaration in `dreego/components/`
- `<@Name>` self-closing + children
- Auto-discovery via `scanComponents()`
- Scoped CSS per component (data-scope)
- `import` statement parsing in `ParseHeader`
- 35 tests

### v0.0.4 — Blueprints + Scaffolding
- `dreego init <path>` with `//go:embed` blueprints in `cmd/dreego/blueprints/default/`
- Docker integration tests: `_tests/Dockerfile` + `_tests/test.sh` orchestrator
- 24 integration tests, `make test` target
- Repo restructured: `pkg/` → `dreego-core/` (single package)
- `.gitattributes` `export-ignore` for `go get`
- Added YAML frontmatter with `type` field to all files
- Replaced `[[wiki-links]]` with standard markdown links
- Renamed `_index.md` to `index.md` with OKF child-list format
- Moved `thinking-list.md` to `PlanTODO.md`
- Moved `gemini-50-tipps.md` to `tips.md`
- Added `open-knowledge-format.md` skill guide
- Restructured `ROADMAP.md` as release pipeline
- Split `TODO.md` into immediate actions only

## 2026-07-24

- Decision: Error handling with typed errors, recovery middleware
- Decision: Routing & Components — hybrid file-based + programmatic routing
- Decision: Output strategy comparison (per-dir dree.go confirmed by GLM)
- Implemented gen/dree.go central import file
- Created dreego/config.json with redirects, rewrites, logging config
- Added RequestLogging middleware (JSONL, core-conditional)
- CLI commands: generate, build, run (-d, -t)

## 2026-07-23

- Initial knowledge base created
- Decision: Name "dreego", file extension ".dreego"
- Decision: Technology stack (Go 1.22+, net/http, HTMX, Alpine.js)
- Decision: Compile-Time Transpiler approach
- Decision: Transpiler Pipeline (Lexer → Parser → AST → CodeGen)
- Decision: TypeScript deferred to V2
- Decision: 5 sections in .dreego files
- Decision: SSR-First architecture
- Decision: No catch tag — errors via Go idioms
- Decision: File-based routing
- Decision: SSG & Wails in V2
- Researched: React/Svelte ecosystems, Phoenix/Laravel/Django, Rust frameworks, Solid/Astro, Blazor
- Concepts drafted: Architecture, Sections, Template Logic, Addon Ecosystem
