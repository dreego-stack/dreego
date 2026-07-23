# Entscheidung: Session-Management — Interface im Core, Store als Addon

**Datum:** 23.07.2026
**Status:** Akzeptiert

## Kontext

SSR-Apps brauchen Sessions: Request-Zuordnung, CSRF-Schutz, Flash-Messages, Warenkorb, Form-Wizards. Auth (Login, `c.User()`) ist ein ANDERES Thema und gehört in dreego-auth.

## Entscheidung

**Session-Interface im Core. Konkrete Stores als Addons.**

```go
// Core: dreego/session (Interface)
type Store interface {
    Get(r *http.Request, key string) (string, error)
    Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error
    Delete(w http.ResponseWriter, r *http.Request, key string) error
    Destroy(w http.ResponseWriter, r *http.Request) error
}
```

```go
// Core: Built-in Cookie-Store (Default, keine externe Dependency)
type CookieStore struct { ... }
func NewCookieStore(secret []byte) *CookieStore

// Addon: dreego-session-redis
// Addon: dreego-session-db (PostgreSQL/SQLite)
```

## dreego-auth baut darauf auf

```go
func (p *AuthPlugin) login(w, r) {
    session.Set(w, r, "user_id", user.ID, nil)
}

func (p *AuthPlugin) Middlewares() []func(http.Handler) http.Handler {
    return []func(http.Handler) http.Handler{
        func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w, r) {
                uid, _ := session.Get(r, "user_id")
                ctx := context.WithValue(r.Context(), "user_id", uid)
                next.ServeHTTP(w, r.WithContext(ctx))
            })
        },
    }
}
```

## Warum nicht nur Addon?

- CSRF-Schutz (Core) braucht Sessions
- Flash-Messages (Core) brauchen Sessions
- Ohne Session-Interface im Core müsste jedes Addon eigene Session-Logik mitbringen
- Cookie-Store ist 50 Zeilen Go — kein Grund, das auszulagern

## Warum nicht komplett im Core?

- Redis/DB-Stores sind optional — nicht jede App braucht sie
- Session-Interface ist stabil, Stores können beliebig wachsen
- Core bleibt schlank (Cookie-Store ist der einzige Built-in Store)

## Konsequenzen

- `dreego.New()` erzeugt automatisch einen Cookie-Store
- Apps können via `app.Use(session.NewRedisStore(...))` upgraden
- dreego-auth nutzt `session.Store` Interface — kein Vendor-Lock-in
- Alle Addons (nicht nur Auth) können Sessions nutzen
