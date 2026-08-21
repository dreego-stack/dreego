---
type: Index
title: Dreego Knowledge Base
description: Agent Knowledge Bundle for Dreego — Go web framework v0.0.63
tags: [index, v0.0.63]
timestamp: 2026-08-15T00:00:00Z
---

# Dreego Knowledge Base

Status legend: **Current** = the v0.1 contract, **Provisional** = accepted for
v0.1 but not yet stable, **Research** = background material, **Superseded** =
historical decisions kept for context. A future agent can identify the current
v0.1 contract from the Current and Provisional sections alone.

## Decisions — Current

- [Name "dreego"](decisions/name-dreego.md) — Naming and package convention
- [Compile-Time Transpiler](decisions/transpiler-vs-runtime.md) — Build-time code generation, no runtime parsing
- [Transpiler Pipeline](decisions/transpiler-pipeline.md) — Lexer → Parser → AST → CodeGen (SSR-only for v0.1)
- [SSR First](decisions/ssr-first.md) — Server-side rendering as default; HTMX/Alpine/plain JS for interactivity
- [Context Design](decisions/context-design.md) — `dreego.Context` interface, no hard `*http.Request` in `<go>` blocks
- [File-based Routing](decisions/file-based-routing.md) — dreego/routes/*.dreego (amended by routing-and-components)
- [Routing & Components](decisions/routing-and-components.md) — App-bound registration, one route file per URL, typed components
- [Session Management](decisions/session-management.md) — Interface in core, store as plugin
- [Error Handling](decisions/error-handling.md) — Typed errors, recovery, request logging
- [Middleware System](decisions/middleware-system.md) — Core-fixed vs core-conditional vs plugin
- [Form Actions](decisions/form-actions.md) — g-action + generated pipeline (built-in validators)
- [No Catch Tag](decisions/no-catch-tag.md) — Errors via {#if hasError}, no separate tag
- [Sections in Dreego](decisions/sections-in-dreego.md) — head, go, div, script, style
- [Line Limit 300](decisions/line-limit-300.md) — Max 300 lines per file
- [Transpiler as internal subpackage](decisions/transpiler-subpackage.md) — Transpiler in `internal/transpiler/`, core is runtime-only
- [Core runtime split into internal subpackages](decisions/core-internal-subpackages.md) — `core/internal/` holds session/server/middleware/context/validate; `core/` re-exports the public API
- [Per-directory dree.go output](decisions/per-directory-dreego.md) — `dreego/gen/` is gone; each directory with `.dreego` sources gets its own `dree.go`, website roots are marked by `dreego.config.json`

## Decisions — Provisional

- [Plugin Registration](decisions/plugin-interface.md) — App-bound `Register(app, options)` functions; no stable Plugin interface before v1
- [TypeScript V2](decisions/typescript-v2.md) — TypeScript deferred to V2; V1 uses plain JavaScript

## Decisions — Superseded (historical)

- [Technology Stack](decisions/technology-stack.md) — Chi, validator, Tailwind CLI; superseded by net/http + built-in validators + dependency-free Core
- [Monorepo Plugin Layout](decisions/monorepo-plugin-layout.md) — Official plugins in this repo; superseded by separate plugin repositories
- [SSG & Wails V2](decisions/ssg-wails-v2.md) — Multi-target V2 timeline; superseded by SSR-only until v1
- [Transpiler Pipeline (3-target codegen)](decisions/transpiler-pipeline.md) — SSR/SSG/Wails codegen; superseded by SSR-only for v0.1
- [Context Design (3 target structs)](decisions/context-design.md) — SSR/SSG/Wails contexts; superseded by SSRContext-only for v0.1
- [SSR First (interactivity stack)](decisions/ssr-first.md) — Datastar/SSE and Tailwind as core options; superseded by HTMX/Alpine/plain JS + plugin boundary
- [Form Actions (validator dep)](decisions/form-actions.md) — go-playground/validator; superseded by built-in validators

## Concepts — Current

- [Dreego Architecture](concepts/dreego-architecture.md) — Architecture overview (SSR-only, dependency-free Core)
- [Dreego Sections](concepts/dreego-sections.md) — The 5 sections of a .dreego file
- [Template Logic](concepts/template-logic.md) — {#if}, {#each}, {#switch}, {#slot}
- [Plugin Ecosystem](concepts/plugin-ecosystem.md) — App-bound plugin registration
- [Form Actions](concepts/form-actions.md) — Form handling concept
- [Component Call Contract](concepts/component-call-contract.md) — Named props, slots, component contract (draft)

## Concepts — Research

- [Signals & Runes](concepts/signals-and-runes.md) — Signals and Svelte runes in Dreego (historical; narrowed for v0.1)
- [Output Strategy](concepts/output-strategy-comparison.md) — Comparison of output strategies
- [Form Actions Spec](concepts/form-actions-spec.md) — Form actions specification (v0.0.16)
- [Component Correctness](concepts/component-correctness.1.md) — Component correctness implementation plan (draft)

## Reference (KB) — Research

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
- [Progressive Enhancement](../_docs/progressive-enhancement.md) — HTMX, Alpine.js, plain JS guide
