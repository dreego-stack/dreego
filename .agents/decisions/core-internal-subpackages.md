---
type: Decision
title: Core runtime split into internal subpackages
description: core/internal/ holds session, server, middleware, context, validate; core/ re-exports the public API
tags: [architecture]
timestamp: 2026-08-20T00:00:00Z
---
# Core runtime split into internal subpackages

**Date:** 2026-08-20
**Status:** Accepted

## Context

After the transpiler moved to `internal/transpiler/`, the runtime framework
stayed as one flat `core` package. Everything shared unexported symbols, so
nothing could be imported from outside the package, but the package had no
structural boundary enforcement: session crypto, middleware internals, and
the App state machine all lived side by side. `core/app.go` exceeded the
300-line limit.

The user explicitly preferred internal subpackages over public ones
(`core/session`, `core/server`, ...): public subpackages would have doubled
the public API contract surface before v0.1. The user also wanted the
transpiler to remain at `internal/transpiler/` (not under the CLI), which
stays unchanged.

## Decision

Split the runtime into internal subpackages, keeping a thin public facade
at `core/`:

```
core/facade.go            ← public API: type aliases + wrapper funcs
core/internal/server/     ← App, routing, redirects, rewrites, server config
core/internal/session/    ← Store interface, CookieStore, crypto, policy
core/internal/middleware/ ← Recovery, CSRF, Compress, RequestID, logging, body limit, security headers
core/internal/context/    ← SSRContext, Context, JSON/XML/Bind/Write
core/internal/validate/   ← BindForm, ValidateForm, SaveOld, SaveErrors
```

The facade preserves the single import path for applications and generated
code: `dreego "github.com/dreego-stack/dreego/core"` keeps working without
changes. Only codegen-emitted symbols (`SSRContext`, `Safe*`, ...) are
re-exported; internal packages are not importable from outside the module.

## Dependency Direction

- `core/internal/context` → `core/internal/session`, `core/internal/middleware`
- `core/internal/middleware` → `core/internal/session` (Store, IsTLS)
- `core/internal/server` → `core/internal/session`, `core/internal/middleware`, `core/internal/validate`
- `core/internal/validate` → (none)
- `core/` → all internal packages

No cycles.

## Consequences

- The public API is unchanged: type aliases preserve methods and receivers;
  wrapper functions preserve signatures. Generated code and `_tests/go`
  integration tests compile unmodified.
- Unexported helpers (`responseWriter`, `jsonlHandler`, `gzipBuffer`,
  `isTLS`, `deriveKeys`, `encryptPayload`) are no longer accessible outside
  their package; a few had to be exported for sibling-package use
  (`DefaultCSP`, `ValidatorFunc`, `Atoi`, `MaxCookieSize`, `IsTLS`,
  `RequestIDKey`).
- Mixed test files were split at the package boundary; tests using unexported
  internals moved into the internal packages, facade tests stayed at `core/`.
- CI (`_tests/test.sh`, `pull-request-check.yml`) already runs `./core/...`,
  which covers the internal packages; `_scripts/check-core-deps.sh` now checks
  `./core/...` too.
- The 300-line rule is now enforced per file inside the split.
