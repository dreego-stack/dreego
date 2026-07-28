
---
type: Concept
title: "Dreego Architecture"
description: "Compile-Time Webframework-Architektur mit Transpiler, Router und Plugin-System"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Architecture

## Overview

Dreego ist ein Compile-Time Webframework für Go. Es besteht aus zwei Hauptkomponenten:

1. **Dreego Transpiler** (`dreego generate`) — Wandelt `.dreego`-Dateien in Go-Code um
2. **Dreego Runtime** (`import "github.com/.../dreego"`) — Bietet Router, Context, Plugin-System

## Architektur-Diagramm

```
                              .dreego Datei
                                   │
          ┌────────────────────────┼────────────────────────┐
          ▼                        ▼                        ▼
      <head> Block           Template & <style>        <script> Block
  (Meta/Assets pro Comp.)   (HTML + Scoped CSS)    (Vanilla JS für Client)
          │                        │                        │
          └────────────────────────┼────────────────────────┘
                                   ▼
                         Dreego Transpiler Engine
                                   │
                                   ▼
                          Generierte .go Dateien
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
│              Externe Dependencies                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────────┐  │
│  │ chi      │ │ validator│ │ Tailwind │ │ HTMX/Alpine │  │
│  └──────────┘ └──────────┘ └──────────┘ └────────────┘  │
└─────────────────────────────────────────────────────────┘
```

## Datenfluss (Request → Response)

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
                       │ Fetching │    │ HTML     │    │ Injektion│
                       └──────────┘    └──────────┘    └──────────┘
                              │                │                │
                              └────────────────┼────────────────┘
                                               ▼
                                    ┌─────────────────┐
                                    │ Finales HTML     │
                                    │ (inkl. <script>, │
                                    │  <style>, Assets)│
                                    └─────────────────┘
                                               │
                                               ▼
                                       Browser Response
```

## Projektstruktur (geplant)

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
- [x] 3 Sektionen: `<go>`, Template, `<style>`
- [x] Template-Logik: `{#if}`, `{#each}`, `{#else}`, `{#each else}`, `$loop`, `{#verbatim}`, `{var|raw|upper}`
- [x] File-based Routing (net/http 1.22+ enhanced routing)
- [x] Component-System: `dreego/components/`, `<@Name>`, Named Slots, Scoped CSS
- [x] Static Assets: `dreego/static/` → inline Handler + MIME-Types + Collision-Check
- [x] Single Binary via `go build`
- [ ] Dev-Server mit Hot Reload
- [ ] Plugin-System (minimal)

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

- MIME-Type via file extension
- Collision detection: if static path overlaps with route → `dreego generate` error
- Inline `[]byte` registration via `core.RegisterStatic()`
