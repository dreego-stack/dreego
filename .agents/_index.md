# dreego Knowledge Base

## Overview

**dreego** — Ein Go-Webframework, das das Beste aus Svelte und Go vereint.

Package: `dreego` | Dateiendung: `.dreego` | CLI: `dreego`

Ziel: Ein SSR-First Webframework für Go, das auf Augenhöhe mit Phoenix, Next.js und SvelteKit spielt — aber als Single Binary deploybar.

## Aktuelle Phase: 0.0.x — Implementierung

Transpiler + Routing + Layout + CSS-Scoping funktionieren. Siehe [[../ROADMAP]].

**Wichtige Dokumente:**
- [[../ROADMAP]] — Aktuelle Roadmap mit Phasen 0–3
- [[../CLI]] — CLI-Kommandos
- [[thinking-list]] — Detaillierte Feature-Liste

## Knowledge Base

### Referenz-Material
- [[KB/dreego-concept]] — Gemini-Chat-Konzept (PDF-Quelle)
- [[KB/ecosystem-analysis]] — React/Svelte: Was übernehmen, was nicht
- [[KB/ecosystem-research-2025-2026]] — State of JS 2025 Research
- [[KB/solid-astro-mdx-research]] — Solid.js, Astro, MDX
- [[KB/framework-research-phoenix-laravel-django]] — Phoenix, Laravel, Django
- [[KB/rust-frameworks-analysis]] — Rust Frameworks (Leptos, Dioxus, Yew)
- [[KB/blazor-research]] — C# Blazor & ASP.NET Core

### Konzepte
- [[concepts/dreego-architecture]] — Architektur-Übersicht
- [[concepts/dreego-sections]] — Die 5 Sektionen einer .dreego-Datei
- [[concepts/template-logic]] — Template-Logik: {#if}, {#each}, {#switch}, {#slot}, {#await}
- [[concepts/addon-ecosystem]] — Addon/Plugin-Architektur
- [[concepts/affons-ecosystem]] — Gesamtes Ecosystem (CLI, Docs, Registry, Community)
- [[concepts/signals-and-runes]] — Signals, Svelte Runes und wie Dreego sie abbildet
- [[concepts/gap-analysis]] — Was fehlt Dreego? Wo muss nachgebessert werden?
- [[concepts/roadmap]] — V1, V2, und darüber hinaus

## Decisions

- [[decisions/name-dreego]] — Namensgebung: dreego / .dreego
- [[decisions/technology-stack]] — Tech-Stack: Go, net/http, Tailwind, HTMX, Alpine.js
- [[decisions/transpiler-vs-runtime]] — Compile-Time Transpiler
- [[decisions/transpiler-pipeline]] — 0.0.1: Single-Pass Scanner, dann formale Pipeline
- [[decisions/typescript-v2]] — TypeScript auf V2
- [[decisions/sections-in-dreego]] — 5 Sektionen
- [[decisions/ssr-first]] — SSR-First
- [[decisions/no-catch-tag]] — Kein `<catch>`-Tag
- [[decisions/file-based-routing]] — File-based Routing mit net/http
- [[decisions/ssg-wails-v2]] — SSG & Wails in V2
- [[decisions/context-design]] — Context: Interface + Target-Structs
- [[decisions/plugin-interface]] — Plugin: Capability-basiert
- [[decisions/transpiler-pipeline]] — Lexer→Parser→AST→CodeGen
- [[decisions/session-management]] — Session: Interface im Core
- [[decisions/middleware-system]] — Middleware: Core vs Plugin
- [[decisions/form-actions]] — g-action + generierte Pipeline

## Tech Stack (aktuell)

| Bereich            | Wahl                                      |
|--------------------|-------------------------------------------|
| Sprache            | Go 1.22+                                  |
| HTTP-Router        | net/http (Go 1.22+ PathValue)            |
| Template Engine    | Dreego Transpiler (.dreego → Go-Code)     |
| CSS Scoping        | `data-scope` via Source-Hash             |
| Binary Packaging   | Go `//go:embed`                          |
| Dev-Server         | `air` (Hot Reload)                       |
