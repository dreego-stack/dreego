
---
type: Concept
title: "Dreego Architecture"
description: "Compile-time web framework architecture with transpiler, router, and plugin system"
tags: [v0.0.21]
timestamp: 2026-07-31T07:00:00Z
---
# Dreego Architecture

> **Current and planned boundaries:** This document primarily explains the
> released SSR implementation and therefore uses the current `<server>`, `<body>`,
> and `<client>` syntax. The accepted v0.x direction introduces a target-neutral
> App and renderer, explicit SSR/SSG/Wails target packages, DreeJS, and the
> planned `<server>`, `<body>`, and `<client>` names. See the
> [target decision](decisions/target-neutral-application-and-first-party-targets.md),
> [section decision](decisions/semantic-sections-and-language-processors.md),
> and [`_plan/`](../_plan/README.md).

> **Historical note:** The layered architecture is current. The monorepo plugin
> layout (`plugins/` + `go.work`, v0.0.21) is superseded — official plugins use
> separate repositories and modules. Tailwind is not a core dependency; it is a
> plugin. Core is dependency-free (stdlib only). See AGENTS.md "Core and Plugin
> Boundary".

## Overview

Dreego is a compile-time web framework for Go. It consists of two main components:

1. **Dreego Transpiler** (`dreego generate`) — Converts `.dreego` files into Go code
2. **Dreego Runtime** (`import "github.com/dreego-stack/dreego/core"`) — Provides router, context, plugin system

## Architecture Diagram

```
                              .dreego File
                                   │
          ┌────────────────────────┼────────────────────────┐
          ▼                        ▼                        ▼
      <head> Block           Template & <style>        <client> Block
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
│                    Dreego Plugins                          │
│  plugins/auth, plugins/map, plugins/ui, plugins/admin, ... │
├─────────────────────────────────────────────────────────┤
│                    Dreego Core                            │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Router   │ │ Context  │ │ Plugin   │ │ Middleware  │  │
│  │(net/http)│ │ (Req/Res)│ │ System   │ │ System      │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
├─────────────────────────────────────────────────────────┤
│              Dreego Transpiler (dreego generate)           │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ Lexer    │ │ Parser   │ │ AST      │ │ Code Gen   │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
├─────────────────────────────────────────────────────────┤
│              External Dependencies                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ net/http │ │ (stdlib) │ │ Tailwind │ │ HTMX/Alpine │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Data Flow (Request → Response)

```
Browser Request
      │
      ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ net/http    │ ──▶ │ Middleware  │ ──▶ │ Dreego        │
│ ServeMux    │     │ (Auth, etc) │     │ Handler      │
└─────────────┘     └─────────────┘     └──────┬──────┘
                                               │
                              ┌────────────────┼────────────────┐
                              ▼                ▼                ▼
                       ┌──────────┐    ┌──────────┐    ┌──────────┐
                       │ <server>     │    │ Template │    │ <head>   │
                       │ Data     │    │ Render   │    │ Assets   │
                       │ Fetching │    │ HTML     │    │ Injection│
                       └──────────┘    └──────────┘    └──────────┘
                              │                │                │
                              └────────────────┼────────────────┘
                                               ▼
                                    ┌─────────────────┐
                                    │ Final HTML       │
                                    │ (incl. <client>, │
                                    │  <style>, Assets)│
                                    └─────────────────┘
                                               │
                                               ▼
                                       Browser Response
```

## Project Structure

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go           # CLI: dreego generate, dreego dev, dreego build
├── core/                      # Core library (single package, stdlib only)
├── plugins/                   # No optional implementations; official plugins use separate repos
├── _tests/                    # Integration tests
│   └── core/<Category>/
├── go.mod
├── go.work
└── go.sum
```

## V1 Scope (MVP)

- [x] Transpiler: `.dreego` → `.go`
- [x] 3 Sections: `<server>`, Template, `<style>`
- [x] Template Logic: `{#if}`, `{#each}`, `{#else}`, `{#else if}`, `{#each else}`, `$loop`, `{#verbatim}`, and `{{ expression }}`
- [x] File-based Routing (net/http 1.22+ enhanced routing)
- [x] Component System: `dreego/components/`, `<@Name>`, Named Slots, Scoped CSS
- [x] Static Assets: `dreego/static/` → inline handler + MIME types + collision check
- [x] Formatter: `dreego fmt` (v0.0.12)
- [x] Scaffolding: `dreego new` with landing blueprint (v0.0.13)
- [x] Split-Gen: `routes.go` + `components.go` + `dree.go` with `isUpToDate()` caching (v0.0.13)
- [x] Health Checks: `/health` + `/ready` endpoints (v0.0.14)
- [x] Security Headers: nosniff, frame-options, referrer-policy, permissions-policy (v0.0.14)
- [x] CSP header via `core.SetCSP` (v0.0.20)
- [x] Gzip Compression: compress/gzip middleware (v0.0.14)
- [x] Content-Type Routing: `<server type="json|xml">` with c.JSON/XML/Bind/Write (v0.0.15)
- [x] Form Actions: `g-action` + auto-validation + redirect (v0.0.16)
- [x] Production Deployment: graceful shutdown, cross-compile, Docker (v0.0.17)
- [x] Request-ID Middleware (v0.0.17)
- [x] CSRF cookie Secure flag TLS-aware (v0.0.20)
- [x] Monorepo plugin layout: `plugins/` + `go.work` (v0.0.21)
- [x] Single Binary via `go build`
- [ ] Dev server with hot reload
- [ ] Plugin interface validated by real external plugins before the v1 stability promise

## Components (v0.0.5+)

Reusable `.dreego` components with scoped styles:

```
dreego/components/Card.dreego:
  Component Card (title string)
  <body><h2>{{ title }}</h2>{#slot}</body>

dreego/routes/page.dreego:
  <body><@Card title="Hello">content</@Card></body>
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
- Section ordering: `<server>`, `<head>`, `<body>`, `<style>`, `<client>`
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
