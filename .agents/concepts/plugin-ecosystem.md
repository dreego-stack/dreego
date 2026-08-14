---
type: Concept
title: "Plugin Ecosystem"
description: "Go-based plugin system with compile-time safety and tree-shaking"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

## Design Philosophy

Dreego plugins are ordinary Go packages with explicit App-bound registration.
There is no dynamic loading, runtime discovery, or central Plugin interface.

## Registration contract

```go
type Options struct {
    LoginPath string
}

func Register(app *dreego.App, options Options) error {
    if err := app.Use(authMiddleware(options)); err != nil {
        return err
    }
    return app.Register(http.MethodGet, options.LoginPath, loginHandler)
}
```

## Extension Points

A plugin registers only the App behavior it needs:

1. **Middleware** — HTTP wrappers (e.g. auth checks)
2. **Routes** — Register new paths (e.g. `/auth/login`)
3. **Context data** — request-local typed helpers owned by the plugin

Assets, lifecycle hooks, and transpiler extensions require separate proven
contracts. Runtime registration does not imply a transpiler extension API.

## Example: dreego-auth

```go
package auth

import "github.com/.../dreego"

type Options struct {
    SecretKey string
}

func Register(app *dreego.App, options Options) error {
    if err := app.Register(http.MethodPost, "/api/auth/login", handleLogin(options)); err != nil {
        return err
    }
    return app.Register(http.MethodGet, "/auth/login", renderLoginPage(options))
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
	if err := auth.Register(app, auth.Options{SecretKey: "super-secret-key"}); err != nil {
		log.Fatal(err)
	}
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
```

## Advantages of the Go Plugin System

1. **No dependency hell:** `go.mod` resolves dependencies strictly
2. **Compile-Time Safety:** Build breaks on incompatibilities
3. **Tree-Shaking:** Unused code is removed by the compiler
4. **Installation:** `go get github.com/dreego-ecosystem/dreego-auth` — one command

## Plugin Ideas (complete)

### Auth & Security
| Plugin           | Description                                      |
|------------------|--------------------------------------------------|
| dreego-auth       | Login, Register, Sessions, OAuth, Passkeys       |
| dreego-csrf       | CSRF protection (if not in core)                 |
| dreego-2fa        | Two-factor authentication                        |

### UI & Components
| Plugin           | Description                                      |
|------------------|--------------------------------------------------|
| dreego-ui         | Component library (Shadcn-like)                  |
| dreego-map        | MapLibre/Leaflet Integration                     |
| dreego-charts     | Charts (Chart.js wrapper)                        |
| dreego-icons      | Icon Library                                    |
| dreego-markdown   | Markdown Rendering                               |

### Data & Backend
| Plugin           | Description                                      |
|------------------|--------------------------------------------------|
| dreego-db         | DB Integration (SQLite, Turso, PG)               |
| dreego-storage    | File Uploads (S3, R2, local)                     |
| dreego-jobs       | Background Jobs & Cron                           |
| dreego-search     | Full-text search (Bleve/Meilisearch)             |
| dreego-cache      | Caching (Redis, In-Memory)                       |

### Business
| Plugin           | Description                                      |
|------------------|--------------------------------------------------|
| dreego-stripe     | Stripe Payments & Webhooks                       |
| dreego-mail       | Email delivery with .dreego templates            |
| dreego-pdf        | PDF Generation                                  |
| dreego-i18n       | Multi-language support                          |
| dreego-seo        | Meta-Tags, Sitemap, OpenGraph                    |
| dreego-analytics  | Privacy-friendly Analytics                       |

### DX & Tools
| Plugin           | Description                                      |
|------------------|--------------------------------------------------|
| dreego-admin      | Auto-generated admin dashboard                   |
| dreego-pwa        | Progressive Web App                              |
| dreego-sitemap    | Automatic sitemap generation                     |
| dreego-devtools   | Debug toolbar (like Laravel Debugbar)            |

## Transpiler extensions

Runtime plugin registration does not extend `.dreego` syntax. A future
transpiler extension must first prove a typed processor boundary, diagnostics,
conflict handling, and deterministic builds. It is not part of the v0.1 plugin
contract.
