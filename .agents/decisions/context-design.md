
---
type: Decision
title: Context-Design (Interface + Target-Structs)
description: Schlankes Context-Interface plus drei konkrete Structs pro Rendering-Target
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Context-Design (Interface + Target-Structs)

**Datum:** 23.07.2026
**Status:** Akzeptiert
**Review:** GLM-5.2 Expert Review (.tmp/output1.md)

## Kontext

Der `<go>`-Block in `.dreego`-Dateien braucht einen Context für Datenzugriff. Dreego unterstützt drei Rendering-Targets:

1. **SSR** — `*http.Request` + `http.ResponseWriter`
2. **SSG** — Build-Zeit, kein HTTP
3. **Wails v3** — WebView + In-Memory IPC, kein HTTP

## Problem

Der `<go>`-Block wird vom Transpiler generiert — das Target ist zur Compile-Zeit bekannt. Der Entwickler soll eine konsistente API haben, aber die Targets unterscheiden sich fundamental.

## Entscheidung

**Schlankes Interface fürs gemeinsame Minimum + drei konkrete Structs pro Target.**

```go
// Core-Interface — NUR was alle drei Targets brauchen
type Context interface {
    context.Context                    // Deadline/Cancel Propagation
    Param(name string) string          // URL-Params / Frontmatter / Route-Params
    Data(key string) any               // Meta-Daten aus Frontmatter / Build-Info
    Render(name string, data any) error // Verschachteltes Rendering
}
```

```go
// Drei konkrete Implementierungen
type SSRContext struct {
    R *http.Request
    W http.ResponseWriter
}
func (c *SSRContext) User() User { ... }          // NUR SSR
func (c *SSRContext) Session() Session { ... }     // NUR SSR
func (c *SSRContext) FormValue(k string) string { ... } // NUR SSR
func (c *SSRContext) Redirect(url string) { ... }  // NUR SSR

type SSGContext struct {
    Meta     Frontmatter
    OutputDir string
}
func (c *SSGContext) Frontmatter() Frontmatter { ... } // NUR SSG

type WailsContext struct {
    App    *app.App
    Window *window.Window
}
func (c *WailsContext) SystemAPI() SystemAPI { ... } // NUR Wails
```

## Warum nicht die Alternativen

| Option                    | Verworfen weil...                                     |
|---------------------------|------------------------------------------------------|
| B: Ein Struct mit nil-Feldern | Mischt Concerns, Runtime-nil-Crash-Gefahr            |
| C: Go's context.Context   | Kein API-Vertrag, Type-Assertion-Hölle                |
| Generics                  | Nur 3 fixe Typen, keine n Targets. Interface reicht   |

## Plugin-Erweiterung

Pro Plugin ein Context-Wrapper mit Embedding + typisiertem Getter:

```go
// Im Plugin (dreego-auth)
type AuthPlugin struct { dreego.Plugin }
type AuthContext struct {
    dreego.Context
    user User
}

// Nutzung im <go>-Block
c := dreego.Get[AuthContext](ctx)
user := c.User()
```

Kein `c.Get("auth")` mit `any`-Cast. Kein String-Key. Typsicher.

## Begründung

1. **Compile-Time-Safety:** Ruft Entwickler `c.Session()` in SSG → Compile-Error. Kein Runtime-Panic.
2. **Der Codegen-Pass kennt das Target:** Generierter Code ruft direkt auf dem konkreten Typ auf, nicht durchs Interface.
3. **Das Interface ist Dokumentation & Test-Mock**, nicht Dispatch-Mechanismus.
4. **Go-idiomatisch:** Wie `database/sql` — ein Interface, verschiedene Driver-Impls.

## Konsequenzen

- Drei Codegen-Passes: SSR, SSG, Wails
- Jeder Pass generiert Code gegen den konkreten Context-Typ
- Target-spezifische Methoden erscheinen nur in ihrem Pass
- Interface für gemeinsame Tests nutzbar
- Plugins wrappen via Embedding, keine zentrale Registry nötig
