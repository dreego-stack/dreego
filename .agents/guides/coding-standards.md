
---
type: Guide
title: Coding Standards for Dreego
description: Code conventions: file size limits, imports, error handling, testing for Dreego
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Coding Standards for Dreego

## General

- **Max 120 lines per file** — hard limit
- **One logical thing per file**
- **No comments** — code speaks for itself
- **Package names** short, clean, without hyphens
- **Go 1.26+**, prefer standard library

## Build & Run

- **Never `go build` directly** — Build via `make build` or `dreego build`
- Dev server: `dreego dev`
- Tests: `make test`
- Generated files (`*_dreego.go`) are not committed

## Imports

Standard library first, then external, then internal. Group with blank lines:

```go
import (
    "context"
    "log"

    "github.com/go-chi/chi/v5"

    core "codeberg.org/dreego/dreego/dreego-core"
)
```

## Error Handling

- Always explicit: `if err != nil { return err }`
- No `panic()` except in `init()` and tests
- `fmt.Errorf` with `%w` for wrapping

## Tests

- Test files next to the code being tested (`foo_test.go`)
- Prefer table-driven tests
- Test fixtures in `testdata/`
- On failing test: fix code, not the test
