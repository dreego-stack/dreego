---
type: Index
title: Dreego Knowledge Base
description: Agent Knowledge Bundle for Dreego — Go web framework v0.0.22
tags: [index, v0.0.21]
timestamp: 2026-07-31T07:00:00Z
---

# Dreego Knowledge Base

## Decisions

- [Name "dreego"](decisions/name-dreego.md) — Naming and package convention
- [Monorepo Plugin Layout](decisions/monorepo-plugin-layout.md) — Official plugins live in this repo under plugins/ (supersedes separate-repo model)
- [Technology Stack](decisions/technology-stack.md) — Tech stack: Go, net/http, HTMX, Alpine.js
- [Transpiler vs Runtime](decisions/transpiler-vs-runtime.md) — Compile-time code generation
- [Transpiler Pipeline](decisions/transpiler-pipeline.md) — Lexer → Parser → AST → CodeGen
- [TypeScript V2](decisions/typescript-v2.md) — TypeScript deferred to V2
- [5 Sections](decisions/sections-in-dreego.md) — head, go, div, script, style
- [SSR First](decisions/ssr-first.md) — Server-side rendering as default
- [No Catch Tag](decisions/no-catch-tag.md) — Errors via {#if hasError}, no separate tag
- [File-based Routing](decisions/file-based-routing.md) — dreego/routes/*.dreego
- [Routing & Components](decisions/routing-and-components.md) — Hybrid routing, plugin routes, components
- [SSG & Wails V2](decisions/ssg-wails-v2.md) — Static site generation and desktop in V2
- [Context Design](decisions/context-design.md) — Interface + embedding per target
- [Plugin Interface](decisions/plugin-interface.md) — Capability-based plugin system
- [Session Management](decisions/session-management.md) — Interface in core, store as plugin
- [Error Handling](decisions/error-handling.md) — Typed errors, recovery, request logging
- [Middleware System](decisions/middleware-system.md) — Core-fixed vs core-conditional vs plugin
- [Form Actions](decisions/form-actions.md) — g-action + generated pipeline

## Concepts

- [Dreego Architecture](concepts/dreego-architecture.md) — Architecture overview
- [Dreego Sections](concepts/dreego-sections.md) — The 5 sections of a .dreego file
- [Template Logic](concepts/template-logic.md) — {#if}, {#each}, {#switch}, {#slot}
- [Plugin Ecosystem](concepts/plugin-ecosystem.md) — Plugin architecture
- [Signals & Runes](concepts/signals-and-runes.md) — Signals and Svelte runes in Dreego
- [Output Strategy](concepts/output-strategy-comparison.md) — Comparison of output strategies
- [Form Actions](concepts/form-actions.md) — Form handling concept

## Reference (KB)

- [Dreego Concept](KB/dreego-concept.md) — Gemini chat concept (source)
- [Ecosystem Analysis](KB/ecosystem-analysis.md) — React/Svelte: What Dreego adopts
- [Ecosystem Research 2025-2026](KB/ecosystem-research-2025-2026.md) — State of JS 2025
- [Solid, Astro, MDX](KB/solid-astro-mdx-research.md) — Solid.js, Astro, MDX research
- [Phoenix, Laravel, Django](KB/framework-research-phoenix-laravel-django.md) — Framework comparison
- [Rust Frameworks](KB/rust-frameworks-analysis.md) — Leptos, Dioxus, Yew
- [Blazor Research](KB/blazor-research.md) — C# Blazor & ASP.NET Core

## Guides

- [Architecture Guide](guides/architecture.md) — Project structure and module boundaries
- [Coding Standards](guides/coding-standards.md) — Coding rules for Dreego
- [Knowledge Base](guides/knowledge-base.md) — Skill: Maintaining the KB
- [Changelog](guides/changelog.md) — Skill: CHANGELOG and versioning
- [Open Knowledge Format](guides/open-knowledge-format.md) — Skill: OKF conventions

## Plans & Tips

- [Open work](../_todo/) — One item file per open task or idea
- [Tips](tips.md) — 50 tips + checklist

## External Docs

- [README](../README.md) — Project README
- [CHANGELOG](../CHANGELOG.md) — Version history
- [Docs](../_docs/index.md) — Public documentation
