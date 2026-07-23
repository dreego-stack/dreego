# Entscheidung: Plugin-Interface (Capability-basiert)

**Datum:** 23.07.2026
**Status:** Akzeptiert
**Review:** GLM-5.2 Expert Review (.tmp/output2.md)

## Kontext

Dreego braucht ein Plugin-System für sein Addon-Ökosystem. Das Interface muss von Anfang an stabil sein — spätere Änderungen sind Breaking Changes für alle Addons.

## Entscheidung

**Capability-basierte Interfaces.** Statt eines fetten Single-Interface implementieren Plugins nur die Fähigkeiten, die sie brauchen — über Go's implizite Interface-Satisfaction.

## Interface-Definition

```go
package dreego

import (
    "context"
    "io/fs"
    "net/http"
)

type PluginID string

// Base — JEDES Plugin MUSS das implementieren
type Plugin interface {
    ID() PluginID
    Init(app *App, cfg map[string]any) error
}

// SSR-Middleware injizieren (optional)
type MiddlewareProvider interface {
    Middlewares() []func(http.Handler) http.Handler
}

// Routen registrieren — SSR + SSG (optional)
type RouteRegistrar interface {
    RegisterRoutes(r Router) error
}

// Embedded Assets (CSS/JS/Images) — target-agnostisch via fs.FS (optional)
type AssetProvider interface {
    Assets() fs.FS
}

// Custom-Tags wie <dreego:map /> im Transpiler (optional)
type TranspilerHook interface {
    Namespace() string
    ParseTag(node TagNode) (TagRenderer, error)
}

// Request-Context erweitern (optional)
type ContextExtender interface {
    BindContext(b *ContextBuilder)
}

// Startup/Shutdown (optional)
type Lifecycle interface {
    OnStart(ctx context.Context) error
    OnShutdown(ctx context.Context) error
}
```

## Beispiel: dreego-auth

```go
package auth

import (
    "embed"
    "github.com/dreego/dreego"
)

//go:embed assets/*
var assets embed.FS

type AuthPlugin struct {
    secretKey string
}

func New(secretKey string) *AuthPlugin {
    return &AuthPlugin{secretKey: secretKey}
}

func (p *AuthPlugin) ID() dreego.PluginID { return "dreego.io/auth" }

func (p *AuthPlugin) Init(app *dreego.App, cfg map[string]any) error {
    p.secretKey = cfg["secret"].(string)
    return nil
}

// AuthPlugin implementiert MiddlewareProvider
func (p *AuthPlugin) Middlewares() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{p.authMiddleware}
}

// AuthPlugin implementiert RouteRegistrar
func (p *AuthPlugin) RegisterRoutes(r dreego.Router) error {
    r.Get("/auth/login", p.loginPage)
    r.Post("/api/auth/login", p.handleLogin)
    return nil
}

// AuthPlugin implementiert AssetProvider
func (p *AuthPlugin) Assets() fs.FS { return assets }

// AuthPlugin implementiert ContextExtender
func (p *AuthPlugin) BindContext(b *dreego.ContextBuilder) {
    b.Set(UserKey, func(c dreego.Context) *User {
        return c.Get(UserKey).(*User)
    })
}
```

## Plugin-Reihenfolge

Explizit durch Registrierungsreihenfolge — keine Magic:

```go
func main() {
    app := dreego.New()
    app.Use(session.New())   // 1. Session (MUSS vor Auth)
    app.Use(auth.New("key")) // 2. Auth
    app.Use(admin.New())     // 3. Admin
    app.Listen(":8080")
}
```

Middleware wird in Durchlaufreihenfolge gewrappt (LIFO wie Chi).

## Warum Capability-basiert

| Vorteil                                 | Erklärung                                     |
|-----------------------------------------|-----------------------------------------------|
| Plugins implementieren nur was sie brauchen | `dreego-map` braucht kein Middleware         |
| Neue Capabilities additiv hinzufügbar   | `Lifecycle` kam später — kein Breaking Change |
| Go-idiomatisch                          | `io.Reader`, `fs.FS` sind dasselbe Muster     |
| Keine leeren Methoden                   | Kein `return nil` für ungenutzte Features     |

## Warum `fs.FS` statt `embed.FS`

`fs.FS` ist Go's Standard-Interface für Dateisysteme. Es funktioniert target-agnostisch:
- SSR: `embed.FS` implementiert `fs.FS`
- SSG: `os.DirFS` implementiert `fs.FS` (Dateien auf Disk)
- Wails: Custom `fs.FS`-Implementierung möglich

## Konsequenzen

- `reego.Plugin` Interface im Core — nie mehr ändern
- Neue Capabilities als separate Interfaces hinzufügbar
- Kein zentrales Plugin-Registry nötig — Nutzer registriert explizit
- Addons starten mit `go get` + `app.Use()`
