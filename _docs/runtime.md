# Runtime API

This guide covers the primary runtime APIs available after `import dreego "github.com/dreego-stack/dreego/core"`. The exported Go package remains the complete API reference.

## SSRContext

Available as **`c`** in routes (including `<go>` blocks and error pages) and as **`ctx`** in components (see [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md)). The name is fixed by the generated code: routes are generated as `func renderX(c *dreego.SSRContext)`, components as `func(ctx *dreego.SSRContext)`.

| Method | Returns | Description |
|--------|---------|-------------|
| `c.Param("id")` | `string` | URL path parameter from `[id]` segment |
| `c.Query("ref")` | `string` | URL query parameter `?ref=x` |
| `c.FormValue("name")` | `string` | POST form value (calls `ParseForm`; returns `""` on parse failure) |
| `c.FormError()` | `error` | Non-nil when the last `FormValue` call failed to parse the form body; `nil` for a valid empty form |
| `c.Data("key")` | `any` | Arbitrary data stored in context |
| `c.Set("key", val)` | — | Store data for use between nested calls |
| `c.Get("key")` | `string` | Retrieve string data (used for slot passing) |
| `c.SessionVal("key")` | `string` | Session value (requires `SetSessionStore`; `""` on store failure) |
| `c.SetSessionVal("k","v")` | — | Write session value (secure defaults) |
| `c.DelSessionVal("key")` | — | Delete single session key |
| `c.DestroySession()` | — | Destroy entire session |
| `c.CSRFToken()` | `string` | Current CSRF token (from session) |
| `c.SessionError()` | `error` | Non-nil when the last session read/write/delete/destroy call failed; `nil` otherwise |
| `c.R` | `*http.Request` | Raw request (use sparingly) |
| `c.W` | `http.ResponseWriter` | Raw writer (use sparingly) |

Session and form failures are not silently dropped. `FormError()` and
`SessionError()` expose the underlying cause to the application so a route can
decide how to respond. Generated handlers and middleware respond with a
generic `500 Internal Server Error` and never disclose filesystem paths,
database details, or Go type errors to clients; the internal cause remains
available to the application and to structured logs.

## Server

| Method | Description |
|----------|-------------|
| `app.Listen(":8080")` | Start HTTP server with the App middleware chain |
| `app.Shutdown(ctx)` | Gracefully shut down a running server (drains active requests) |
| `app.Handler()` | Freeze configuration and return the App's `http.Handler` |
| `app.ServeHTTP(w, r)` | Use the App directly as an `http.Handler` |

`Listen` blocks until the server stops. It installs SIGINT/SIGTERM handlers,
drains active requests within the shutdown deadline, and returns a non-nil
error if draining did not finish in time so the caller can decide how to
report the unfinished work. Signal subscriptions are released when `Listen`
returns, so repeated server lifecycles do not leak goroutines.

### Server timeouts and limits

`app.SetServerConfig(ServerConfig)` tunes the HTTP server before build:

| Field | Default | Effect |
|-------|---------|--------|
| `ReadHeaderTimeout` | 10s | Aborts clients that send headers too slowly (slowloris protection) |
| `ReadTimeout` | 30s | Caps the total time to read the request, including the body |
| `WriteTimeout` | 30s | Caps the time from end-of-headers to fully writing the response |
| `IdleTimeout` | 120s | Closes keep-alive connections that are idle longer than this |
| `MaxHeaderBytes` | 1 MiB | Rejects oversized request headers with a 431 response |
| `ShutdownTimeout` | 10s | Deadline for draining active requests during shutdown |

Connection and header timeouts are app-wide server policy for v0.1; they
cannot be relaxed for a single route. Zero values fall back to the secure
defaults in the table above, so timeouts and limits cannot be disabled for
v0.1 — enforced defaults are the point. Without them, a missing
`ReadHeaderTimeout` would remove slowloris protection, a missing
`ReadTimeout` would let clients hold a connection open while streaming the
body indefinitely, and a missing `WriteTimeout` would let slow clients block
response goroutines.

