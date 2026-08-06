# Plugins

Dreego plugins live in the same repository under `plugins/`. Each plugin that needs external dependencies has its own `go.mod`, so Core stays dependency-free and users only pull in what they import. Core never imports any plugin package — plugins depend on Core, never the other way around.

## Architecture

```
codeberg.org/dreego/dreego              ← Root module: core/ + cmd/dreego/
codeberg.org/dreego/dreego/core         ← Core package (no external deps beyond stdlib)
codeberg.org/dreego/dreego/plugins/auth ← Plugin: OAuth2, JWT, sessions (own go.mod)
codeberg.org/dreego/dreego/plugins/db   ← Plugin: SQL drivers, migrations (own go.mod)
```

Plugins without external dependencies can also be plain packages inside the root module, but once a plugin needs a third-party dependency it gets its own `go.mod`.

The repository root contains a `go.work` file that links the root module and every plugin module for local development. Consumers of the framework only see the modules they explicitly import.

## Plugin Interface (v1, frozen)

Core defines a single `Plugin` interface. Plugins import Core, satisfy the interface, and register themselves with the runtime via `dreego.UsePlugin(p)`. Core never imports a plugin.

```go
// Defined in core
type Plugin interface {
    Name() string
    RegisterRoutes() // plugin calls dreego.Register(...) internally
    Middlewares() []func(http.Handler) http.Handler
    Assets() fs.FS
    OnStart(ctx context.Context) error
    OnShutdown(ctx context.Context) error
}
```

### Registration

`dreego.UsePlugin(p)` is the central v1 API. It is called at package level (typically from `main.go`), not on an app object. It registers the plugin's routes, middleware, assets and lifecycle hooks with the core runtime:

```go
func UsePlugin(p Plugin)
```

A plugin package imports Core and implements the interface:

```go
// codeberg.org/dreego/dreego/plugins/auth
package auth

import (
    "context"
    "io/fs"
    "net/http"

    dreego "codeberg.org/dreego/dreego/core"
)

type Auth struct{ secret string }

func New(secret string) *Auth { return &Auth{secret: secret} }

func (a *Auth) Name() string { return "auth" }

func (a *Auth) RegisterRoutes() {
    dreego.Register("GET", "/login", a.handleLogin)
    dreego.Register("POST", "/logout", a.handleLogout)
}

func (a *Auth) Middlewares() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{a.sessionMiddleware}
}

func (a *Auth) Assets() fs.FS { return nil }

func (a *Auth) OnStart(ctx context.Context) error    { return nil }
func (a *Auth) OnShutdown(ctx context.Context) error { return nil }
```

The application registers the plugin in `main.go`:

```go
package main

import (
    dreego "codeberg.org/dreego/dreego/core"
    "codeberg.org/dreego/dreego/plugins/auth"
)

func main() {
    dreego.UsePlugin(auth.New("secret"))
    // ...
}
```

### Lifecycle

Plugins are started and shut down in registration order via `dreego.StartPlugins(ctx)` and `dreego.ShutdownPlugins(ctx)`:

```go
func StartPlugins(ctx context.Context) error   // OnStart on every plugin
func ShutdownPlugins(ctx context.Context) error // OnShutdown on every plugin
```

**Abort behavior:** `StartPlugins` returns on the **first** `OnStart` error and stops — it does **not** start the remaining plugins. Plugins whose `OnStart` already succeeded before the failure are **not** shut down automatically, which can leak resources (open connections, goroutines, background workers). Callers must handle cleanup explicitly on an error from `StartPlugins`. The same early-abort applies to `ShutdownPlugins`.

### Middleware hooks (FIFO)

A plugin's `Middlewares()` are appended to the runtime middleware chain in **FIFO order** — the first registered plugin runs first on request entry, then the next, then the handler:

```go
dreego.UsePlugin(pluginA)  // A runs first
dreego.UsePlugin(pluginB)  // then B
```

The chain is **fixated on the first `Build()`**: registering a plugin after the handler is already built does **not** reorder the stack. To change middleware order you must register plugins before the first build/serve.

### Route hooks (programmatic routes)

A plugin registers routes by calling `dreego.Register(...)` **inside its `RegisterRoutes()`**. All such routes are served by `dreego.ServeMux()` alongside the file-based routes:

```go
func (a *Auth) RegisterRoutes() {
    dreego.Register("GET", "/login", a.handleLogin)
    dreego.Register("POST", "/logout", a.handleLogout)
}
```

Because `dreego.Register` is idempotent (re-registering a `method`+`pattern` replaces the handler), a later-registered plugin can override a route deterministically.

## Layout

```
plugins/
├── sample/                 ← minimal example plugin
│   ├── go.mod              → module codeberg.org/dreego/dreego/plugins/sample
│   ├── sample.go           → implements dreego.Plugin or other core interfaces
│   └── README.md
├── auth/                   ← future official plugin
├── db/
└── ...
```

## Rules

1. **Core never imports a plugin.** This is the invariant that keeps the dependency graph clean.
2. **Plugins import Core.** They use `codeberg.org/dreego/dreego/core`.
3. **Plugins with external deps get their own `go.mod`.** Plugins without external deps can live as plain packages in the root module or in `plugins/` with or without their own module.
4. **One repo, many modules.** Releases use directory-prefix tags (e.g. `plugins/auth/v0.0.1`) if a plugin has its own `go.mod`.

## Planned Plugins

| Plugin | Description | Dependencies |
|--------|-------------|-------------|
| `auth` | OAuth2, JWT, session management | `golang.org/x/oauth2` |
| `metrics` | Prometheus `/metrics` endpoint | `prometheus/client_golang` |
| `tracing` | OpenTelemetry spans | `go.opentelemetry.io/otel` |
| `cache` | Redis + in-memory cache | `redis/go-redis` |
| `storage` | S3, R2, local files | `aws/aws-sdk-go-v2` |
| `mail` | SMTP, Resend, Postmark | `resend/resend-go` |
| `db` | SQL driver registration + migrations | `golang-migrate/migrate` |
| `i18n` | Translation loading + templating | `nicksnyder/go-i18n` |
| `seo` | Meta tags, sitemap, JSON-LD | none |
| `markdown` | Markdown rendering with frontmatter | `yuin/goldmark` |
| `search` | Full-text search (Meilisearch, Typesense) | `meilisearch/meilisearch-go` |
| `jobs` | Background jobs + cron | `robfig/cron` |
| `pwa` | Service worker generation, manifest | none |
| `pdf` | PDF generation from HTML | `chromedp/chromedp` |
| `charts` | Chart.js/Canvas server-side components | none |
| `icons` | Lucide/Heroicons as components | none |

Community plugins can also be published in separate repositories as templates, but the official plugins live here.
