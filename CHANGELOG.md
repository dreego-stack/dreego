# Changelog

## v0.0.3 (unreleased)

- Session-Interface: `session.Store` im Core mit `CookieStore` (HMAC-signiert)

## v0.0.2 (2026-07-25)

- Route-Segmente: `[id]` (eckige Klammern) als Konvention für dynamische Segmente, kompatibel mit Next.js/SvelteKit/Astro
- Route-Groups: `(group)/` — Ordner, die nicht in der URL erscheinen (Layout/Middleware-Gruppierung)
- Flat Gen-Package: alle Route-Handler in `gen/routes.go` (keine `_ "import"` mehr), löst Go-Import-Path-Problem mit Sonderzeichen
- Context refactoring: `map[string]string` → Interface + Embedding (`Context` interface + `SSRContext` struct)
- Recovery-Middleware: Panic → 500 mit Stack-Trace-Logging via slog
- XSS-Schutz: Auto-Escaping aller `{variable}`-Ausdrücke via `html.EscapeString`
- Custom Error-Pages: `404.dreego` + `500.dreego`

## v0.0.1 (2026-07-25)

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
