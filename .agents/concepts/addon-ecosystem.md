
---
type: Concept
title: "Addon/Plugin Ecosystem"
description: "Go-based plugin system with compile-time safety and tree-shaking"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

## Design Philosophy

Dreego addons are Go packages that fulfill the `dreego.Plugin` interface. No dynamic plugins, no runtime magic — pure, compile-time-safe Go.

## Plugin Interface

```go
package dreego

type Plugin interface {
    Name() string
    RegisterRoutes(app *App)
    Middlewares() []func(http.Handler) http.Handler
    Assets() *embed.FS
}
```

## Extension Points

An addon can hook into the framework at 5 points:

1. **Middleware** — HTTP wrappers (e.g. auth checks)
2. **Routes** — Register new paths (e.g. `/auth/login`)
3. **Assets** — CSS/JS/Images via `//go:embed`
4. **Transpiler** — Custom tags (e.g. `<dreego:map />`)
5. **Context** — Extension of the request context (e.g. `c.User()`)

## Example: dreego-auth

```go
package auth

import "github.com/.../dreego"

//go:embed assets/*
var authAssets embed.FS

type AuthPlugin struct {
    SecretKey string
}

func New(secret string) *AuthPlugin {
    return &AuthPlugin{SecretKey: secret}
}

func (p *AuthPlugin) Name() string { return "DreegoAuth" }

func (p *AuthPlugin) RegisterRoutes(app *dreego.App) {
    app.POST("/api/auth/login", p.handleLogin)
    app.GET("/auth/login", p.renderLoginPage)
}

func (p *AuthPlugin) Assets() *embed.FS {
    return &authAssets
}
```

## Usage in main.go

```go
import (
    "github.com/.../dreego"
    "github.com/dreego-ecosystem/dreego-auth"
)

func main() {
    app := dreego.New()
    app.UsePlugin(auth.New("super-secret-key"))
    app.Listen(":8080")
}
```

## Advantages of the Go Addon System

1. **No dependency hell:** `go.mod` resolves dependencies strictly
2. **Compile-Time Safety:** Build breaks on incompatibilities
3. **Tree-Shaking:** Unused code is removed by the compiler
4. **Installation:** `go get github.com/dreego-ecosystem/dreego-auth` — one command

## Addon Ideas (complete)

### Auth & Security
| Addon            | Description                                      |
|------------------|--------------------------------------------------|
| dreego-auth       | Login, Register, Sessions, OAuth, Passkeys       |
| dreego-csrf       | CSRF protection (if not in core)                 |
| dreego-2fa        | Two-factor authentication                        |

### UI & Components
| Addon            | Description                                      |
|------------------|--------------------------------------------------|
| dreego-ui         | Component library (Shadcn-like)                  |
| dreego-map        | MapLibre/Leaflet Integration                     |
| dreego-charts     | Charts (Chart.js wrapper)                        |
| dreego-icons      | Icon Library                                    |
| dreego-markdown   | Markdown Rendering                               |

### Data & Backend
| Addon            | Description                                      |
|------------------|--------------------------------------------------|
| dreego-db         | DB Integration (SQLite, Turso, PG)               |
| dreego-storage    | File Uploads (S3, R2, local)                     |
| dreego-jobs       | Background Jobs & Cron                           |
| dreego-search     | Full-text search (Bleve/Meilisearch)             |
| dreego-cache      | Caching (Redis, In-Memory)                       |

### Business
| Addon            | Description                                      |
|------------------|--------------------------------------------------|
| dreego-stripe     | Stripe Payments & Webhooks                       |
| dreego-mail       | Email delivery with .dreego templates            |
| dreego-pdf        | PDF Generation                                  |
| dreego-i18n       | Multi-language support                          |
| dreego-seo        | Meta-Tags, Sitemap, OpenGraph                    |
| dreego-analytics  | Privacy-friendly Analytics                       |

### DX & Tools
| Addon            | Description                                      |
|------------------|--------------------------------------------------|
| dreego-admin      | Auto-generated admin dashboard                   |
| dreego-pwa        | Progressive Web App                              |
| dreego-sitemap    | Automatic sitemap generation                     |
| dreego-devtools   | Debug toolbar (like Laravel Debugbar)            |

## Transpiler Hook for Custom Tags

Addons can register their own HTML tags:

```html
<!-- In a .dreego file -->
<dreego:map lat="52.52" lng="13.40" />
```

The transpiler:
1. Finds `<dreego:map />`
2. Checks `dreego.config.json` for installed addons
3. Replaces with Go code: `dreegomap.RenderMap(dreegomap.Props{Lat: 52.52, Lng: 13.40})`
