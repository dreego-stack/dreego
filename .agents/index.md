---
type: Index
title: Dreego Knowledge Base
description: Agent Knowledge Bundle for the Dreego Go-Webframework
tags: [index, v0.0.1]
timestamp: 2026-07-25T00:00:00Z
---

# Dreego Knowledge Base

## Decisions

- [Name "dreego"](decisions/name-dreego.md) — Namensgebung und Package-Konvention
- [Technology Stack](decisions/technology-stack.md) — Tech-Stack: Go, net/http, HTMX, Alpine.js
- [Transpiler vs Runtime](decisions/transpiler-vs-runtime.md) — Compile-Time Code-Generation
- [Transpiler Pipeline](decisions/transpiler-pipeline.md) — Lexer → Parser → AST → CodeGen
- [TypeScript V2](decisions/typescript-v2.md) — TypeScript auf V2 verschoben
- [5 Sektionen](decisions/sections-in-dreego.md) — head, go, div, script, style
- [SSR First](decisions/ssr-first.md) — Server-Side Rendering als Default
- [No Catch Tag](decisions/no-catch-tag.md) — Fehler via {#if hasError}, kein eigenes Tag
- [File-based Routing](decisions/file-based-routing.md) — dreego/routes/*.dreego
- [Routing & Components](decisions/routing-and-components.md) — Hybrides Routing, Plugin-Routes, Komponenten
- [SSG & Wails V2](decisions/ssg-wails-v2.md) — Static Site Generation und Desktop in V2
- [Context Design](decisions/context-design.md) — Interface + Embedding pro Target
- [Plugin Interface](decisions/plugin-interface.md) — Capability-basiertes Plugin-System
- [Session Management](decisions/session-management.md) — Interface im Core, Store als Addon
- [Error Handling](decisions/error-handling.md) — Typisierte Fehler, Recovery, RequestLogging
- [Middleware System](decisions/middleware-system.md) — Core-Fixed vs Core-Conditional vs Plugin
- [Form Actions](decisions/form-actions.md) — g-action + generierte Pipeline

## Concepts

- [Dreego Architecture](concepts/dreego-architecture.md) — Architektur-Übersicht
- [Dreego Sections](concepts/dreego-sections.md) — Die 5 Sektionen einer .dreego-Datei
- [Template Logic](concepts/template-logic.md) — {#if}, {#each}, {#switch}, {#slot}
- [Addon Ecosystem](concepts/addon-ecosystem.md) — Plugin-Architektur
- [Affons Ecosystem](concepts/affons-ecosystem.md) — CLI, Docs, Registry, Community
- [Signals & Runes](concepts/signals-and-runes.md) — Signals und Svelte Runes in Dreego
- [Output Strategy](concepts/output-strategy-comparison.md) — Vergleich der Output-Strategien
- [Gap Analysis](concepts/gap-analysis.md) — Was fehlt, was kann besser werden
- [Roadmap](concepts/roadmap.md) — V1, V2, und daruber hinaus
- [Form Actions](concepts/form-actions.md) — Form-Handling-Konzept

## Reference (KB)

- [Dreego Concept](KB/dreego-concept.md) — Gemini-Chat-Konzept (Quelle)
- [Ecosystem Analysis](KB/ecosystem-analysis.md) — React/Svelte: Was ubernimmt Dreego
- [Ecosystem Research 2025-2026](KB/ecosystem-research-2025-2026.md) — State of JS 2025
- [Solid, Astro, MDX](KB/solid-astro-mdx-research.md) — Solid.js, Astro, MDX Research
- [Phoenix, Laravel, Django](KB/framework-research-phoenix-laravel-django.md) — Framework-Vergleich
- [Rust Frameworks](KB/rust-frameworks-analysis.md) — Leptos, Dioxus, Yew
- [Blazor Research](KB/blazor-research.md) — C# Blazor & ASP.NET Core

## Guides

- [Architecture Guide](guides/architecture.md) — Projektstruktur und Modul-Grenzen
- [Coding Standards](guides/coding-standards.md) — Coding-Regeln fur Dreego
- [Knowledge Base](guides/knowledge-base.md) — Skill: KB pflegen
- [Changelog](guides/changelog.md) — Skill: CHANGELOG und Versionierung
- [Open Knowledge Format](guides/open-knowledge-format.md) — Skill: OKF-Konventionen

## Plans & Tips

- [PlanTODO](PlanTODO.md) — Vollstandiger Plan aller Features
- [Tips](tips.md) — 50 Tipps + Beachtungsliste

## External Docs

- [README](../README.md) — Projekt-Readme
- [TODO](../TODO.md) — Nachste Code-Anderungen
- [ROADMAP](../ROADMAP.md) — Release-Pipeline
- [CHANGELOG](../CHANGELOG.md) — Versionshistorie
- [Docs](../_docs/index.md) — Offentliche Dokumentation
