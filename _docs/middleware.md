# Middleware

## Core Middleware

Dreego has a middleware chain with fixed order:

```
[Recovery → RequestID → RealIP → RequestLogging*]
  → Redirect/Rewrite
    → Router → Handler
```

\* `RequestLogging` is Core-Conditional: default on, deactivatable via `dreego/config.json`.

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

## Planned (V1)

- `Recovery`: Panic → 500 (Core-Fixed)
- `RequestID`: X-Request-ID Header (Core-Fixed)
- `RealIP`: X-Forwarded-For evaluation (Core-Fixed)
- `CSRF`: CSRF protection (Core-Conditional)
