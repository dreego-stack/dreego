# Plugin Interfaces

Dreego's plugin system is based on Go interfaces. Every interface is a contract: the core defines, plugins implement.

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

Every plugin implements `Init()` for registering routes, middleware, services.

### Middleware Hooks (Plugin)

```go
type MiddlewareProvider interface {
    Middlewares() []func(http.Handler) http.Handler
}
```

Plugins that inject HTTP middleware (CSRF, auth, rate-limiting, etc.).

### Route Hooks (Plugin)

```go
type RouteProvider interface {
    Routes() []Route
}
```

Plugins that register their own URL paths (`/admin/*`, `/api/auth/*`).

## Plugin Interfaces (not yet implemented)

### Storage Interface

```go
type Storage interface {
    Put(ctx context.Context, key string, r io.Reader) error
    Get(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    URL(ctx context.Context, key string) (string, error)
}
```

Implementations: `dreego-storage-s3`, `dreego-storage-local`.

### Email Interface

```go
type Mailer interface {
    Send(ctx context.Context, msg Message) error
}
```

Implementations: `dreego-mail-smtp`, `dreego-mail-resend`.

### Queue Interface

```go
type Queue interface {
    Dispatch(ctx context.Context, job Job) error
    Worker(name string, handler JobHandler)
}
```

Implementations: `dreego-jobs-redis`, `dreego-jobs-memory`.

### Cache Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

Implementations: `dreego-cache-redis`, `dreego-cache-memory`.

### Event Bus Interface

```go
type EventBus interface {
    Publish(ctx context.Context, topic string, event any) error
    Subscribe(topic string, handler EventHandler) error
}
```

Implementations: `dreego-eventbus-redis`, `dreego-eventbus-nats`.

## Plugin Repo Structure

```
codeberg.org/dreego/dreego              ← Core
codeberg.org/dreego/dreego-session-redis ← Plugin
codeberg.org/dreego/dreego-auth          ← Plugin
```

Or in your own project repo:

```
myapp/
├── go.mod           (module myapp)
├── main.go          (import _ "myapp/plugins/auth")
├── plugins/
│   └── auth/
│       └── auth.go  (implements dreego.Plugin)
└── dreego/
    └── routes/
```

Then import:
```go
import _ "myapp/plugins/auth"
```

## Cluster Plugin (planned)

`dreego-cluster` — distributed state for multi-node deployments. Combines:

```
dreego-cluster
├── Node-Discovery       (memberlist / DNS / static)
├── Shared-Session-Store (dreego-session-redis)
├── PubSub-Sync          (dreego-eventbus-redis)
└── Distributed-Cache    (dreego-cache-redis)
```

A load balancer distributes requests across N Go instances. `dreego-cluster` ensures all instances see the same state. No Kubernetes needed — Valkey/Redis as backend is sufficient.
