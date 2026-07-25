
---
type: Decision
title: Routing, Plugin-Routes & Komponenten-System
description: Routing, Plugin-Routen-Registrierung und Komponenten-System mit Namespace-Hierarchie
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Routing, Plugin-Routes & Komponenten-System

**Datum:** 24.07.2026 (aktualisiert nach Review)
**Status:** Akzeptiert
**Ersetzt:** [file-based-routing](file-based-routing.md) (aktualisiert)

## Kontext / Offene Fragen

1. Go-Package-System: Jedes Verzeichnis = ein Package. `dreego/routes/about/` ist Package `about`. Aktuell muss `main.go` jedes Route-Paket einzeln importieren (`demo/main.go:4-7`). Das skaliert nicht.

2. Plugin-Routes: `dreego generate` scannt `"."` — findet nie externe Packages im Module-Cache oder Vendor. Wie kommen Plugin-Routen ins Binary?

3. Komponenten-Namespace: Wenn User `components/Button.dreego` hat UND `dreego-ui` auch `Button` anbietet — wie wird das aufgelost?

## Entscheidung 1: Generierte Route-Import-Datei statt manueller Imports

`dreego generate` erzeugt EINE zentrale Import-Datei, die alle Route-Pakete importiert. Der User importiert nur diese eine Datei.

```
dreego/
├── routes/                    ← User schreibt hier
│   ├── get.dreego           → dreego → init() registriert GET /
│   ├── about/get.dreego       → dreego → init() registriert GET /about
│   ├── users/[id]/get.dreego  → dreego → init() registriert GET /users/{id}
│   └── ...
├── gen/                       ← GENERIERT (wird committed)
│   └── routes.go              → importiert ALLE Route-Pakete
│
main.go importiert NUR `_ "myapp/dreego/gen"`
```

```go
// dreego/gen/routes.go  (GENERIERT)
package gen

import (
    _ "myapp/dreego/routes"          // index, about, ...
    _ "myapp/dreego/routes/about"    // falls about/ ein Unterordner ist
    _ "myapp/dreego/routes/users/_id_"
    _ "github.com/dreego/dreego-auth" // Plugin mit init()-Registrierung
)
```

Jedes Route-Paket (auch Plugin-Pakete) registriert sich via `init()` → `runtime.Register()`. Das `gen/routes.go` importiert alle — `main.go` importiert nur `gen`.

### Warum `init()` beibehalten?

- Go-idiomatisch: `database/sql` Treiber machen das genauso
- Plugin-Pakete brauchen keine spezielle Behandlung
- Kein Laufzeit-Scanning, kein Reflection
- `go build` linkt nur importierte Pakete → Tree-Shaking

## Entscheidung 2: Plugin-Routen via init() — Kein dreego generate nötig

Der Plugin-Author committed generierte `dree.go`-Dateien IM Plugin-Repo.

```
dreego-auth/
├── routes/
│   ├── login.go          ← pre-generiert (enthält init() + runtime.Register)
│   └── ...
├── go.mod
```

Plugin-Entwickler-Workflow:
```bash
cd dreego-auth
dreego generate               # erzeugt routes/*.go
git add routes/*.go && git commit
git tag v0.1.0 && git push
```

User-Workflow:
```bash
go get github.com/dreego/dreego-auth@v0.1.0
# dreego generate fügt den Import in gen/routes.go ein:
#   _ "github.com/dreego/dreego-auth"
```

`dreego generate` erkennt Plugin-Pakete über `go.mod` + `go list -m -json all` und fügt deren Import in `gen/routes.go` ein. Kein Scannen von Vendor/Modul-Cache nötig — `go build` holt automatisch die richtige Version.

## Entscheidung 3: Routing-Konventionen

| Syntax                | Pfad                    | Go-Param              |
|-----------------------|-------------------------|-----------------------|
| `get.dreego`        | `/`                     | —                     |
| `about.dreego`        | `/about`                | —                     |
| `[id]/get.dreego`     | `/users/{id}`           | `c.Param("id")`       |
| `[...catchall].dreego`| `/blog/{catchall}`      | `c.Param("catchall")` |
| `[[lang]]/get.dreego` | `/docs/{lang}` (optional)| `c.Param("lang")`    |
| `(auth)/login.dreego` | `/login` (Gruppe im Pfad ignoriert) |           |

Priorität: Statisch > Dynamisch > Optional > Catch-All

Konflikt-Erkennung: `dreego generate` wirft Error, wenn zwei Routen dasselbe Pattern beanspruchen:
```
error: route conflict: /auth/login
  dreego/routes/auth/login.dreego
  plugin: dreego-auth (github.com/dreego/dreego-auth)
```

### API-Routen & HTTP-Methoden

HTTP-Methode wird aus dem Dateinamen abgeleitet:

```
routes/api/users.get.dreego   → GET  /api/users
routes/api/users.post.dreego  → POST /api/users
routes/api/users.put.dreego   → PUT  /api/users
routes/api/users.delete.dreego→ DELETE /api/users
```

Alternativ via `<go method="post">` in der `.dreego`-Datei (wie aktuell).

