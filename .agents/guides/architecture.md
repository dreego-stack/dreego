---
type: Guide
title: Architecture Guide
description: Project structure, module boundaries, and architectural rules for Dreego
tags: [v0.0.21]
timestamp: 2026-07-31T07:00:00Z
---
# Architecture Guide

## Project Structure

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go              # CLI entry point
├── core/                        # Core library (single package, stdlib only)
│   ├── lexer.go
│   ├── parser.go
│   ├── ast.go
│   ├── codegen.go
│   ├── router.go
│   ├── context.go
│   └── middleware.go
├── plugins/                     # Official plugins (each with own go.mod if deps needed)
│   └── sample/
├── _tests/                      # Integration tests
│   ├── core/<Category>/         # Core/framework tests
│   └── plugins/<name>/          # Plugin tests
├── _docs/                       # Public documentation
├── go.mod                       # Root module (core + cmd/dreego)
├── go.work                      # Links root module + plugin modules
├── Makefile
└── Dockerfile
```

## Module Boundaries

| Module      | Responsibility                                    | Dependencies          |
|-------------|---------------------------------------------------|-----------------------|
| core        | .dreego → Go code, Router, Context, Middleware     | net/http (stdlib)     |
| plugins/*   | Official plugins (auth, db, cache, etc.)          | core + plugin-specific deps |

## Rules

- **Core never imports a plugin package** — plugins depend on Core, never the other way around
- **Plugins with external dependencies get their own `go.mod`** inside `plugins/<name>/`
- **Every package is independently testable**
- **No circular dependencies**
- **`core/` for stable, public APIs**
- **`plugins/` for official plugins**
- **`cmd/` for entry points only**
- **`go.work` links all local modules for development**