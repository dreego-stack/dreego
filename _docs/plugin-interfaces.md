# Plugin-Interfaces

Dreegos Plugin-System basiert auf Go-Interfaces. Jedes Interface ist ein Vertrag: der Core definiert, Plugins implementieren.

## Core Interfaces (in `codeberg.org/dreego/dreego`)

### session.Store

```go
type Store interface {
    Get(r *http.Request, key string) (string, error)
    Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error
    Delete(w http.ResponseWriter, r *http.Request, key string) error
    Destroy(w http.ResponseWriter, r *http.Request) error
}
```

Built-in: `CookieStore`. Plugins: `dreego-session-redis`.

### Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Init(app *App) error
}
```

Jedes Plugin implementiert `Init()` zur Registrierung von Routes, Middleware, Services.

### Middleware-Hooks (Plugin)

```go
type MiddlewareProvider interface {
    Middlewares() []func(http.Handler) http.Handler
}
```

Plugins die HTTP-Middleware injecten (CSRF, Auth, Rate-Limiting, etc.).

### Route-Hooks (Plugin)

```go
type RouteProvider interface {
    Routes() []Route
}
```

Plugins die eigene URL-Pfade registrieren (`/admin/*`, `/api/auth/*`).

## Plugin-Interfaces (noch nicht implementiert)

### Storage Interface

```go
type Storage interface {
    Put(ctx context.Context, key string, r io.Reader) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    URL(ctx context.Context, key string) (string, error)
}
```

Implementierungen: `dreego-storage-s3`, `dreego-storage-local`.

### Email Interface

```go
type Mailer interface {
    Send(ctx context.Context, msg Message) error
}
```

Implementierungen: `dreego-mail-smtp`, `dreego-mail-resend`.

### Queue Interface

```go
type Queue interface {
    Dispatch(ctx context.Context, job Job) error
    Worker(name string, handler JobHandler)
}
```

Implementierungen: `dreego-jobs-redis`, `dreego-jobs-memory`.

### Cache Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

Implementierungen: `dreego-cache-redis`, `dreego-cache-memory`.

### Event-Bus Interface

```go
type EventBus interface {
    Publish(ctx context.Context, topic string, event any) error
    Subscribe(topic string, handler EventHandler) error
}
```

Implementierungen: `dreego-eventbus-redis`, `dreego-eventbus-nats`.

## Plugin-Repo-Struktur

```
codeberg.org/dreego/dreego              ← Core
codeberg.org/dreego/dreego-session-redis ← Plugin
codeberg.org/dreego/dreego-auth          ← Plugin
```

Oder im eigenen Projekt-Repo:

```
myapp/
├── go.mod           (module myapp)
├── main.go          (import _ "myapp/plugins/auth")
├── plugins/
│   └── auth/
│       └── auth.go  (implementiert dreego.Plugin)
└── dreego/
    └── routes/
```

Import dann:
```go
import _ "myapp/plugins/auth"
```

## Cluster-Plugin (geplant)

`dreego-cluster` — verteilter State für Multi-Node-Deployments. Kombiniert:

```
dreego-cluster
├── Node-Discovery       (memberlist / DNS / static)
├── Shared-Session-Store (dreego-session-redis)
├── PubSub-Sync          (dreego-eventbus-redis)
└── Distributed-Cache    (dreego-cache-redis)
```

Ein Loadbalancer verteilt Requests auf N Go-Instanzen. `dreego-cluster` sorgt dass alle Instanzen den gleichen State sehen. Kein Kubernetes nötig — reicht Valkey/Redis als Backend.
