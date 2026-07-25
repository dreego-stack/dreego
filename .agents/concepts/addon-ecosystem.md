
---
type: Concept
title: "Addon/Plugin-Ökosystem"
description: "Go-basiertes Plugin-System mit Compile-Time Safety und Tree-Shaking"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---

## Design-Philosophie

Dreego-Addons sind Go-Packages, die das `dreego.Plugin`-Interface erfüllen. Keine dynamischen Plugins, keine Laufzeit-Magie — reines, kompilierungszeit-sicheres Go.

## Plugin-Interface

```go
package dreego

type Plugin interface {
    Name() string
    RegisterRoutes(app *App)
    Middlewares() []func(http.Handler) http.Handler
    Assets() *embed.FS
}
```

## Einklinkpunkte

Ein Addon kann sich an 5 Stellen im Framework einklinken:

1. **Middleware** — HTTP-Wrapper (z.B. Auth-Checks)
2. **Routes** — Neue Pfade registrieren (z.B. `/auth/login`)
3. **Assets** — CSS/JS/Images via `//go:embed`
4. **Transpiler** — Custom-Tags (z.B. `<dreego:map />`)
5. **Context** — Erweiterung des Request-Kontexts (z.B. `c.User()`)

## Beispiel: dreego-auth

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

## Nutzung in main.go

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

## Vorteile des Go-Addon-Systems

1. **Keine Versionierungs-Hölle:** `go.mod` löst Abhängigkeiten strikt auf
2. **Compile-Time Safety:** Build bricht bei Inkompatibilitäten ab
3. **Tree-Shaking:** Nicht genutzter Code wird vom Compiler entfernt
4. **Installation:** `go get github.com/dreego-ecosystem/dreego-auth` — ein Befehl

## Addon-Ideen (vollständig)

### Auth & Security
| Addon            | Beschreibung                                    |
|------------------|------------------------------------------------|
| dreego-auth       | Login, Register, Sessions, OAuth, Passkeys     |
| dreego-csrf       | CSRF-Schutz (falls nicht im Core)              |
| dreego-2fa        | Zwei-Faktor-Authentifizierung                  |

### UI & Komponenten
| Addon            | Beschreibung                                    |
|------------------|------------------------------------------------|
| dreego-ui         | Komponenten-Bibliothek (Shadcn-ähnlich)        |
| dreego-map        | MapLibre/Leaflet Integration                   |
| dreego-charts     | Diagramme (Chart.js Wrapper)                   |
| dreego-icons      | Icon-Library                                   |
| dreego-markdown   | Markdown-Rendering                             |

### Daten & Backend
| Addon            | Beschreibung                                    |
|------------------|------------------------------------------------|
| dreego-db         | DB-Integration (SQLite, Turso, PG)             |
| dreego-storage    | File-Uploads (S3, R2, local)                   |
| dreego-jobs       | Hintergrund-Jobs & Cron                        |
| dreego-search     | Volltextsuche (Bleve/Meilisearch)              |
| dreego-cache      | Caching (Redis, In-Memory)                     |

### Business
| Addon            | Beschreibung                                    |
|------------------|------------------------------------------------|
| dreego-stripe     | Stripe Payments & Webhooks                     |
| dreego-mail       | E-Mail-Versand mit .dreego-Templates            |
| dreego-pdf        | PDF-Generierung                                |
| dreego-i18n       | Mehrsprachigkeit                               |
| dreego-seo        | Meta-Tags, Sitemap, OpenGraph                  |
| dreego-analytics  | Privacy-friendly Analytics                     |

### DX & Tools
| Addon            | Beschreibung                                    |
|------------------|------------------------------------------------|
| dreego-admin      | Auto-generiertes Admin-Dashboard               |
| dreego-pwa        | Progressive Web App                            |
| dreego-sitemap    | Automatische Sitemap-Generierung               |
| dreego-devtools   | Debug-Toolbar (wie Laravel Debugbar)           |

## Transpiler-Hook für Custom-Tags

Addons können eigene HTML-Tags registrieren:

```html
<!-- In einer .dreego-Datei -->
<dreego:map lat="52.52" lng="13.40" />
```

Der Transpiler:
1. Findet `<dreego:map />`
2. Prüft `dreego.config.json` auf installierte Addons
3. Ersetzt durch Go-Code: `dreegomap.RenderMap(dreegomap.Props{Lat: 52.52, Lng: 13.40})`
