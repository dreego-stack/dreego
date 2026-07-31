# Middleware

## Core Middleware

Dreego has a middleware chain with fixed order:

```
[Recovery → SecurityHeaders → Compression → RequestLogging*]
  → Session → CSRF
    → Redirect/Rewrite
      → Router (mux)
```

\* `RequestLogging` is Core-Conditional: default on, deactivatable via `dreego/config.json`.

## Health Checks (v0.0.14)

Built-in `GET /health` and `GET /ready` endpoints, always available:

- `GET /health` → 200 `ok` — process is alive
- `GET /ready` → 200 `ready` or 503 `not ready` — traffic readiness

```go
core.SetReady(false) // signal not ready (e.g., during startup)
core.SetReady(true)  // signal ready
```

Health endpoints are registered before user routes — they cannot be overridden.

## Security Headers (v0.0.14, CSP v0.0.20)

Core-fixed middleware that sets security headers on every response:

| Header | Value |
|--------|-------|
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `DENY` |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | `geolocation=(), microphone=(), camera=()` |
| `Content-Security-Policy` | permissive default (see below) |

The default CSP allows `self`, `unsafe-inline` for scripts/styles (so HTMX/Alpine.js and scoped CSS work out of the box), and common CDN/font sources. Override it for stricter setups:

```go
core.SetCSP("default-src 'self'")
```

Call `core.SetCSP` before `core.Listen`. An empty string falls back to `default-src 'self'`.

Always on. Applied after Recovery, before Compression.

## Compression (v0.0.14)

Gzip compression for all responses, core-fixed:

- Checks `Accept-Encoding: gzip` header
- Compresses response body via `compress/gzip`
- Sets `Content-Encoding: gzip`

Applied after Security Headers, before RequestLogging.

## RequestLogging

Logs each request as a JSONL line:

```jsonl
{"time":"2026-07-24T21:42:12","method":"GET","path":"/","status":200,"ip":"[::1]:56365","duration":"13.584µs"}
```

Fields: `time`, `method`, `path`, `status`, `ip`, `duration`.

Configuration:

```json
{ "logging": { "enabled": true } }
```

## Redirect/Rewrite

Executed before the router. Redirects redirect (301/302), rewrites change the path transparently.

Configured in `dreego/config.json` → `redirects` and `rewrites`.

## Plugin Middleware

Plugins implement `MiddlewareProvider` and inject their own middleware into the chain. Order = `app.Use()` order.
