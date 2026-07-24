# dreego Knowledge Base

## Overview

**dreego** — Ein Go-Webframework, das das Beste aus Svelte und Go vereint.

Package: `dreego` | Dateiendung: `.dreego` | CLI: `dreego`

Ziel: Ein SSR-First Webframework fur Go, das auf Augenhohe mit Phoenix, Next.js und SvelteKit spielt — aber als Single Binary deploybar.

## Aktuelle Phase: 0.0.1 — Erster Prototyp

Transpiler + Routing + Layout + CSS-Scoping + CLI + Middleware funktionieren. v0.0.1 getaggt.

**Wichtige Dokumente:**
- [[../README]] — Projekt-Readme
- [[../CHANGELOG]] — Versionshistorie
- [[../ROADMAP]] — Roadmap mit Phasen 0–3
- [[../_docs/index]] — Dokumentation
- [[thinking-list]] — Detaillierte Feature-Liste

## Knowledge Base

### Referenz-Material
- [[KB/dreego-concept]] — Gemini-Chat-Konzept (PDF-Quelle)
- [[KB/ecosystem-analysis]] — React/Svelte: Was ubernehmen, was nicht
- [[KB/ecosystem-research-2025-2026]] — State of JS 2025 Research
- [[KB/solid-astro-mdx-research]] — Solid.js, Astro, MDX
- [[KB/framework-research-phoenix-laravel-django]] — Phoenix, Laravel, Django
- [[KB/rust-frameworks-analysis]] — Rust Frameworks (Leptos, Dioxus, Yew)
- [[KB/blazor-research]] — C# Blazor & ASP.NET Core

### Konzepte
- [[concepts/dreego-architecture]] — Architektur-Ubersicht
- [[concepts/dreego-sections]] — Die 5 Sektionen einer .dreego-Datei
- [[concepts/template-logic]] — Template-Logik
- [[concepts/addon-ecosystem]] — Addon/Plugin-Architektur
- [[concepts/affons-ecosystem]] — Gesamtes Ecosystem
- [[concepts/signals-and-runes]] — Signals & Svelte Runes
- [[concepts/output-strategy-comparison]] — Output-Strategie-Vergleich
- [[concepts/gap-analysis]] — Gap-Analyse
- [[concepts/roadmap]] — V1, V2, und daruber hinaus

## Decisions

- [[decisions/name-dreego]] — Namensgebung
- [[decisions/technology-stack]] — Tech-Stack
- [[decisions/transpiler-vs-runtime]] — Compile-Time Transpiler
- [[decisions/transpiler-pipeline]] — Lexer→Parser→AST→CodeGen
- [[decisions/typescript-v2]] — TypeScript auf V2
- [[decisions/sections-in-dreego]] — 5 Sektionen
- [[decisions/ssr-first]] — SSR-First
- [[decisions/no-catch-tag]] — Kein `<catch>`-Tag
- [[decisions/file-based-routing]] — File-based Routing
- [[decisions/routing-and-components]] — Hybrides Routing + Komponenten
- [[decisions/ssg-wails-v2]] — SSG & Wails in V2
- [[decisions/context-design]] — Context-Design
- [[decisions/plugin-interface]] — Plugin-Interface
- [[decisions/session-management]] — Session-Management
- [[decisions/error-handling]] — Error-Handling + RequestLogging
- [[decisions/middleware-system]] — Middleware-System
- [[decisions/form-actions]] — Form Actions

## Tech Stack (aktuell)

| Bereich            | Wahl                                   |
|--------------------|----------------------------------------|
| Sprache            | Go 1.22+                               |
| HTTP-Router        | net/http (Go 1.22+ PathValue)          |
| Template Engine    | Dreego Transpiler (.dreego → Go-Code)  |
| CSS Scoping        | `data-scope` via Source-Hash           |
| CLI                | dreego generate, build, run            |
| Logging            | JSONL via slog (Core-Conditional)      |
| Dev-Server         | `dreego run -d -t N`                   |
