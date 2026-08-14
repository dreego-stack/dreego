# Runtime API

Full API surface available after `import dreego "github.com/dreego-stack/dreego/core"`.

## SSRContext

Available as **`c`** in routes (including `<go>` blocks and error pages) and as **`ctx`** in components (see [Components](https://github.com/dreego-stack/dreego/blob/main/_docs/components.md)). The name is fixed by the generated code: routes are generated as `func renderX(c *dreego.SSRContext)`, components as `func(ctx *dreego.SSRContext)`.

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

| Method | Description |
|----------|-------------|
| `app.Listen(":8080")` | Start HTTP server with the App middleware chain |
| `app.Handler()` | Freeze configuration and return the App's `http.Handler` |
| `app.ServeHTTP(w, r)` | Use the App directly as an `http.Handler` |

## Session

| Function | Description |
|----------|-------------|
| `dreego.NewCookieStore([]byte("secret-32-bytes"))` | Create HMAC-signed cookie session store |
| `app.SetSessionStore(store)` | Enable sessions for this App before build |

Session cookies use secure defaults: `HttpOnly: true`, `Secure: TLS-aware`, `Path: "/"`.

For AES-256-GCM session encryption see [Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md).

## Configuration

| Method | Description |
|----------|-------------|
| `app.SetLogging(bool)` | Enable/disable request logging (JSONL format) before build |
| `app.SetCSRF(bool)` | Enable/disable CSRF protection for this App before build (default: on) |
| `app.SetCSP(value string)` | Override Content-Security-Policy before build |
| `app.SetErrorHandler(code, handler)` | Configure an App-local error handler before build |
| `app.SetReady(bool)` | Change readiness dynamically, including after build |

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
	store := dreego.NewCookieStore([]byte("super-secret-key-32-bytes!"))
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
