---
type: Log
title: Knowledge Base Changelog
description: Record of all changes to the .agents/ knowledge bundle
tags: [log]
timestamp: 2026-07-25T00:00:00Z
---

# log

## 2026-07-25

- Context refactoring: `map[string]string` → Interface + Embedding
  - New `Context` interface embedding `context.Context` with `Param`, `Data`, `Render`
  - New `SSRContext` struct with `NewSSR()`, `Set`, `Get`, `Query`, `FormValue`
  - Updated codegen to generate `*context.SSRContext` + `context.NewSSR(w, r)`
  - Deleted stale `demo/dreego/layouts/dree.go` (layout is now inlined per route)
- Recovery-Middleware: Panic → 500 + Stack-Trace-Logging via slog
  - New `pkg/middleware/recovery.go` — defer recover() with JSON log
  - Integrated as outermost Core-Fixed middleware in runtime pipeline
- Converted entire knowledge base to Open Knowledge Format (OKF) v0.1
- Added YAML frontmatter with `type` field to all files
- Replaced `[[wiki-links]]` with standard markdown links
- Renamed `_index.md` to `index.md` with OKF child-list format
- Moved `thinking-list.md` to `PlanTODO.md`
- Moved `gemini-50-tipps.md` to `tips.md`
- Added `open-knowledge-format.md` skill guide
- Restructured `ROADMAP.md` as release pipeline
- Split `TODO.md` into immediate actions only

## 2026-07-24

- Decision: Error-Handling with typed errors, recovery middleware
- Decision: Routing & Components — hybrid file-based + programmatic routing
- Decision: Output strategy comparison (per-dir dree.go confirmed by GLM)
- Implemented gen/dree.go central import file
- Created dreego/config.json with redirects, rewrites, logging config
- Added RequestLogging middleware (JSONL, Core-Conditional)
- CLI commands: generate, build, run (-d, -t)

## 2026-07-23

- Initial knowledge base created
- Decision: Name "dreego", file extension ".dreego"
- Decision: Technology stack (Go 1.22+, net/http, HTMX, Alpine.js)
- Decision: Compile-Time Transpiler approach
- Decision: Transpiler Pipeline (Lexer → Parser → AST → CodeGen)
- Decision: TypeScript deferred to V2
- Decision: 5 sections in .dreego files
- Decision: SSR-First architecture
- Decision: No catch tag — errors via Go idioms
- Decision: File-based routing
- Decision: SSG & Wails in V2
- Researched: React/Svelte ecosystems, Phoenix/Laravel/Django, Rust frameworks, Solid/Astro, Blazor
- Concepts drafted: Architecture, Sections, Template Logic, Addon Ecosystem
