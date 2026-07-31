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

## Plugin Interface (planned for v0.1.0)

Core defines Go interfaces. Plugins register implementations.

```go
// Defined in core
type MetricsProvider interface {
    Middleware() func(http.Handler) http.Handler
    Handler() http.Handler
}

type CacheProvider interface {
    Get(key string) (any, bool)
    Set(key string, value any, ttl time.Duration)
    Delete(key string)
}
```

Plugins register via `init()`:

```go
// codeberg.org/dreego/dreego/plugins/auth
package auth

import "codeberg.org/dreego/dreego/core"

func init() {
    core.RegisterAuth(&OAuth2Provider{})
}
```

## Layout

```
plugins/
├── sample/                 ← minimal example plugin
│   ├── go.mod              → module codeberg.org/dreego/dreego/plugins/sample
│   ├── sample.go           → implements core.Plugin or other core interfaces
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
