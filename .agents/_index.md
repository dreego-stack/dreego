# dreego Knowledge Base

## Overview

**dreego** — Ein Go-Webframework, das das Beste aus Svelte und Go vereint.

Package: `dreego` | Dateiendung: `.dreego` | CLI: `dreego`

Ziel: Ein SSR-First Webframework für Go, das auf Augenhöhe mit Phoenix, Next.js und SvelteKit spielt — aber als Single Binary deploybar.

## Aktuelle Phase: Planung / Konzeption

Dreego befindet sich in der Konzeptionsphase. Noch kein Code, nur Architektur-Planung.

**Wichtige Dokumente:**
- [[../thinking-list]] — Offene Punkte, die vor/nach Code-Start geklärt werden müssen
- [[concepts/gap-analysis]] — Was fehlt Dreego? Wo muss nachgebessert werden?

## Knowledge Base

### Referenz-Material
- [[KB/dreego-concept]] — Vollständiges Konzept aus dem Gemini-Chat (PDF-Quelle)
- [[KB/ecosystem-research-2025-2026]] — React/Svelte Ecosystem Research (State of JS 2025)
- [[KB/ecosystem-analysis]] — Was wir von React/Svelte übernehmen und was nicht
- [[KB/solid-astro-mdx-research]] — Solid.js, Astro, MDX Deep-Dive
- [[KB/framework-research-phoenix-laravel-django]] — Phoenix, Laravel, Django Analyse
- [[KB/rust-frameworks-analysis]] — Rust Frontend-Frameworks (Leptos, Dioxus, Yew)
- [[KB/blazor-research]] — C# Blazor & ASP.NET Core Analyse
- [[KB/framework-research-phoenix-laravel-django]] — Phoenix/Laravel/Django Feature-Analyse für Dreego
- [[KB/rust-frameworks-analysis]] — Rust-Webframeworks (Leptos, Dioxus, Yew): Patterns für Go-SSR
- [[KB/solid-astro-mdx-research]] — Solid.js, Astro & MDX — Deep-Dive (23.07.2026)

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

- [[decisions/name-dreego]] — Namensgebung: edreego / .dreego
- [[decisions/technology-stack]] — Tech-Stack: Go, Chi, Tailwind, HTMX, Alpine.js, Datastar
- [[decisions/transpiler-vs-runtime]] — Compile-Time Transpiler statt Runtime-Parsing
- [[decisions/typescript-v2]] — TypeScript auf V2 verschoben, V1 nur Vanilla JS
- [[decisions/sections-in-dreego]] — 5 Sektionen: `<head>`, `<go>`, Template, `<script>`, `<style>`
- [[decisions/ssr-first]] — SSR-First Ansatz, kein Client-Side Framework
- [[decisions/no-catch-tag]] — Kein `<catch>`-Tag, Fehler via Go-Idiome behandeln
- [[decisions/file-based-routing]] — File-based Routing mit Chi
- [[decisions/ssg-wails-v2]] — SSG & Wails erst in V2, Architektur in V1 vorbereitet
- [[decisions/context-design]] — Context: Interface + Target-Structs (GLM-5.2 Review)

## Guides

- [[guides/architecture]] — Projektstruktur und Modul-Grenzen
- [[guides/coding-standards]] — Coding-Regeln für Dreego

## Tech Stack (V1)

| Bereich            | Wahl                                      |
|--------------------|-------------------------------------------|
| Sprache            | Go 1.26+                                  |
| HTTP-Router        | go-chi/chi                                |
| Template Engine    | Dreego Transpiler (.dreego → Go-Code)       |
| Interaktivität     | HTMX + Alpine.js + Datastar (optional)    |
| CSS                | Tailwind CLI (embedded)                   |
| Binary Packaging   | Go `//go:embed` (Single Binary)           |
| Validierung        | go-playground/validator                   |
| Dev-Server         | SSE-basiertes Hot Reload                  |

## Affons Ecosystem

Siehe [[concepts/affons-ecosystem]] für die Planung des umgebenden Tool-Ökosystems.
