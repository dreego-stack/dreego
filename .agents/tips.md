
---
type: Reference
title: 50 Tips + Checklist
description: 50 development tips and checklist covering DX, architecture, templating, plugins, and performance
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# 50 Tips + Checklist

**Source:** Gemini Chat, 2026-07-25

---

## 1. Developer Experience (DX) & Ergonomics (1–10)

| # | Tip | Status |
|---|------|--------|
| 1 | Errors with `.dreego` line numbers, never from `dree.go` | planned |
| 2 | Colored CLI output (green checkmarks, red errors) | planned |
| 3 | `dreego init my-app` — working minimal project | planned |
| 4 | Allow flat folder structure | ✅ `dreego/routes/get.dreego` is enough |
| 5 | Color-highlighted HTTP logs in dev mode | planned (currently: JSONL) |
| 6 | Serve `static/` folder via `embed.FS` | planned |
| 7 | Compact CLI output, only show changes | planned |
| 8 | Hot reload in <1s | planned (V2) |
| 9 | `dreego build` → single binary | ✅ |
| 10 | `/health` endpoint activatable via config | planned |

## 2. Architecture Design & Routing (11–20)

| # | Tip | Status |
|---|------|--------|
| 11 | Explicit HTTP method suffixes | ✅ `get.dreego`, `post.dreego` |
| 12 | Type-safe path parameters (int/string) | planned |
| 13 | Wildcard routes `[...all].dreego` | planned (decision made) |
| 14 | `_middleware.go` per folder | planned (decision made) |
| 15 | No global state, avoid race conditions | refactoring (runtime globals) |
| 16 | Pass through `context.Context` | ✅ `c.R.Context()` |
| 17 | Custom `404.dreego` and `500.dreego` | planned |
| 18 | Form data parser (`x-www-form-urlencoded`) | planned |
| 19 | SSE / Streaming from `<go>` block | planned |
| 20 | `dreego.Redirect(ctx, "/login")` | planned |

## 3. Templating & Component System (21–30)

| # | Tip | Status |
|---|------|--------|
| 21 | Clear separation: `<go>`, `<script>`, `<style>`, HTML | ✅ |
| 22 | Type-safe props for components | planned (decision made) |
| 23 | Automatic XSS escaping | planned |
| 24 | Raw HTML escape hatch | planned |
| 25 | CSS scope hash for `<style>` block | ✅ `data-scope="hash"` |
| 26 | JS inlining from `<script>` | ✅ |
| 27 | Slot system for layouts | ✅ `{#slot}` |
| 28 | Conditional rendering `{#if}`, `{#each}` | ✅ |
| 29 | Zero-JS mode (no `<script>` section) | ✅ implicitly (nothing = no JS) |
| 30 | Head management (`<title>`, `<meta>`) | ✅ `<head>` block |

## 4. Plugins, Type Safety & Ecosystem (31–40)

| # | Tip | Status |
|---|------|--------|
| 31 | Compile-time plugin pattern (Go packages) | ✅ `dreego.Plugin` Interface |
| 32 | Interface-driven design | ✅ Capability-based |
| 33 | Type-safe context values | planned (decision made) |
| 34 | Validate plugin configuration | planned |
| 35 | Official auth plugin | planned (dreego-auth) |
| 36 | i18n plugin with type safety | planned (V2) |
| 37 | Zero CGO for cross-compiling | ✅ (no CGO dependencies) |
| 38 | CLI plugin hooks | planned |
| 39 | Plugin lifecycle (`OnStart`, `OnShutdown`) | planned (decision made) |
| 40 | Plugin documentation template | planned |

## 5. Performance, Security & Accessibility (41–50)

| # | Tip | Status |
|---|------|--------|
| 41 | Secure HTTP headers (CSP, X-Frame-Options) | ✅ CSP v0.0.20, headers v0.0.14 |
| 42 | Gzip/Brotli compression | ✅ Gzip v0.0.14 |
| 43 | Zero-allocation in hot path | planned (V2) |
| 44 | CSRF protection | ✅ v0.0.3, Secure flag v0.0.20 |
| 45 | Graceful shutdown | ✅ v0.0.17 |
| 46 | VS Code Syntax Highlighting Extension | planned |
| 47 | ARIA warnings in CLI | planned (V3) |
| 48 | Go doc comments for core | planned |
| 49 | Dockerfile template (multi-stage on scratch) | planned |
| 50 | Showcase projects (Todo, Blog) | planned |

## Summary

- ✅ Already implemented: 11 of 50
- Planned (decision made): 22 of 50
- Planned (V2/V3): 17 of 50
