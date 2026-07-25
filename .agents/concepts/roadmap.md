
---
type: Concept
title: "Dreego Roadmap"
description: "Entwicklungsplan von V1 MVP über V2 Production-Ready bis V3 Ecosystem"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Dreego Roadmap

## V1 — MVP (Core)

**Ziel:** Lauffähiger Prototyp — Transpiler + Router + Template

- [ ] **Transpiler Core**
  - Lexer: `.dreego`-Datei parsen
  - Parser: `<go>`, Template, `<style>` Sektionen erkennen
  - Code Generator: Go-Code ausgeben
  - Template-Variablen: `{variable}` Syntax
- [ ] **Template-Logik (V1)**
  - `{#if}` / `{#else}`
  - `{#each}` / `{#else}`
- [ ] **File-based Routing**
  - `routes/`-Ordner scannen
  - `dreego_router.go` generieren
  - Chi-Integration
- [ ] **CLI-Tool**
  - `dreego generate` — Transpiler ausführen
  - `dreego dev` — Dev-Server mit Auto-Reload
  - `dreego build` — Single Binary bauen
- [ ] **Single Binary**
  - `//go:embed` für statische Assets
  - CSS/JS in Binary einbetten

## V2 — Production-Ready

**Ziel:** Vollwertiges Framework auf Phoenix/Next.js-Niveau

- [ ] **Erweiterte Template-Logik**
  - `{#switch}` / `{#case}`
  - `{#slot}` / `{#fill}`
  - `{#await}` (Async/SSE)
- [ ] **TypeScript Support**
  - esbuild-Integration als Go-Binding
  - Types-Sharing: Go-Struct → TS-Interface
- [ ] **SSG (Static Site Generation)**
  - `.dreego` → statisches HTML/JS exportieren
  - Für Landing Pages, Blogs, Docs
- [ ] **Wails Integration**
  - `.dreego`-Komponenten in Desktop-Apps verwenden
  - Gleiche Codebase für Web + Desktop
- [ ] **Dev-Experience**
  - Hot Reload via SSE
  - Error Overlay im Browser
  - Tailwind JIT im Dev-Server
- [ ] **Plugin-System**
  - `dreego.Plugin` Interface stabilisieren
  - Transpiler-Hooks für Custom-Tags
  - `dreego.config.json`

## V3 — Ecosystem

**Ziel:** Addon-Ökosystem etablieren

- [ ] **Core Addons**
  - dreego-auth (MVP)
  - dreego-ui (MVP)
  - dreego-db (SQLite)
- [ ] **Dokumentation**
  - docs.dreego.dev
  - Tutorial-Projekt
  - API-Referenz
- [ ] **Community**
  - Template Registry
  - Addon Registry

## Ideen für später (Backlog)

- Auto-Form Binding (`r.FormValue` ersetzen)
- Inline API Endpoints (`<go type="api">`)
- SSE Direct Directive (`<div g-sse="/api/live">`)
- Automatic Image Optimization (WebP/AVIF)
- View Transitions API Integration
- Global State & Event Bus (regeo.emit)
- Smart Error Boundaries
- Built-in CSRF & XSS Protection
