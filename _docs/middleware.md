# Middleware

## Core Middleware

Dreego has a middleware chain with fixed order:

```
[Recovery → SecurityHeaders → Compression → RequestLogging*]
  → Session → CSRF
    → Redirect/Rewrite
      → Router (mux)
```

\* `RequestLogging` is Core-Conditional: default on, deactivatable via `dreego.config.json`.

## Health Checks

Built-in `GET /health` and `GET /ready` endpoints, always available:

- `GET /health` → 200 `ok` — process is alive
- `GET /ready` → 200 `ready` or 503 `not ready` — traffic readiness

```go
app.SetReady(false) // signal not ready (e.g., during startup)
app.SetReady(true)  // signal ready
```

Health endpoints are registered before user routes — they cannot be overridden.

## Security Headers (CSP since v0.0.20)

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
if err := app.SetCSP("default-src 'self'"); err != nil {
    log.Fatal(err)
}
```

Call `app.SetCSP` before `app.Build`, `app.Handler`, `app.ServeHTTP`, or `ssr.Listen`. An empty string falls back to `default-src 'self'`.

Always on. Applied after Recovery, before Compression.

## Compression

Gzip compression for all responses, core-fixed:

- Negotiates via `Accept-Encoding` q-values: `gzip`, `gzip;q=0.5`, `*` (case-insensitive)
- `gzip;q=0` disables compression, also when a wildcard with a non-zero q-value is present
- Compresses response body via `compress/gzip`
- Sets `Content-Encoding: gzip`
- Always appends `Vary: Accept-Encoding` (existing `Vary` values are preserved)
- Removes a stale `Content-Length` when compressing; preserves it when not compressing
- Skips HEAD responses and responses with status 204, 304, or an existing `Content-Encoding`

Responses are buffered in memory before being sent, so a panic in a downstream
handler can be turned into one plain error response (see Panic Recovery below).
Calling `Flush()` commits the buffered gzip member and continues with a new
member — multi-member gzip streams are valid per RFC 1952 and decompress
transparently. Handlers using `io.Copy` stay compressed. Informational 1xx
responses (e.g. 103 Early Hints) are forwarded to the client immediately and
never become the final status.

The wrapped writer preserves `http.Flusher`, `http.Hijacker`, `http.Pusher`,
`io.ReaderFrom`, and exposes `Unwrap()`. `Hijack`/`Push` return
`http.ErrNotSupported` when the upstream writer does not support them.

Applied after Security Headers, before RequestLogging.

## Panic Recovery

`Recovery` runs outside Compression. When a handler panics, any buffered
compressed bytes are discarded and one plain 500 response (plus `Vary:
Accept-Encoding`) is written with the intended status and no `Content-Encoding`.
Responses already committed to the wire via `Flush` cannot be rewound; the
panic handler output follows the partial stream.

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

Configured in `dreego.config.json` → `redirects` and `rewrites`.

## Plugin Middleware

Plugins register middleware through `app.Use()` in source order. There is no
central `MiddlewareProvider` interface before v1; a plugin exposes an explicit
`Register(app, options) error` function that calls `app.Use` and `app.Register`
directly. See [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md).
