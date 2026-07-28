
---
type: Guide
title: Architecture Guide
description: Project structure, module boundaries, and architectural rules for Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Architecture Guide

## Project Structure

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go              # CLI entry point (max 120 lines)
├── dreego-core/                  # Core library (single package)
│   ├── lexer.go
│   ├── parser.go
│   ├── ast.go
│   ├── codegen.go
│   ├── router.go
│   ├── routes.go
│   ├── plugin.go
│   ├── context.go
│   └── middleware.go
├── dreego-plugin/                # Plugins (future)
├── internal/                     # Non-public packages
│   └── ...
├── testdata/                     # Test fixtures (.dreego files)
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── .kilo/
```

## Module Boundaries

| Module      | Responsibility                                    | Dependencies          |
|-------------|---------------------------------------------------|-----------------------|
| core        | .dreego → Go code, Router, Context, Middleware     | net/http, chi         |
| plugin      | Plugin interface, Registry (future)               | core                  |

## Rules

- **Every package is independently testable**
- **No circular dependencies**
- **`internal/` for implementation details not part of the public API**
- **`dreego-core/` for stable, public APIs**
- **`dreego-plugin/` for plugins (future)**
- **`cmd/` for entry points only**
