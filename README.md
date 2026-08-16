# dreego — Go Web Framework

SSR-First web framework for Go. Write `.dreego` files, transpile to Go code, deploy as single binary. File-based routing, built-in form handling, compile-time validation — no runtime magic.

```html
<!-- dreego/routes/login/post.dreego -->
<head><title>Dreego</title></head>

<go>
    type LoginForm struct {
        Email string `form:"email" validate:"required,email"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        c.SetSessionVal("user", form.Email)
        return c.Redirect("/dashboard", 303)
    }
</go>

<div>
    <h1>Login</h1>
    {#if c.Errors("email")}<p class="error">{{ c.Errors("email") }}</p>{/if}
    <form g-action="Login" method="post">
        <input name="email" type="email" value="{{ c.Old("email") }}">
        <button type="submit">Login</button>
    </form>
</div>
```

```bash
dreego generate && go run .
```

→ **[Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md)**

## Philosophy

Dreego is a **compile-time transpiler**, not a runtime framework. `.dreego` files compile to standard Go code — no reflection-based routers, no runtime template parsing. Your app is a plain Go binary using `net/http`. See [Benchmarks](_docs/benchmarks.md) for measured code-generation and request performance.

Four principles:
1. **SSR-First** — Pages render server-side. HTMX/Alpine.js for progressive enhancement, not required.
2. **File-Based** — The current pre-v0.1 router maps `dreego/routes/login/get.dreego` to `GET /login`. The accepted v0.1 migration will use one route file per URL with method-specific sections.
3. **Type-Safe** — Generated handlers and components use typed Go contracts; dynamic HTTP boundary data stays explicit.
4. **Accessible by Default** — CLI output, diagnostics, blueprints, and official components are designed for screen readers, keyboards, and semantic HTML. Applications still verify their own content and conformance.

## Features

### Core
- **Transpiler Pipeline** — Lexer → Parser → AST → CodeGen. `.dreego` → Go code.
- **File-based Routing** — `dreego/routes/get.dreego` → `GET /`, `dreego/routes/login/post.dreego` → `POST /login`
- **Dynamic Segments** — `[id]` brackets for URL params, `(group)/` for layout groups
- **Single Binary** — `go build` → deploy one file. Zero runtime dependencies beyond `net/http`.
### Template & Components
- **Template Logic** — `{{ value }}`, `{#if}...{#else}...{/if}`, `{#each items as item}...{#each else}...{/each}`
- **Template Helpers** — `{{ $loop.Index }}`, `{{ value|raw }}`, `{{ value|upper }}`, `{#verbatim}`
- **Component System** — `dreego/components/`, `<@Card title="x">...<\@Card>`, named slots, scoped CSS
- **Layout System** — `dreego/layouts/default.dreego` with `{#slot}` + `{#head}`
- **CSS Scoping** — `data-scope` via source hash, automatically applied

### Form Handling (v0.0.16)
- **g-action Forms** — `<form g-action="Login">` generates full handler pipeline
- **Auto-Validation** — `validate:"required,email,min=3"` struct tags, no external deps
- **Auto-Binding** — `form:"email"` struct tags → `r.ParseForm()` + reflection mapping
- **Error Re-render** — `c.Errors("email")` and `c.Old("email")` available in templates after validation failure
- **PRG Pattern** — `c.Redirect(url, 303)` with `ErrRedirect` sentinel for after-success redirect

### API Endpoints (v0.0.15)
- **Content-Type Routing** — `<go type="json">`, `<go type="xml">`, `<go type="custom">`
- **JSON API** — `c.JSON(200, data)`, `c.Bind(&target)`, auto `Content-Type: application/json`
- **XML API** — `c.XML(200, data)`, auto `Content-Type: application/xml`
- **Custom** — `c.Write(status, contentType, body)` for arbitrary formats

### Middleware (v0.0.14, v0.0.20)
- **Health Checks** — `GET /health` → 200, `GET /ready` → 200/503 via `app.SetReady(bool)`
- **Security Headers** — X-Content-Type-Options, X-Frame-Options, Referrer-Policy, Permissions-Policy, Content-Security-Policy (configure via `app.SetCSP` before build)
- **Gzip Compression** — `Accept-Encoding` → compressed response wrapping
- **Recovery** — Panic → 500 with stack trace
- **Request Logging** — JSONL format with duration, IP, status
- **Session** — Cookie store via `app.SetSessionStore()`, `c.SetSessionVal()`
- **CSRF** — Double-submit cookie, auto-validation on POST/PUT/DELETE, Secure flag TLS-aware

### Developer Experience
- **CLI** — `dreego init`, `dreego generate [--force] [--check]`, `dreego fmt [--check]`
- **CI Mode** — `dreego generate --check` exits non-zero when generated files are stale
- **Auto-Imports** — `fmt`, `html`, `strings`, `net/http` added to generated code as needed
- **147 Integration Tests** — Docker-based, all pass

## Quick Start

```bash
go install github.com/dreego-stack/dreego/cli/dreego@latest
dreego init myapp
cd myapp
go mod init myapp
go mod edit -replace github.com/dreego-stack/dreego/core=../dreego/core  # or use go get
dreego generate
go run .

# or build for production
dreego build --target linux/amd64
docker build -t myapp .
```

## Architecture

```
dreego/
├── routes/           # .dreego files → URL routes
│   ├── get.dreego        → GET /
│   ├── login/
│   │   ├── get.dreego    → GET /login
│   │   └── post.dreego   → POST /login
│   └── [id]/get.dreego   → GET /{id}
├── layouts/
│   └── default.dreego    # {#slot} + {#head} wrapper
├── components/
│   └── Card.dreego       # <@Card title="x"/>
├── static/
│   └── style.css         # inlined into binary
├── gen/                  # GENERATED — do not edit
│   ├── routes.go
│   ├── components.go
│   └── dree.go
└── main.go
```

## Plugins

Official plugins live in separate repos under `github.com/dreego-stack/`. Each plugin has its own `go.mod` and requires `github.com/dreego-stack/dreego`; Core stays dependency-free and never imports a plugin package.

```
github.com/dreego-stack/
├── dreego/             # main repo (core + CLI, single module)
├── plugin-example/     # minimal example plugin
├── plugin-auth/        # future: OAuth2, JWT, sessions
├── plugin-db/          # future: SQL drivers, migrations
└── ...
```

→ **[Plugin System](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md)**

## Documentation

| Doc | Topic |
|-----|-------|
| `_docs/routing.md` | File-based routing, dynamic segments, content-type routing |
| `_docs/middleware.md` | Health, security, compression, session, CSRF |
| `_docs/forms.md` | g-action forms, validation, redirects, error handling |
| `_docs/components.md` | Component system, slots, scoped CSS |
| `_docs/hot-reload.md` | Hot reload with Air (.air.toml config) |
| `_docs/runtime.md` | Full Go API reference |
| `_docs/getting-started.md` | Step-by-step tutorial |

## License

MPL-2.0
