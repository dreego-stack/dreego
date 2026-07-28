
---
type: Decision
title: Technology Stack for Dreego V1
description: Tech stack: Go, net/http, HTMX, Alpine.js
tags: [stack, v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Technology Stack for Dreego V1

**Date:** 2026-07-28
**Status:** Accepted

## Context

Dreego should build on proven Go libraries instead of reinventing everything. Each dependency was evaluated against the criterion "Build ourselves vs. use dependency."

## Decisions in Detail

### HTTP Router: go-chi/chi
- **Why not build ourselves?** Waste of time. Chi is ultra-fast, 100% Go stdlib compatible, extremely flexible, and well tested.
- **Alternative:** `net/http` directly (too low-level for file-based routing), gorilla/mux (no longer maintained)

### Template Engine: Dreego Custom Transpiler
- **Why build ourselves?** This is the heart of Dreego. No existing solution offers `.dreego` → Go code transpilation.
- **No alternative:** `a-h/templ` and `gomponents` are cool, but not what Dreego aims to be.

### Interactivity: HTMX + Alpine.js + Datastar
- **Why not build ourselves?** Building a JS framework is a huge project and not the goal.
- **Datastar** uses SSE for signals — ideal for Go's concurrency.
- **HTMX** and **Alpine.js** are extremely light and cover 95% of all interactivity cases.

### CSS: Tailwind CLI
- **Why not build ourselves?** Building a CSS parser/generator is unnecessarily complex.
- Dreego invokes the standalone Tailwind binary in the dev server.

### Validation: go-playground/validator
- **Why not build ourselves?** Too many edge cases. The validator is mature and covers all common cases.

### Binary Packaging: embed (Go Stdlib)
- **Why not build ourselves?** Already exists natively in Go. Packs Tailwind, JS, and templates directly into the binary.

## Rejected Dependencies

| Dependency      | Reason for Rejection                                  |
|-----------------|-------------------------------------------------------|
| Node.js / npm   | Destroys single-binary promise                        |
| esbuild (V1)    | In V1 only Vanilla JS — no bundler needed             |
| TypeScript      | Complexity, scope creep, first in V2                  |

## Consequences

- `go.mod` in V1: chi, validator, goree/dreego (self)
- Tailwind is embedded as a standalone binary (not as a Go dependency)
- HTMX + Alpine.js are shipped as embedded assets
