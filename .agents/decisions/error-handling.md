# Entscheidung: Error-Handling-Strategie (Core, nicht Plugin)

**Datum:** 24.07.2026
**Status:** Akzeptiert

## Kontext

Dreego ist ein Compile-Time Transpiler. Fehler können auf drei Ebenen auftreten:

1. **Build-Zeit** (`dreego generate`, `go build`) — Template-Syntax, fehlende Variablen, Typ-Fehler
2. **Start-Zeit** (`main()`) — Fehlende Config, ungültige Plugins, Port-Konflikt
3. **Laufzeit** (pro Request) — Validierung, DB-Fehler, Panics, Auth-Fehler

Das Framework muss typensicher sein. Fehler sollen so früh wie möglich (Build > Start > Runtime) und so lokal wie möglich (pro Handler, nicht global) auftreten.

`net/http` bietet nur `http.Error(w, msg, code)`. Es gibt keine Recovery, kein strukturiertes Logging, keine typisierten Fehler, keine Flash/Redirect-Semantik.

## Entscheidung

**Typisierte Fehler-Typen im Core.** Recovery-Middleware dispatcht auf Typ, Renderer entscheidet HTML/JSON/Redirect.

### Fehler-Typen (neue Datei `pkg/errors/errors.go`)

```go
package errors

type HTTPError struct {
    Code   int
    Public string  // sichtbar für Endnutzer
    Cause  error   // interner Fehler (Dev-Modus: Stack, Prod: nur Log)
}

type ValidationError struct {
    Fields map[string]string  // Feldname → Fehlermeldung
}

type RedirectError struct {
    URL   string
    Flash FlashBag
}
```

### Error-Level-Pyramide (Build > Start > Runtime)

```
Build-Zeit:  dreego generate bricht ab, go build schlägt fehl
             → Template-Syntax, fehlende Variablen, ungültige Sektionen
             → Kein Output, non-zero exit

Start-Zeit:  server.Listen() bricht ab
             → Fehlende Config, Port belegt, Plugin-Init-Fehler
             → log.Fatal / slog.Error + os.Exit(1)

Laufzeit:    Pro-Request, Recovery fängt alles
             → Panic → Recovery → 500 HTML
             → ValidationError → 422 HTML/JSON
             → RedirectError → 302 + Flash-Cookie
             → HTTPError → entsprechender Status-Code
             → Unerwartet → 500 (Dev: Stack-Trace, Prod: generisch)
```

### Recovery-Middleware (Core-Fixed)

```go
// pkg/middleware/recovery.go
func Recovery(log *slog.Logger) func(http.Handler) http.Handler
```

Ablauf:
1. `defer recover()` — fängt alle Panics
2. Loggt via `slog.Error` mit RequestID
3. Dev-Modus: Rendert Fehlerseite mit Stack-Trace
4. Prod-Modus: Rendert generische `_error.dreego` oder fallback `500.html`

### Middleware-Stapel (Core-Fixed + Core-Conditional, Reihenfolge fix)

```
[Recovery → RequestID → RealIP → RequestLogging*]
  → [User-Middleware / Plugin-Middleware]
    → Router → Handler
```

\* `RequestLogging` ist Core-Conditional: default an, abschaltbar via `dreego/config.json` (`logging.enabled: false`). Ein Plugin (`dreego-logging`) kann es in V2 ersetzen.

- `Recovery` — Panic → 500 + log
- `RequestID` — X-Request-ID Header, in Context und Logs
- `RealIP` — X-Forwarded-For / X-Real-IP Auswertung
- `RequestLogging` — slog.Info pro Request (Methode, Pfad, Status, Dauer)

### Dev vs Prod

```go
func IsDev() bool { return os.Getenv("APP_ENV") != "production" }
```

- Dev: Stack-Traces im Browser, ausführliche Fehlermeldungen, `slog.LevelDebug`
- Prod: Generische 500-Seite, `slog.LevelInfo`, keine internen Details

### Form-Validierungs-Feedback

CodeGen erzeugt Handler, die `ValidationError` zurückgeben. Rendering:

```html
<input name="email" value="{c.Old("email")}" />
{#if c.Errors("email")}
  <p class="error">{c.Errors("email")}</p>
{/if}

{#if c.Flash("success")}
  <div class="success">{c.Flash("success")}</div>
{/if}
```

`c.Old()`, `c.Errors()`, `c.Flash()` sind Context-Methoden (nicht Magic-Variablen).

### Error Boundary (pro Komponente)

Keine Error-Boundary auf Komponenten-Ebene in V1. Begründung:
- `<go>`-Block läuft vor dem Template-Rendering — alle Fehler sind vorher bekannt
- `{#if hasError}` deckt den Use-Case ab
- Error-Boundaries sind ein SPA-Konzept (React), nicht nötig bei SSR

→ Siehe [[no-catch-tag]]: Fehler via `{#if hasError}`, kein spezielles Tag.

### Logging-Strategie

```go
// pkg/middleware/logging.go
func RequestLogging(log *slog.Logger) func(http.Handler) http.Handler
```

- Nutzt Go's `log/slog` (strukturiert, Level-basiert)
- RequestID aus Context in jedem Log-Eintrag
- Log-Level per Config: `APP_LOG_LEVEL=debug|info|warn|error`
- Kein Plugin — Core-fixed, immer aktiv, nicht deaktivierbar

### Integration mit Plugin-System

Addons geben Fehler als typisierte Werte zurück:

```go
// Im Auth-Plugin
func (p *AuthPlugin) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user, err := p.validateToken(r)
        if err != nil {
            dreego.WriteError(w, r, dreego.NewHTTPError(401, "Nicht authentifiziert", err))
            return
        }
        ctx := context.WithValue(r.Context(), "user", user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

`dreego.WriteError(w, r, err)` dispatcht auf den Typ und rendert entsprechend — die Recovery-Middleware ist der zentrale Dispatcher.

## Konsequenzen

- Neue Packages: `pkg/errors/`, `pkg/middleware/` (existiert als Konzept, wird befüllt)
- Alle Handler (auch generierte) geben `error` zurück — Recovery dispatched
- `dreego.Context` wird um `Errors()`, `Old()`, `Flash()` erweitert
- `dreego generate` validiert Template-Syntax → Build-Fehler vor Laufzeit
- `slog` ist Core-Dependency (Go 1.21+ stdlib, Go 1.22+ genutzt)
- Kein Chi — alle Middleware selbst gebaut auf `net/http`
