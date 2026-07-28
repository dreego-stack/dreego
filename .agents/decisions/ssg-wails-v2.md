
---
type: Decision
title: SSG & Wails Integration in V2
description: SSG and Wails integration as equal output modes alongside SSR in V2
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# SSG & Wails Integration in V2

**Date:** 2026-07-28
**Status:** Accepted (planned for V2)

## Context

Dreego V1 is SSR-only (Server-Side Rendering). For many use cases that is sufficient. But two important scenarios need static output:

1. **SSG (Static Site Generation):** Compile `.dreego` pages to static HTML/JS/CSS
2. **Wails:** Use `.dreego` components in Go desktop apps
3. **Mobile (future):** Same components via Wails Mobile / Gomobile

## Decision

**V1: SSR-Only** (no SSG). The focus is on the transpiler + router.
**V2: SSG + Wails** as equal output modes alongside SSR.

## SSG Use Cases

| Use Case               | Example                                  |
|------------------------|-------------------------------------------|
| Cloudflare Pages       | Deploy static HTML files on Edge         |
| GitHub Pages           | Project docs, landing pages               |
| S3/Cloudflare R2       | Pure static sites, 0 server costs         |
| Blog                   | Markdown → .dreego → static HTML         |
| Documentation          | docs.dreego.dev itself built with SSG     |

## Wails Use Cases

| Use Case               | Example                                  |
|------------------------|-------------------------------------------|
| Desktop App            | Go backend + Dreego frontend = Native App |
| Tray App               | Menu bar app with Dreego UI               |
| Mobile (future)        | Same codebase for iOS/Android             |

## Advantage: Code Reuse

```
.dreego Components
       │
       ├── SSR (Web)          — Chi server, HTML on-the-fly
       ├── SSG (Static)       — Static HTML files
       ├── Wails (Desktop)    — Native windows, system APIs
       └── Mobile (later)     — iOS/Android via Wails Mobile
```

The same `.dreego` file renders in four different contexts.
No JS framework can do that — because they all need a JS runtime.

## Architecture Preparation in V1

Even though SSG/Wails come only in V2, the architecture must be prepared in V1:

1. **Transpiler pipeline with output modes:** The code generator has a `Target` interface
   - `TargetSSR` — Go HTTP handler (V1)
   - `TargetSSG` — Static HTML files (V2)
   - `TargetWails` — Wails-compatible Go functions (V2)
2. **`dreego build --static`** — CLI flag already reserved in V1 (does nothing, shows "Coming in V2")
3. **No SSR-specific assumptions in the template:** `<go>` block can only have server code in V1, but the template is target-agnostic

## Inspiration from the Rust World

- **Dioxus:** Same components for Web, Desktop (Blitz), Mobile
- **Leptos:** SSR + Hydration + Islands — shows that multi-target architecture works
- **Yew:** Was CSR-only, added SSR later — harder than designing multi-target from the start

## Consequences

- V1: `dreego build` generates SSR binary (HTTP server)
- V2: `dreego build --static` generates `dist/` with HTML files
- V2: `dreego build --wails` generates Wails-compatible code
- `dreego.config.json` gets a `target` field: `"ssr" | "ssg" | "wails"`
- Same `.dreego` components usable across all targets
