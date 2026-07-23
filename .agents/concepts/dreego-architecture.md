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
├── pkg/
│   ├── transpiler/            # Lexer, Parser, AST, Code-Generator
│   ├── router/                # Chi-Wrapper, File-based Routing
│   ├── plugin/                # Plugin-Interface & Registry
│   ├── context/               # Request-Context, Session, User
│   ├── middleware/             # CSRF, Session, Auth
│   └── html/                  # HTML-Builder (optional, für Pure-Go-Weg B)
├── go.mod
└── go.sum
```

## V1 Scope (MVP)

- [x] Transpiler: `.dreego` → `.go`
- [x] 3 Sektionen: `<go>`, Template, `<style>`
- [x] Template-Logik: `{#if}`, `{#each}`
- [x] Chi-Router Integration
- [x] File-based Routing
- [x] Single Binary via `//go:embed`
- [ ] Dev-Server mit Hot Reload
- [ ] Tailwind CLI Integration
- [ ] Plugin-System (minimal)
