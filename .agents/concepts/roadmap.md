
---
type: Concept
title: "Dreego Roadmap"
description: "Development plan from V1 MVP through V2 Production-Ready to V3 Ecosystem"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Roadmap

## V1 — MVP (Core)

**Goal:** Working prototype — Transpiler + Router + Template

- [ ] **Transpiler Core**
  - Lexer: Parse `.dreego` file
  - Parser: Recognize `<go>`, Template, `<style>` sections
  - Code Generator: Output Go code
  - Template Variables: `{variable}` syntax
- [ ] **Template Logic (V1)**
  - `{#if}` / `{#else}`
  - `{#each}` / `{#else}`
- [ ] **File-based Routing**
  - Scan `routes/` directory
  - Generate `dreego_router.go`
  - Chi integration
- [ ] **CLI Tool**
  - `dreego generate` — Run transpiler
  - `dreego dev` — Dev server with auto-reload
  - `dreego build` — Build single binary
- [ ] **Single Binary**
  - `//go:embed` for static assets
  - Embed CSS/JS in binary

## V2 — Production-Ready

**Goal:** Full-fledged framework at Phoenix/Next.js level

- [ ] **Advanced Template Logic**
  - `{#switch}` / `{#case}`
  - `{#slot}` / `{#fill}`
  - `{#await}` (Async/SSE)
- [ ] **TypeScript Support**
  - esbuild integration as Go binding
  - Type sharing: Go struct → TS interface
- [ ] **SSG (Static Site Generation)**
  - `.dreego` → export static HTML/JS
  - For landing pages, blogs, docs
- [ ] **Wails Integration**
  - Use `.dreego` components in desktop apps
  - Same codebase for Web + Desktop
- [ ] **Dev Experience**
  - Hot Reload via SSE
  - Error overlay in browser
  - Tailwind JIT in dev server
- [ ] **Plugin System**
  - Stabilize `dreego.Plugin` interface
  - Transpiler hooks for custom tags
  - `dreego.config.json`

## V3 — Ecosystem

**Goal:** Establish addon ecosystem

- [ ] **Core Addons**
  - dreego-auth (MVP)
  - dreego-ui (MVP)
  - dreego-db (SQLite)
- [ ] **Documentation**
  - docs.dreego.dev
  - Tutorial project
  - API reference
- [ ] **Community**
  - Template Registry
  - Addon Registry

## Ideas for Later (Backlog)

- Auto-Form Binding (replace `r.FormValue`)
- Inline API Endpoints (`<go type="api">`)
- SSE Direct Directive (`<div g-sse="/api/live">`)
- Automatic Image Optimization (WebP/AVIF)
- View Transitions API Integration
- Global State & Event Bus (regeo.emit)
- Smart Error Boundaries
- Built-in CSRF & XSS Protection
