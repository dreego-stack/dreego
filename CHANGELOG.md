# Changelog

## v0.0.6 (unreleased) — Component Completion

- Children-Slot-Passing: `<@Card>content</@Card>` → `{#slot}` im Component funktioniert
- `dreego generate --check`: CI-Validation — exit non-zero wenn generierte Dateien stale sind
- Named Slots: `{#slot header}` Lexer/Parser vorbereitet (v0.0.7)

## v0.0.5 (2026-07-27) — Component Model

- Component-System: `Component Name (props)` in `dreego/components/`, Aufruf via `<@Name>`
- Self-closing (`<@Icon name="star"/>`) und mit Children (`<@Card>...</@Card>`)
- Default-Slot via `{#slot}` im Component-Template
- Scoped Styles pro Component (`data-scope`)
- File-based Discovery: `dreego/components/Card.dreego` → `<@Card>`
- 6 Component-Integration-Tests + 2 Bug-Tests
- `import "dreego/components/Name"` in Route-Dateien (ParseHeader vor Lex)
- Multi-File Directory Import: `import "dreego/components/button"` → `<@Login/>`

## v0.0.4 (2026-07-27) — Blueprints & Tests

- `dreego init <path>` — scaffold new project from embedded blueprint
- Blueprints via `//go:embed` in CLI binary, keine externen Dateien nötig
- Integration-Tests in `_tests/` via Docker-Container (`make test`)

## v0.0.3 (2026-07-27) — Security & Developer Experience

- Session-Integration: `session.Store` Interface + `CookieStore` (HMAC-signiert) in Runtime eingehängt
- Session-Middleware: Context-basierte Store-Injektion pro Request
- SSRContext: `SessionVal`/`SetSessionVal`/`DelSessionVal`/`DestroySession` mit sicheren Defaults (`HttpOnly`, `Secure` TLS-aware)
- CSRF-Schutz: Double-Submit-Cookie (Core-Conditional, default an) — Token via X-CSRF-Token Header oder csrf_token Form-Feld
- SSRContext: `CSRFToken()` fur Template-Rendering (Hidden-Field)
- VS Code Extension: Syntax-Highlighting + Waschbär-Icon für `.dreego`-Dateien (`make dx`)
- **Breaking:** `pkg/` → `dreego-core/` (single package), einziger Import `import core "codeberg.org/dreego/dreego/dreego-core"`
- `dreego-plugin/` fur zukünftige Plugins (Auth, Redis, DB, etc.)

## v0.0.2 (2026-07-25) — Safety & Structure

- Route-Segmente: `[id]` (eckige Klammern) als Konvention für dynamische Segmente, kompatibel mit Next.js/SvelteKit/Astro
- Route-Groups: `(group)/` — Ordner, die nicht in der URL erscheinen (Layout/Middleware-Gruppierung)
- Flat Gen-Package: alle Route-Handler in `gen/routes.go` (keine `_ "import"` mehr), löst Go-Import-Path-Problem mit Sonderzeichen
- Context refactoring: `map[string]string` → Interface + Embedding (`Context` interface + `SSRContext` struct)
- Recovery-Middleware: Panic → 500 mit Stack-Trace-Logging via slog
- XSS-Schutz: Auto-Escaping aller `{variable}`-Ausdrücke via `html.EscapeString`
- Custom Error-Pages: `404.dreego` + `500.dreego`

## v0.0.1 (2026-07-25) — The Prototype

Erster Prototyp. Transpiler, Routing, Layout, Middleware, CLI.

### Features

- Formale Transpiler-Pipeline: Lexer → Parser → AST → CodeGen
- Alle 5 Sektionen: `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- Template-Logik: `{var}`, `{#if}`, `{#each}`, `{#slot}`, `{#head}`
- File-based Routing: `dreego/routes/*.dreego`
- Dynamische Segmente: `[id]`, `[...catchall]`, `[[optional]]`, `(group)/`
- Layout-System: `dreego/layouts/default.dreego` mit `{#slot}` + `{#head}`
- CSS-Scoping: `data-scope` via Source-Hash (12 Zeichen)
- Zentrale `dreego/gen/dree.go` fur Route-Imports
- `dreego/config.json`: Redirects, Rewrites, Logging-Config
- RequestLogging-Middleware (Core-Conditional, JSONL-Format, IP-Capture)
- Redirect/Rewrite-Middleware
- CLI: `dreego generate [--force]`, `dreego build`, `dreego run [-d] [-t N]`
- Lauffahiger Demo-Server mit net/http 1.22+

### Decisions

- [Error Handling](.agents/decisions/error-handling.md) — Typisierte Fehler, Recovery, Dev/Prod
- [Routing & Components](.agents/decisions/routing-and-components.md) — Hybrides Routing, Plugin-Routes
- [Middleware System](.agents/decisions/middleware-system.md) — Core-Fixed vs Core-Conditional vs Plugin
- GLM-Review: Per-Directory `dree.go` Output-Strategie bestatigt
