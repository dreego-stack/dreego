
---
type: Decision
title: SSR-First Architecture
description: Dreego renders HTML on the server with HTMX and Alpine.js for interactivity
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# SSR-First Architecture

**Date:** 2026-07-28
**Status:** Accepted — SSR-first remains the v0.1 contract; see note below

> **Historical note:** The interactivity stack is narrowed for v0.1. HTMX,
> Alpine.js, and plain JavaScript remain the supported progressive-enhancement
> path. Datastar/SSE-based real-time updates are no longer a core option — SSE
> and WebSockets are plugins, not core packages (AGENTS.md "Core and Plugin
> Boundary"). Tailwind is not a fixed core dependency; it is a plugin. There is
> no internal client runtime before v0.1.

## Context

Dreego should be a web framework. The question: Client-Side Rendering (CSR/SPA) or Server-Side Rendering (SSR)?

## Decision

**SSR-First.** Dreego renders everything on the server. The client receives finished HTML.

Interactivity comes not through a client-side framework, but through:
- **HTMX** for server interactions (partial page updates without reload)
- **Alpine.js** for local UI interactions (dropdowns, tabs, modals)
- **Datastar** (optional) for SSE-based real-time updates

## Rationale

1. **No JS build hell:** No node_modules, no Webpack, no Vite
2. **0 MB state synchronization:** State exists only on the Go server
3. **Perfect SEO:** Everything is static HTML on first load
4. **Fast FCP (First Contentful Paint):** Even on weak mobile devices
5. **Simpler architecture:** No API layer, no JSON serialization, no client state stores
6. **Direct DB access:** `<go>` block can directly access the database

## Comparison

| Aspect              | CSR (React/Svelte)         | SSR (Dreego)                  |
|---------------------|----------------------------|-------------------------------|
| Initial Load        | JS must load + hydrate     | Immediate finished HTML       |
| SEO                 | Difficult (SSR needed)     | Perfect                       |
| State Management    | Client + Server sync       | Server only                   |
| Bundle Size         | 100+ KB JS                 | ~10 KB (HTMX + Alpine)        |
| Deployment          | Node.js + Static Files     | Single Go Binary              |

## Counterarguments

- **"SSR feels slow on interactions"** → HTMX only swaps HTML fragments, no full page reload
- **"No SPA feeling"** → Alpine.js + View Transitions API for smooth transitions
- **"Less interactive"** → Datastar streams DOM updates via SSE (like Phoenix LiveView)

## Consequences

- No client-side router needed
- No API layer between template and database
- HTMX + Alpine.js are the supported progressive-enhancement path; plain JavaScript is always available
- Historical: V2: SSG (Static Site Generation) for purely static pages (deferred past v1; see [ssg-wails-v2](ssg-wails-v2.md))
