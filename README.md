# dreego — Go Web Framework

SSR-First web framework for Go. Write `.dreego` files, transpile to Go code, deploy as single binary. File-based routing, built-in form handling, compile-time validation — no runtime magic.

```html
<!-- www/routes/login/post.dreego -->
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
dreego new myapp && cd myapp && dreego generate && go run .
```

→ **[Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md)**

## Philosophy

Dreego is a **compile-time transpiler**, not a runtime framework. `.dreego` files compile to standard Go code — no reflection-based routers, no runtime template parsing. Your app is a plain Go binary using `net/http`. See [Benchmarks](_docs/benchmarks.md) for measured code-generation and request performance.

Four principles:
1. **SSR-First** — Pages render server-side. HTMX/Alpine.js for progressive enhancement, not required.
2. **File-Based** — The current pre-v0.1 router maps `www/routes/login/get.dreego` to `GET /login`. The accepted v0.1 migration will use one route file per URL with method-specific sections.
3. **Type-Safe** — Generated handlers and components use typed Go contracts; dynamic HTTP boundary data stays explicit.
4. **Accessibility-Aware Tooling** — CLI output and diagnostics are designed for screen readers, and the landing blueprint demonstrates semantic navigation. Applications still verify their own content and conformance.

## Features

### Core
- **Transpiler Pipeline** — Lexer → Parser → AST → CodeGen. `.dreego` → Go code.
- **File-based Routing** — `www/routes/get.dreego` → `GET /`, `www/routes/login/post.dreego` → `POST /login`
- **Dynamic Segments** — `[id]` brackets for URL params, `(group)/` for layout groups
- **Single Binary** — `go build` → deploy one file. Zero runtime dependencies beyond `net/http`.

### Template & Components
- **Template Logic** — `{{ value }}`, `{#if}...{#else}...{/if}`, `{#each items as item}...{#each else}...{/each}`
- **Template Helpers** — `{{ $loop.Index }}`, `{{ value|raw }}`, `{{ value|upper }}`, `{#verbatim}`
- **Component System** — `www/components/`, `<@Card title="x">...<\@Card>`, named slots, scoped CSS
- **Layout System** — `www/layouts/default.dreego` with `{#slot}` + `{#head}`
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
- **Accessibility Checks** — `dreego generate` warns about missing image alternatives and unassociated form labels; CLI output is color-free and screen-reader-linear

## Quick Start

Requires Go 1.22 or newer (`go version`). The `dreego new` command scaffolds a
project, writes a `go.mod` that requires the published `dreego` module, and
runs `go mod tidy` so the build resolves from the public Go proxy. No
repository-local `replace` directive is needed for a release-installed CLI.

```bash
# 1. Install the Dreego CLI
go install github.com/dreego-stack/dreego/cli/dreego@latest

# 2. Scaffold a new project (writes go.mod, main.go, www/ tree, runs go mod tidy)
dreego new myapp

# 3. Generate Go code from .dreego files and run the server
cd myapp
dreego generate
go run .
# → http://localhost:8080

# or build for production
dreego build --target linux/amd64
docker build -t myapp .
```

`main.go` uses the explicit App API — no globals, no hidden state:

```go
package main

import (
	"log"

	dreego "github.com/dreego-stack/dreego/core"
	"myapp/www"
)

func main() {
	app := dreego.New()
	if err := www.Register(app); err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

If `dreego new` cannot resolve the `dreego` dependency (e.g. no network, or the
CLI was built from an untagged checkout), set `DREEGO_LOCAL_REPO=/path/to/dreego`
to point the scaffold at a local checkout via a `replace` directive. This is a
developer escape hatch and is not part of the canonical path.

## VS Code Extension

Syntax highlighting, snippets, and `dreego.config.json` validation for
`.dreego` files. Install via a symlinked clone — always up to date:

```sh
git clone --depth 1 https://github.com/dreego-stack/vscode-dreego ~/dev/vscode-dreego
cd ~/dev/vscode-dreego
./install.sh
```

Restart VS Code afterwards. Uninstall: `./install.sh uninstall` in the same
directory. The install script links `~/.vscode/extensions/dreego` to the
checkout, so use a permanent directory (not `/tmp`, which is cleared on
reboot). Source: https://github.com/dreego-stack/vscode-dreego

## Architecture

A website lives in its own directory — any name, marked by a
`dreego.config.json`. `dreego generate` produces one `dree.go` per directory
with `.dreego` sources; the website root gets a `Register(app)` entry point:

```
www/                       # website root (name is free, marker: dreego.config.json)
├── dreego.config.json     # logging, redirects, rewrites
├── dree.go                # GENERATED — package www, Register(app)
├── routes/                # .dreego files → URL routes
│   ├── get.dreego             → GET /
│   ├── login/
│   │   ├── get.dreego         → GET /login
│   │   └── post.dreego        → POST /login
│   ├── [id]/get.dreego        → GET /{id}
│   └── dree.go             # GENERATED — package routes, handlers + Register
├── layouts/
│   ├── default.dreego      # {#slot} + {#head} wrapper
│   └── dree.go             # GENERATED — package layouts
├── components/
│   ├── Card.dreego         # <@Card title="x"/>
│   └── dree.go             # GENERATED — package components
├── static/
│   └── style.css           # inlined into binary
└── main.go                 # imports "myapp/www", calls www.Register(app)
```

Multiple websites can share one module — each directory with a
`dreego.config.json` is an independent website.

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

### Getting Started

| Doc | Topic |
|-----|-------|
| `_docs/index.md` | Documentation index and navigation |
| `_docs/getting-started.md` | Step-by-step tutorial |
| `_docs/cli.md` | CLI Reference |
| `_docs/config.md` | dreego.config.json |

### Guides

| Doc | Topic |
|-----|-------|
| `_docs/routing.md` | File-based routing, dynamic segments, content-type routing |
| `_docs/layouts.md` | Layouts, `{#slot}` and `{#head}` |
| `_docs/middleware.md` | Health, security, compression, session, CSRF |
| `_docs/components.md` | Component system, slots, scoped CSS |
| `_docs/forms.md` | g-action forms, validation, redirects, error handling |
| `_docs/session-encryption.md` | AES-256-GCM encrypted session cookies |
| `_docs/progressive-enhancement.md` | HTMX, Alpine.js, plain JavaScript |
| `_docs/security.md` | Context-aware escaping, output safety |
| `_docs/accessibility.md` | Accessibility guarantees and blueprint defaults |

### Reference

| Doc | Topic |
|-----|-------|
| `_docs/runtime.md` | Practical guide to the main Go runtime APIs |
| `_docs/plugin-interfaces.md` | Plugin interface contracts |
| `_docs/plugins.md` | Plugin model, middleware and route hooks |
| `_docs/compatibility.md` | Breaking-change policy and the v0.1 promise |
| `_docs/roadmap.md` | Product direction and release phases |

### Testing & Ops

| Doc | Topic |
|-----|-------|
| `_docs/testing.md` | Integration test strategy |
| `_docs/benchmarks.md` | Code generation and request benchmarks |
| `_docs/reference-apps.md` | Reference applications under `_tests/fixtures/` |
| `_docs/dev-server.md` | `dreego dev` watcher and auto-reload |
| `_docs/hot-reload.md` | Hot reload with Air (.air.toml config) |
| `_docs/deployment.md` | Build, cross-compile, containers |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the PR workflow, commit
conventions, and local development setup.

## License

MPL-2.0
    
