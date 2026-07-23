# Entscheidung: Middleware-System — Core vs Plugin

**Datum:** 23.07.2026
**Status:** Akzeptiert
**Review:** GLM-5.2 Expert Review (.tmp/output4.md)

## Entscheidung

### Core-Fixed (immer aktiv, nicht deaktivierbar)

```go
// internal/middleware/core.go
func coreMiddleware() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{
        Recovery(),       // Panic → 500 + log
        RequestID(),      // X-Request-ID Header
        RealIP(),         // Chi-built-in, hinter Proxy
        RequestLogging(), // slog, nutzt RequestID
    }
}
```

### Core-Conditional (default an, opt-out via Config)

```go
// Per dreego.config.json abschaltbar
CSRF()        // Braucht nur Session-Interface (Core)
CORS()        // Default: same-origin restriktiv
Compress()    // gzip/deflate via Chi
```

### Plugin/Addon (niemals Core)

| Middleware         | Warum Plugin                                   |
|--------------------|------------------------------------------------|
| Rate-Limiting/DDoS | Braucht Redis/Backend → Infrastruktur-abhängig |
| Auth (OAuth/JWT)   | Policy, nicht jeder braucht Auth               |
| Prometheus/Metrics | Monitoring-Infrastruktur                       |
| Tracing (OTel)     | Externer Collector nötig                       |

## Reihenfolge-Modell

```
[Recovery → RequestID → RealIP → Logging]   ← Core-fixed
  → [CORS → CSRF → Compress]                 ← Core-conditionals
    → [Plugin-Middleware A, B, C…]           ← app.Use() FIFO
      → [Router / Handler]
```

**V1:** Reine Registrierungsreihenfolge (FIFO via `app.Use()`)
**V2:** Constraint-Sortierung (`before`/`after` im Plugin-Manifest)

Middleware-Stack wird beim ersten `ListenAndServe` gelockt — späteres `app.Use()` → Panic (deterministisch).

## Middleware-Signatur

```go
type Middleware func(ctx Context, next http.Handler) http.Handler
```

Nimmt `dreego.Context`, nicht rohes `*http.Request` — target-agnostisch (gilt für SSR, SSG, Wails).

## Konsequenzen

- Core kompiliert OHNE externe Dienste — deterministisch, testbar
- CSRF ist Core weil es nur Session braucht (keine externe Infrastruktur)
- DDoS/Rate-Limiting ist Plugin — braucht Redis/Backend
- Plugins via `MiddlewareProvider` registrieren Middleware in `app.Use()`-Reihenfolge