API-Routen rendern KEIN Layout — nur das `<div>`-Fragment. Erkennung: Pfad enthält `api/` → `layout = nil`.

### Middleware pro Route (V1)

Pro Ordner eine `_middleware.go` (wird NICHT generiert — der User schreibt sie):

```go
// dreego/routes/admin/_middleware.go
package admin

import "github.com/dreego/dreego-auth"

func init() {
    runtime.RegisterMiddleware("/admin/", auth.RequireRole("admin"))
}
```

### Redirects & Rewrites

```json
// dreego.config.json
{
  "redirects": [
    { "from": "/old-blog", "to": "/blog", "status": 301 }
  ],
  "rewrites": [
    { "from": "/api/v1/*", "to": "/api/v2/*" }
  ]
}
```

Wird in `gen/routes.go` als Middleware-Logik vor den File-based Routes generiert.

## Entscheidung 4: Komponenten-System

### Drei Quellen — eine Namespace-Hierarchie

```
Priorität (höchste zuerst):
1. dreego/components/Button.dreego     ← User-Komponente (shadowed Plugin)
2. dreego/layouts/default.dreego       ← Layouts (Spezialfall)
3. Plugin-Assets via fs.FS             ← @dreego-ui/Button
```

Explizite Disambiguierung:
```
{#use Button from "components/Button.dreego"}     ← explizit User
{#use Button from "@dreego-ui/Button"}            ← explizit Plugin
```

Ohne `from`-Angabe wird zuerst im User-Verzeichnis gesucht, dann in Plugins:
```
{#use Button}   ← sucht Button.dreego in components/, dann in plugins
```

### Wie findet dreego generate Plugin-Komponenten?

Plugins legen `.dreego`-Komponenten in einem bekannten Pfad ab:

```
dreego-ui/
├── components/               ← KONVENTION
│   ├── Button.dreego
│   ├── Card.dreego
│   └── Alert.dreego
├── dreego.go                 ← Plugin-Interface impl
├── go.mod
```

`dreego generate`:
1. Liest `go.mod` → findet `github.com/dreego/dreego-ui`
2. `go list -m -json github.com/dreego/dreego-ui` → `Dir: /home/.../pkg/mod/...`
3. Sucht `<Dir>/components/*.dreego`
4. Bei Vendor-Mode: `vendor/github.com/dreego/dreego-ui/components/*.dreego`

Kein Plugin-Loading, kein Reflection, kein Import im CLI-Binary. Nur Dateisystem-Scan.

### Komponenten-Verwendung im Template

```html
<!-- dreego/routes/get.dreego -->
{#use Button}              <!-- findet components/Button.dreego -->
{#use Alert}               <!-- findet @dreego-ui/Alert (kein User-Alert) -->

<div>
    <Alert type="success">Geschafft!</Alert>
    <Button class="primary" disabled={isLoading}>
        Absenden
    </Button>
</div>
```

### Komponenten-CodeGen

Der Transpiler erzeugt:
```go
// <Button class="primary" disabled={isLoading}>Absenden</Button>
// wird zu:
renderButton(c, ButtonProps{Class: "primary", Disabled: isLoading}, func(c *Context) string {
    return "Absenden"
})
```

- Props: Alle HTML-Attribute werden als `ComponentProps`-Struct gebündelt
- Children: Inhalt zwischen Tags als Closure (entspricht `{#slot}`)
- Self-closing: `<Alert type="info" />` → kein Children-Parameter
- Scoping: Jede Komponente hat eigenen Scope-Hash (kein CSS-Leaking)

### Plugin-Komponenten ohne .dreego-Source (V2)

Wenn ein Plugin aus Performance-Gründen nur Go-Funktionen bereitstellt (keine `.dreego`-Transpilation nötig):

```go
// dreego-ui registriert eine Named-Components-Map
func (p *UIPlugin) Components() map[string]ComponentFunc {
    return map[string]ComponentFunc{
        "Button": renderButton,
        "Card":   renderCard,
    }
}
```

`{#use Button from "@dreego-ui/Button"}` → direkter Funktionsaufruf, kein Transpiler-Pass nötig.

## Entscheidung 5: dreego/routes/ vs dreego/pages/

`dreego/routes/` bleibt. Der Name ist etabliert (SvelteKit, SolidStart, Next.js Pages Router). `/pages/` wäre genauso valide, aber wir bleiben bei der bestehenden Konvention.

Konfigurierbar via `dreego.config.json`:
```json
{ "routeDir": "dreego/pages" }
```

## Konsequenzen

- `dreego/gen/routes.go` wird generiert und committed
- Enthält ALLE Route-Imports (File-based + Plugin)
- `main.go` importiert nur `_ "myapp/dreego/gen"`
- Plugin-Entwickler committen `dree.go`-Dateien (Pre-Generated)
- Plugin-Komponenten liegen unter `<plugin>/components/` (Konvention)
- `dreego generate` findet sie via `go list` + Dateisystem
- User-Komponenten shadowen Plugin-Komponenten (expliziter Namespace fallback)
- `[...catchall]`, `[[optional]]`, `(group)/` werden im Lexer/Parser ergänzt
- Duplicate-Route-Erkennung wirft Build-Fehler