### Request body limits

Request-body limits are **the application's responsibility** — Dreego does not
apply an implicit body cap by default (only header limits are enforced by the
server). An unguarded route that accepts a request body is therefore exposed to
unbounded uploads; wrap your POST/PUT routes with `MaxBodyReader` to set a sane
limit.

Use the `MaxBodyReader(max int64)` middleware to cap the request body for a
specific route without weakening unrelated routes:

```go
upload := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if _, err := io.Copy(io.Discard, r.Body); err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	w.WriteHeader(http.StatusOK)
})

app.Use(core.MaxBodyReader(1 << 20)) // 1 MiB application-wide default
app.Register("POST", "/upload", upload) // inherits the default
large := core.MaxBodyReader(64 << 20)(upload)
app.Register("POST", "/large-upload", large.ServeHTTP)
```

Application middleware runs before CSRF form parsing, so an app-wide
`MaxBodyReader` limit also protects CSRF-protected form posts. CSRF maps an
oversized form to `413 Payload Too Large`. Other handlers receive the stdlib
`*http.MaxBytesError` and must map it to a response themselves.
Streaming and upload exceptions are intentionally not built into core; raise
the limit per-route through handler composition when a real plugin proves the
requirement.

## Session

| Function | Description |
|----------|-------------|
| `dreego.NewCookieStore([]byte("01234567890123456789012345678901"))` | Create HMAC-signed cookie session store (secret must be at least 32 bytes) |
| `app.SetSessionStore(store)` | Enable sessions for this App before build |

Session cookies use secure defaults: `HttpOnly: true`, `SameSite: Lax`, `Secure: TLS-aware`, `Path: "/"`.

Configure a non-root path once with `CookiePolicy.Path`. Per-call path
overrides are rejected with `ErrCookiePathOverride`, ensuring `Delete` and
`Destroy` always expire the same browser cookie that `Set` created.

For AES-256-GCM session encryption see [Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md).

## Configuration

| Method | Description |
|----------|-------------|
| `app.SetLogging(bool)` | Enable/disable request logging (JSONL format) before build |
| `app.SetCSRF(bool)` | Enable/disable CSRF protection for this App before build (default: on) |
| `app.SetCSP(value string)` | Override Content-Security-Policy before build |
| `app.SetErrorHandler(code, handler)` | Configure an App-local error handler before build |
| `app.SetReady(bool)` | Change readiness dynamically, including after build |
| `app.SetServerConfig(ServerConfig)` | Tune HTTP server timeouts and limits before build |

Configuration freezes on the first call to `Build`, `Handler`, `ServeHTTP`, or
`Listen`. Later configuration calls return `dreego.ErrAppBuilt`; readiness is
the only intentionally dynamic setting.

## Configuration File

`dreego/config.json` controls redirects, rewrites, and logging:

```json
{
    "logging": {"enabled": true},
    "redirects": [
        {"from": "/old", "to": "/new", "status": 301}
    ],
    "rewrites": [
        {"from": "/api/*", "to": "/v2/*"}
    ]
}
```

## Static Assets

```go
if err := app.RegisterStatic("/style.css", "text/css", []byte("body{color:red}")); err != nil {
    log.Fatal(err)
}
```

Generated automatically from `dreego/static/` by `dreego generate`. MIME type detected from file extension.

## main.go Pattern

```go
package main

import (
	"log"

	"myapp/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	app := dreego.New()
	store := dreego.NewCookieStore([]byte("01234567890123456789012345678902"))
	if err := app.SetSessionStore(store); err != nil {
		log.Fatal(err)
	}
	if err := gen.Register(app); err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

Install core: `go get github.com/dreego-stack/dreego/core@latest`

## See Also

- [Docs Index](https://github.com/dreego-stack/dreego/blob/main/_docs/index.md)
- [Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md)
- [Config Reference](https://github.com/dreego-stack/dreego/blob/main/_docs/config.md)
