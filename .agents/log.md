---
type: Log
title: Knowledge Base Changelog
description: Record of all changes to the .agents/ knowledge bundle
tags: [log]
timestamp: 2026-07-25T00:00:00Z
---

# log

## 2026-07-25

- Route-Segmente: `[id]` (eckige Klammern) als Konvention — Demo von `_id_` migriert
- Route-Groups: `(group)/` — Ordner ohne URL-Präsenz, `patternSegment` gibt "" zurück
- Flat Gen-Package: alle Handler in `gen/routes.go`, `gen/dree.go` ohne Route-Imports
  - Löst Go-Import-Path-Problem mit `[` `(` in Verzeichnisnamen
  - `sanitizeImportPath` entfernt, kein `_ "import"` mehr

- Context refactoring: `map[string]string` → Interface + Embedding
  - New `Context` interface embedding `context.Context` with `Param`, `Data`, `Render`
  - New `SSRContext` struct with `NewSSR()`, `Set`, `Get`, `Query`, `FormValue`
  - Updated codegen to generate `*context.SSRContext` + `context.NewSSR(w, r)`
  - Deleted stale `demo/dreego/layouts/dree.go` (layout is now inlined per route)
- Recovery-Middleware: Panic → 500 + Stack-Trace-Logging via slog
  - New `dreego-core/recovery.go` — defer recover() with JSON log
  - Integrated as outermost Core-Fixed middleware in runtime pipeline
- XSS-Schutz: Auto-Escaping aller `{variable}`-Template-Ausdrücke via `html.EscapeString`
  - Expression nodes generieren jetzt `html.EscapeString(fmt.Sprintf("%v", expr))`
  - `"html"` import wird nur bei Expressions eingefügt (conditional in codegen)
- Custom Error-Pages: `404.dreego` + `500.dreego`
  - `GenerateErrorHandler` in codegen: kein Layout, Custom init (catch-all / SetErrorHandler)
  - Per-Directory 404: Go Mux Pattern-Precedence selektiert spezifischsten 404
    - `dreego/routes/users/404.dreego` → `/users/{p...}` (nur unter `/users/*`)
    - `dreego/routes/404.dreego` → `/{p...}` (globaler Fallback)
    - Kein 404 vorhanden → Standard-HTTP-404-Text
  - 500: `runtime.SetErrorHandler(500, handler)` → Recovery-Middleware rendert bei Panic
 - Converted entire knowledge base to Open Knowledge Format (OKF) v0.1

## 2026-07-28

### v0.0.10 — Static Assets
- `dreego/static/` Ordner: Files werden beim Generate eingelesen und inline als `[]byte` registriert
- MIME-Type via Extension (.css, .js, .svg, .png, .ico, .html, .json, .woff2)
- Kollision-Check: statischer Pfad vs Route-Pattern → Error
- `core.RegisterStatic(path, mime, content)` in Runtime
- 3 Static-Tests (basic, subdir, collision), 71 Integration-Tests total

### v0.0.9 — Template-Primitives
- `{#verbatim}` Block: Raw-Output für JS-Templates, Content 1:1 ausgegeben
- `{#each}` mit `$loop`-Variable: `$loop.Index`, `.First`, `.Last`, `.Even`, `.Odd`
- Template-Filter: `{var|raw}` (kein Escaping), `{var|upper}` (uppercase). Pipe-Syntax.
- `{#else}` in `{#if}`-Block
- `{#each else}`: Empty-List Fallback — `if len(items) > 0 { for } else { }`
- Fix: `<header>`, `<main>`, `<footer>` prefix-match Bug in scanTag (tag terminator check)
- Test-System: `PASS/FAIL <path>` Format, `DREEGO_FILTER=<pattern>`, Docker-Build-Logs suppressed

### v0.0.8 — Named Slots
- Named Slots: `{#slot header}...{/slot}` Block-Syntax
- Component: `{#slot header}{/slot}` — Platzhalter, Route: `{#slot header}<content>{/slot}` — Definition
- `c.Set("slot_name", ...)` / `c.Get("slot_name")` in Codegen
- Default-Slot `{#slot}` bleibt ohne `{/slot}`
- 4 Positiv-Tests + 2 Negativ-Tests

### v0.0.7 — Test Coverage
- 41+ Integration-Tests (up from 36): edge cases, negative tests
- `_docs/testing.md`: 60+ Test-Ideen
- Bugfixes: `extractAttrValues` für `{expr}` Props, `scanComponents` path filter

### v0.0.6 — Component Children
- Children slot passing: `<@Card>content</@Card>` → `{#slot}`
- `dreego generate --check`: timestamp comparison, exit non-zero if stale
- 36 Tests

### v0.0.5 — Component System
- `Component Name (props)` Declaration in `dreego/components/`
- `<@Name>` self-closing + children
- Auto-discovery via `scanComponents()`
- Scoped CSS per Component (data-scope)
- `import` statement parsing in `ParseHeader`
- 35 Tests

### v0.0.4 — Blueprints + Scaffolding
- `dreego init <path>` mit `//go:embed` Blueprints in `cmd/dreego/blueprints/default/`
- Docker-Integration-Tests: `_tests/Dockerfile` + `_tests/test.sh` Orchestrator
- 24 Integration-Tests, `make test` target
- Repo restructured: `pkg/` → `dreego-core/` (single package)
- `.gitattributes` `export-ignore` für `go get`
- Added YAML frontmatter with `type` field to all files
- Replaced `[[wiki-links]]` with standard markdown links
- Renamed `_index.md` to `index.md` with OKF child-list format
- Moved `thinking-list.md` to `PlanTODO.md`
- Moved `gemini-50-tipps.md` to `tips.md`
- Added `open-knowledge-format.md` skill guide
- Restructured `ROADMAP.md` as release pipeline
- Split `TODO.md` into immediate actions only

## 2026-07-24

- Decision: Error-Handling with typed errors, recovery middleware
- Decision: Routing & Components — hybrid file-based + programmatic routing
- Decision: Output strategy comparison (per-dir dree.go confirmed by GLM)
- Implemented gen/dree.go central import file
- Created dreego/config.json with redirects, rewrites, logging config
- Added RequestLogging middleware (JSONL, Core-Conditional)
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
