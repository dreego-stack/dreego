# Runtime API

Full API surface available after `import dreego "codeberg.org/dreego/dreego/core"`.

## SSRContext

Available as **`c`** in routes (including `<go>` blocks and error pages) and as **`ctx`** in components (see [Components](https://codeberg.org/dreego/dreego/src/branch/main/_docs/components.md)). The name is fixed by the generated code: routes are generated as `func renderX(c *dreego.SSRContext)`, components as `func(ctx *dreego.SSRContext)`.

| Method | Returns | Description |
|--------|---------|-------------|
| `c.Param("id")` | `string` | URL path parameter from `[id]` segment |
| `c.Query("ref")` | `string` | URL query parameter `?ref=x` |
| `c.FormValue("name")` | `string` | POST form value (calls `ParseForm`) |
| `c.Data("key")` | `any` | Arbitrary data stored in context |
| `c.Set("key", val)` | — | Store data for use between nested calls |
| `c.Get("key")` | `string` | Retrieve string data (used for slot passing) |
| `c.SessionVal("key")` | `string` | Session value (requires `SetSessionStore`) |
| `c.SetSessionVal("k","v")` | — | Write session value (secure defaults) |
| `c.DelSessionVal("key")` | — | Delete single session key |
| `c.DestroySession()` | — | Destroy entire session |
| `c.CSRFToken()` | `string` | Current CSRF token (from session) |
| `c.R` | `*http.Request` | Raw request (use sparingly) |
| `c.W` | `http.ResponseWriter` | Raw writer (use sparingly) |

## Server

| Function | Description |
|----------|-------------|
| `dreego.Listen(":8080")` | Start HTTP server with full middleware chain |
| `dreego.ServeMux()` | Build `http.Handler` with all routes, middleware, session, CSRF |

## Session

| Function | Description |
|----------|-------------|
| `dreego.NewCookieStore([]byte("secret-32-bytes"))` | Create HMAC-signed cookie session store |
| `dreego.SetSessionStore(store)` | Enable sessions for all requests |

Session cookies use secure defaults: `HttpOnly: true`, `Secure: TLS-aware`, `Path: "/"`.

For AES-256-GCM session encryption see [Session Encryption](https://codeberg.org/dreego/dreego/src/branch/main/_docs/session-encryption.md).

## Configuration

| Function | Description |
|----------|-------------|
| `dreego.SetLogging(bool)` | Enable/disable request logging (JSONL format) |
| `dreego.SetCSRF(bool)` | Enable/disable CSRF protection (default: on) |
| `dreego.SetCSP(value string)` | Override the Content-Security-Policy header (empty falls back to `default-src 'self'`) |
| `dreego.SetErrorHandler(code, handler)` | Custom handler for HTTP errors (500 used by Recovery) |

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
dreego.RegisterStatic("/style.css", "text/css", []byte("body{color:red}"))
```

Generated automatically from `dreego/static/` by `dreego generate`. MIME type detected from file extension.

## main.go Pattern

```go
package main

import (
    _ "myapp/dreego/gen"
    dreego "codeberg.org/dreego/dreego/core"
)

func main() {
    store := dreego.NewCookieStore([]byte("super-secret-key-32-bytes!"))
    dreego.SetSessionStore(store)
    dreego.Listen(":8080")
}
```

Install core: `go get codeberg.org/dreego/dreego/core@latest`

## See Also

- [Docs Index](https://codeberg.org/dreego/dreego/src/branch/main/_docs/index.md)
- [Getting Started](https://codeberg.org/dreego/dreego/src/branch/main/_docs/getting-started.md)
- [Config Reference](https://codeberg.org/dreego/dreego/src/branch/main/_docs/config.md)
