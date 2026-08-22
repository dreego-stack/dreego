
---
type: Decision
title: Middleware System — Core vs Plugin
description: Three-class middleware: Core-Fixed, Core-Conditional, and Plugin middleware
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Middleware System — Core vs Plugin

## Decision

### Core-Fixed (always active, cannot be deactivated)

```go
// internal/middleware/core.go
func coreMiddleware() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{
        Recovery(),       // Panic → 500 + log
        RequestID(),      // X-Request-ID header
        RealIP(),         // Chi-built-in, behind proxy
        RequestLogging(), // slog, uses RequestID
    }
}
```

### Core-Conditional (on by default, opt-out via config)

```go
// Disableable via dreego.config.json
CSRF()        // Only needs Session interface (Core)
CORS()        // Default: same-origin restrictive
Compress()    // gzip/deflate via Chi
```

### Plugin (never Core)

| Middleware         | Why Plugin                                   |
|--------------------|----------------------------------------------|
| Rate-Limiting/DDoS | Needs Redis/backend → Infrastructure-dependent |
| Auth (OAuth/JWT)   | Policy, not everyone needs auth              |
| Prometheus/Metrics | Monitoring infrastructure                     |
| Tracing (OTel)     | External collector needed                     |

## Order Model

```
[Recovery → RequestID → RealIP → Logging]      ← Core-fixed
  → [CORS → CSRF → Compress]                    ← Core-conditionals
    → [Plugin middleware A, B, C…]              ← app.Use() FIFO
      → [Router / Handler]
```

**V1:** Pure registration order (FIFO via `app.Use()`)
**V2:** Constraint-based sorting (`before`/`after` in plugin manifest)

Middleware stack is locked at the first `ListenAndServe` — later `app.Use()` → Panic (deterministic).

## Middleware Signature

```go
type Middleware func(ctx Context, next http.Handler) http.Handler
```

Takes `dreego.Context`, not raw `*http.Request` — target-agnostic (applies to SSR, SSG, Wails).

## Consequences

- Core compiles WITHOUT external services — deterministic, testable
- CSRF is Core because it only needs Session (no external infrastructure)
- DDoS/Rate-Limiting is Plugin — needs Redis/backend
- Plugins register middleware via `MiddlewareProvider` in `app.Use()` order
