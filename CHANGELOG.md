# Changelog

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

- [[.agents/decisions/error-handling]] — Typisierte Fehler, Recovery, Dev/Prod
- [[.agents/decisions/routing-and-components]] — Hybrides Routing, Plugin-Routes, Komponenten
- [[.agents/decisions/middleware-system]] — Core-Fixed vs Core-Conditional vs Plugin
- GLM-Review: Per-Directory `dree.go` Output-Strategie bestatigt
