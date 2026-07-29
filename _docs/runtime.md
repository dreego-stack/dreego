# Runtime API

Full API surface available after `import core "codeberg.org/dreego/dreego/core"`.

## SSRContext

Available as `c` in `<go>` blocks and component render functions.

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
| `core.Listen(":8080")` | Start HTTP server with full middleware chain |
| `core.ServeMux()` | Build `http.Handler` with all routes, middleware, session, CSRF |

## Session

| Function | Description |
|----------|-------------|
| `core.NewCookieStore([]byte("secret-32-bytes"))` | Create HMAC-signed cookie session store |
| `core.SetSessionStore(store)` | Enable sessions for all requests |

Session cookies use secure defaults: `HttpOnly: true`, `Secure: TLS-aware`, `Path: "/"`.

## Configuration

| Function | Description |
|----------|-------------|
| `core.SetLogging(bool)` | Enable/disable request logging (JSONL format) |
| `core.SetCSRF(bool)` | Enable/disable CSRF protection (default: on) |
| `core.SetErrorHandler(code, handler)` | Custom handler for HTTP errors (500 used by Recovery) |

## Configuration File

`dreego/config.json` controls redirects, rewrites, and logging:

```json
{
    "logging": true,
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
core.RegisterStatic("/style.css", "text/css", []byte("body{color:red}"))
```

Generated automatically from `dreego/static/` by `dreego generate`. MIME type detected from file extension.

## main.go Pattern

```go
package main

import (
    _ "myapp/dreego/gen"
    core "codeberg.org/dreego/dreego/core"
)

func main() {
    store := core.NewCookieStore([]byte("super-secret-key-32-bytes!"))
    core.SetSessionStore(store)
    core.Listen(":8080")
}
```

## See Also

- [Docs Index](https://codeberg.org/dreego/dreego/src/branch/main/_docs/index.md)
- [Getting Started](https://codeberg.org/dreego/dreego/src/branch/main/_docs/getting-started.md)
- [Config Reference](https://codeberg.org/dreego/dreego/src/branch/main/_docs/config.md)
