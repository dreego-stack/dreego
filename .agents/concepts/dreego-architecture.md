
---
type: Concept
title: "Dreego Architecture"
description: "Compile-time web framework architecture with transpiler, router, and plugin system"
tags: [v0.0.13]
timestamp: 2026-07-28T21:33:00Z
---
# Dreego Architecture

## Overview

Dreego is a compile-time web framework for Go. It consists of two main components:

1. **Dreego Transpiler** (`dreego generate`) — Converts `.dreego` files into Go code
2. **Dreego Runtime** (`import "github.com/.../dreego"`) — Provides router, context, plugin system

## Architecture Diagram

```
                              .dreego File
                                   │
          ┌────────────────────────┼────────────────────────┐
          ▼                        ▼                        ▼
      <head> Block           Template & <style>        <script> Block
  (Meta/Assets per Comp.)  (HTML + Scoped CSS)    (Vanilla JS for Client)
          │                        │                        │
          └────────────────────────┼────────────────────────┘
                                   ▼
                         Dreego Transpiler Engine
                                   │
                                   ▼
                          Generated .go Files
                                   │
                                   ▼
                            go build
                                   │
                                   ▼
                          Single Go Binary
                     (100% Type-Safe & SSR-First)
```

## Layered Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Dreego Addons                          │
│  dreego-auth, dreego-map, dreego-ui, dreego-admin, ...      │
├─────────────────────────────────────────────────────────┤
│                    Dreego Core                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Router   │ │ Context  │ │ Plugin   │ │ Middleware  │  │
│  │ (Chi)    │ │ (Req/Res)│ │ System   │ │ System      │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
├─────────────────────────────────────────────────────────┤
│              Dreego Transpiler (dreego generate)           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Lexer    │ │ Parser   │ │ AST      │ │ Code Gen   │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
├─────────────────────────────────────────────────────────┤
│              External Dependencies                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ chi      │ │ validator│ │ Tailwind │ │ HTMX/Alpine │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Data Flow (Request → Response)

```
Browser Request
      │
      ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Chi Router  │ ──▶ │ Middleware  │ ──▶ │ Dreego        │
│             │     │ (Auth, etc) │     │ Handler      │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                              ┌────────────────┼────────────────┐
                              ▼                ▼                ▼
                       ┌──────────┐    ┌──────────┐    ┌──────────┐
                       │ <go>     │    │ Template │    │ <head>   │
                       │ Data     │    │ Render   │    │ Assets   │
                       │ Fetching │    │ HTML     │    │ Injection│
                       └──────────┘    └──────────┘    └──────────┘
                              │                │                │
                              └────────────────┼────────────────┘
                                               ▼
                                    ┌─────────────────┐
                                    │ Final HTML       │
                                    │ (incl. <script>, │
                                    │  <style>, Assets)│
                                    └─────────────────┘
                                               │
                                               ▼
                                       Browser Response
```

## Project Structure (planned)

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go           # CLI: dreego generate, dreego dev, dreego build
├── dreego-core/               # Core library (single package)
├── dreego-plugin/             # Plugins (future)
├── go.mod
└── go.sum
```

## V1 Scope (MVP)

- [x] Transpiler: `.dreego` → `.go`
- [x] 3 Sections: `<go>`, Template, `<style>`
- [x] Template Logic: `{#if}`, `{#each}`, `{#else}`, `{#each else}`, `$loop`, `{#verbatim}`, `{var|raw|upper}`
- [x] File-based Routing (net/http 1.22+ enhanced routing)
- [x] Component System: `dreego/components/`, `<@Name>`, Named Slots, Scoped CSS
- [x] Static Assets: `dreego/static/` → inline handler + MIME types + collision check
- [x] Formatter: `dreego fmt` (v0.0.12)
- [x] Scaffolding: `dreego new` with landing blueprint (v0.0.13)
- [x] Split-Gen: `routes.go` + `components.go` + `dree.go` with `isUpToDate()` caching (v0.0.13)
- [x] Health Checks: `/health` + `/ready` endpoints (v0.0.14)
- [x] Security Headers: CSP, nosniff, frame-options (v0.0.14)
- [x] Gzip Compression: compress/gzip middleware (v0.0.14)
- [x] Content-Type Routing: `<go type="json|xml">` with c.JSON/XML/Bind/Write (v0.0.15)
- [x] Single Binary via `go build`
- [ ] Dev server with hot reload
- [ ] Plugin system (minimal)

## Components (v0.0.5+)

Reusable `.dreego` components with scoped styles:

```
dreego/components/Card.dreego:
  Component Card (title string)
  <div><h2>{title}</h2>{#slot}</div>

dreego/routes/get.dreego:
  <div><@Card title="Hello">content</@Card></div>
```

- Auto-Discovery via `scanComponents()` — no import needed
- `<@Name>` syntax with `@`-prefix to distinguish from HTML
- Self-closing: `<@Icon name="star"/>`
- Children → default `{#slot}`
- Named slots: `{#slot header}...{/slot}` via `c.Set`/`c.Get`
- Scoped CSS: `data-scope` per component (SHA256 hash)

## Static Assets (v0.0.10)

Files in `dreego/static/` are inlined into generated code and served automatically:

```
dreego/static/style.css  → GET /style.css  (text/css)
dreego/static/logo.svg   → GET /logo.svg   (image/svg+xml)
```

- MIME type via file extension
- Collision detection: if static path overlaps with route → `dreego generate` error
- Inline `[]byte` registration via `core.RegisterStatic()`

## Formatter (v0.0.12)

`dreego fmt` formats `.dreego` files in-place, `--check` for CI, `--stdout` for piping:

- Normalizes component headers, expressions, control flow spacing
- Section ordering: `<go>`, `<head>`, template, `<style>`, `<script>`
- Idempotent: formatting twice produces same output

## Scaffolding + Split-Gen (v0.0.13)

**Scaffolding:**
- `dreego new <name>` copies landing blueprint with `§$name$§` placeholders
- Auto-runs `go mod init` + `go mod edit` — project ready in seconds

**Split-Gen:**
- `dreego generate` now produces three files:
  - `gen/routes.go` — HTTP handler registration (one handler per route)
  - `gen/components.go` — component functions
  - `gen/dree.go` — config loading + static asset registration
- `isUpToDate()` per-file caching: file written only when content changes
- All in `gen` package for single-import compatibility

**Landing Blueprint:**
- Tailwind CDN, Hero + FeatureCard components, layout with `{#head}` + `{#slot}`
- Nav, pricing, CTA, footer sections
- Dockerfile (golang:1.22-alpine → distroless nonroot)
- `.gitignore` configured for Dreego project
