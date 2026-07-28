
---
type: Decision
title: Compile-Time Transpiler Instead of Runtime Parsing
description: Build-time code generation, no runtime parsing
tags: [transpiler, v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Compile-Time Transpiler Instead of Runtime Parsing

**Date:** 2026-07-28
**Status:** Accepted

## Context

`.dreego` files must be converted to executable code. Two approaches are available:

1. **Runtime parsing:** Server reads `.dreego` files at runtime (similar to `html/template`)
2. **Compile-time transpiler:** `dreego generate` converts `.dreego` → `.go` before the build

## Decision

**Compile-Time Transpiler** (Path A from the Gemini chat).

`dreego generate` reads `.dreego` files and generates Go code from them.

## Rationale

| Criterion              | Runtime Parsing      | Compile-Time (chosen)        |
|------------------------|----------------------|------------------------------|
| Performance            | Slower (parsing)     | Maximum (no runtime overhead) |
| Single Binary          | via `//go:embed`     | Everything in binary, no parsing |
| Error detection        | At runtime (crash)   | At build time (`go build` fails) |
| Type safety            | None                 | Full Go type safety          |
| DevX                   | No build step        | `dreego generate` in watcher  |
| Debugging              | Difficult            | Normal (generated Go code)   |

## Consequences

- Build step: `dreego generate` must run before `go build`
- Generated `*_dreego.go` files are not committed
- 100% compile-time safety: No template error reaches production
- Dev server automatically runs `dreego generate` on file changes
- Go has no macros — code generation is the Go-idiomatic way
