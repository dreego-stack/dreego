
---
type: Decision
title: Error Handling Strategy (Core, not Plugin)
description: Typed error types and recovery middleware for build-start-runtime errors
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Error Handling Strategy (Core, not Plugin)

**Date:** 2026-07-28
**Status:** Accepted

## Context

Dreego is a compile-time transpiler. Errors can occur at three levels:

1. **Build time** (`dreego generate`, `go build`) — Template syntax, missing variables, type errors
2. **Start time** (`main()`) — Missing config, invalid plugins, port conflict
3. **Runtime** (per request) — Validation, DB errors, panics, auth errors

The framework must be type-safe. Errors should be as early as possible (Build > Start > Runtime) and as local as possible (per handler, not global).

`net/http` only offers `http.Error(w, msg, code)`. There is no recovery, no structured logging, no typed errors, no flash/redirect semantics.

## Decision

**Typed error types in the core.** Recovery middleware dispatches on type, renderer decides HTML/JSON/Redirect.

### Error Types (new file `dreego-core/errors.go`)

```go
package errors

type HTTPError struct {
    Code   int
    Public string  // visible to end user
    Cause  error   // internal error (Dev mode: Stack, Prod: only log)
}

type ValidationError struct {
    Fields map[string]string  // Field name → Error message
}

type RedirectError struct {
    URL   string
    Flash FlashBag
}
```

### Error Level Pyramid (Build > Start > Runtime)

```
Build time:  dreego generate aborts, go build fails
             → Template syntax, missing variables, invalid sections
             → No output, non-zero exit

Start time:  server.Listen() aborts
             → Missing config, port occupied, plugin init error
             → log.Fatal / slog.Error + os.Exit(1)

Runtime:     Per-request, recovery catches everything
             → Panic → Recovery → 500 HTML
             → ValidationError → 422 HTML/JSON
             → RedirectError → 302 + Flash-Cookie
             → HTTPError → corresponding status code
             → Unexpected → 500 (Dev: Stack trace, Prod: generic)
```

### Recovery Middleware (Core-Fixed)

```go
// dreego-core/recovery.go
func Recovery(log *slog.Logger) func(http.Handler) http.Handler
```

Flow:
1. `defer recover()` — catches all panics
2. Logs via `slog.Error` with RequestID
3. Dev mode: Renders error page with stack trace
4. Prod mode: Renders generic `_error.dreego` or fallback `500.html`

### Middleware Stack (Core-Fixed + Core-Conditional, fixed order)

```
[Recovery → RequestID → RealIP → RequestLogging*]
  → [User middleware / Plugin middleware]
    → Router → Handler
```

\* `RequestLogging` is Core-Conditional: on by default, can be disabled via `dreego/config.json` (`logging.enabled: false`). A plugin (`dreego-logging`) can replace it in V2.

- `Recovery` — Panic → 500 + log
- `RequestID` — X-Request-ID header, in context and logs
- `RealIP` — X-Forwarded-For / X-Real-IP evaluation
- `RequestLogging` — slog.Info per request (method, path, status, duration)

### Dev vs Prod

```go
func IsDev() bool { return os.Getenv("APP_ENV") != "production" }
```

- Dev: Stack traces in browser, detailed error messages, `slog.LevelDebug`
- Prod: Generic 500 page, `slog.LevelInfo`, no internal details

### Form Validation Feedback

CodeGen generates handlers that return `ValidationError`. Rendering:

```html
<input name="email" value="{{ c.Old("email") }}" />
{#if c.Errors("email")}
  <p class="error">{{ c.Errors("email") }}</p>
{/if}

{#if c.Flash("success")}
  <div class="success">{{ c.Flash("success") }}</div>
{/if}
```

`c.Old()`, `c.Errors()`, `c.Flash()` are context methods (not magic variables).

### Error Boundary (per component)

No error boundary at component level in V1. Rationale:
- `<go>` block runs before template rendering — all errors are known beforehand
- `{#if hasError}` covers the use case
- Error boundaries are an SPA concept (React), not needed for SSR

→ See [no-catch-tag](no-catch-tag.md): errors via `{#if hasError}`, no special tag.

### Logging Strategy

```go
// dreego-core/logging.go
func RequestLogging(log *slog.Logger) func(http.Handler) http.Handler
```

- Uses Go's `log/slog` (structured, level-based)
- RequestID from context in every log entry
- Log level via config: `APP_LOG_LEVEL=debug|info|warn|error`
- Not a plugin — Core-fixed, always active, cannot be deactivated

### Integration with Plugin System

Plugins return errors as typed values:

```go
// In auth plugin
func (p *AuthPlugin) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, err := p.validateToken(r)
        if err != nil {
            dreego.WriteError(w, r, dreego.NewHTTPError(401, "Not authenticated", err))
            return
        }
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

`dreego.WriteError(w, r, err)` dispatches on type and renders accordingly — the recovery middleware is the central dispatcher.

## Consequences

- New files: `dreego-core/errors.go`, `dreego-core/recovery.go`, `dreego-core/logging.go` (exists as concept, will be populated)
- All handlers (including generated ones) return `error` — Recovery dispatches
- `dreego.Context` is extended with `Errors()`, `Old()`, `Flash()`
- `dreego generate` validates template syntax → build errors before runtime
- `slog` is a core dependency (Go 1.21+ stdlib, Go 1.22+ used)
- No Chi — all middleware built from scratch on `net/http`
