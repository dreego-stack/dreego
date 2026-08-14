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
├── _tests/                      # Integration tests
│   └── go/                      # Core/framework integration tests
├── _docs/                       # Public documentation
├── go.mod                       # Root module (core + cmd/dreego)
├── Makefile
└── Dockerfile
```

## Module Boundaries

| Module      | Responsibility                                    | Dependencies          |
|-------------|---------------------------------------------------|-----------------------|
| core        | .dreego → Go code, Router, Context, Middleware     | net/http (stdlib)     |
| external plugins | Optional capabilities in separate repositories | core + plugin-specific deps |

## Rules

- **Core never imports a plugin package** — plugins depend on Core, never the other way around
- **Every optional plugin gets a separate repository and `go.mod`**, including dependency-free plugins
- **Every package is independently testable**
- **No circular dependencies**
- **`core/` for stable, public APIs**
- **Separate repositories for official plugins**
- **`cmd/` for entry points only**
