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
├── core/                        # Public framework API (facade re-exports, stdlib only)
│   └── internal/               # Runtime implementation (not importable from outside the module)
│       ├── server/             # App, routing, redirects, rewrites, server config
│       ├── session/            # Store interface, CookieStore, crypto, policy
│       ├── middleware/         # Recovery, CSRF, Compress, RequestID, logging, body limit, security
│       ├── context/            # SSRContext, Context, JSON/XML/Bind/Write
│       └── validate/           # BindForm, ValidateForm, SaveOld, SaveErrors
├── internal/transpiler/       # .dreego → Go code (used by CLI and dreegotest)
├── cli/dreego/                # CLI entry point
├── dreegotest/                # Test helpers for integration tests
├── _tests/                    # Integration tests
├── _docs/                     # Public documentation
├── go.mod                     # Root module (one tag per release)
├── Makefile
└── _scripts/                  # check-core-deps.sh, release prep
```

## Module Boundaries

| Module            | Responsibility                                    | Dependencies          |
|-------------------|---------------------------------------------------|-----------------------|
| core (facade)     | Public API: App, SSRContext, middleware, sessions, Safe* | stdlib only        |
| core/internal/*   | Runtime implementation                              | stdlib only           |
| internal/transpiler | .dreego → Go code generation                     | stdlib only           |
| external plugins  | Optional capabilities in separate repositories     | core + plugin-specific deps |

## Rules

- **Core never imports a plugin package** — plugins depend on Core, never the other way around
- **Every optional plugin gets a separate repository and `go.mod`**, including dependency-free plugins
- **Every package is independently testable**
- **No circular dependencies**: context → session/middleware → server; validate and safety are leaves
- **`core/` is the only public import path for applications** — `core/internal/*` is unimportable outside the module, `internal/transpiler` likewise
- **Separate repositories for official plugins**
- **`cli/` for entry points only**
- **`core/internal/*` may not import `core/` (facade)** — the facade imports them
