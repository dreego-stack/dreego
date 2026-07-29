# Plugins

Dreego plugins are separate Go modules in their own repositories under the `codeberg.org/dreego` organization. Core never imports plugins — plugins implement Core interfaces.

## Architecture

```
codeberg.org/dreego/dreego      ← Core (no external deps beyond stdlib)
codeberg.org/dreego/auth         ← Plugin: OAuth2, JWT, sessions
codeberg.org/dreego/metrics       ← Plugin: Prometheus /metrics
codeberg.org/dreego/cache         ← Plugin: Redis/Memory caching
codeberg.org/dreego/storage       ← Plugin: S3, R2, local file storage
codeberg.org/dreego/mail          ← Plugin: SMTP, Resend, Postmark
codeberg.org/dreego/db            ← Plugin: SQL drivers, migrations
```

Each plugin has its own `go.mod` with dependencies isolated from Core.

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
// codeberg.org/dreego/metrics
package metrics

import "codeberg.org/dreego/dreego/core"

func init() {
    core.RegisterMetrics(&PrometheusProvider{})
}
```

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
