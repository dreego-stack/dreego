# Plugin Interfaces

Dreego's plugin system is based on Go interfaces. Every interface is a contract: the core defines, plugins implement.

Official plugins live under `plugins/` in the dreego repository. Each plugin with external dependencies has its own `go.mod`, while dependency-free plugins can be plain packages in the root module. Core never imports any plugin package.

## Core Interfaces (in `github.com/dreego-stack/dreego/core`)

### session.Store

```go
type Store interface {
    Get(r *http.Request, key string) (string, error)
    Set(w http.ResponseWriter, r *http.Request, key, value string, opts *Options) error
    Delete(w http.ResponseWriter, r *http.Request, key string) error
    Destroy(w http.ResponseWriter, r *http.Request) error
}
```

Built-in: `CookieStore`. Plugin: `github.com/dreego-stack/dreego/plugins/session-redis` (planned).

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

Implementations: `github.com/dreego-stack/dreego/plugins/storage-s3`, `github.com/dreego-stack/dreego/plugins/storage-local`.

### Email Interface

```go
type Mailer interface {
    Send(ctx context.Context, msg Message) error
}
```

Implementations: `github.com/dreego-stack/dreego/plugins/mail-smtp`, `github.com/dreego-stack/dreego/plugins/mail-resend`.

### Queue Interface

```go
type Queue interface {
    Dispatch(ctx context.Context, job Job) error
    Worker(name string, handler JobHandler)
}
```

Implementations: `github.com/dreego-stack/dreego/plugins/jobs-redis`, `github.com/dreego-stack/dreego/plugins/jobs-memory`.

### Cache Interface

```go
type Cache interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

Implementations: `github.com/dreego-stack/dreego/plugins/cache-redis`, `github.com/dreego-stack/dreego/plugins/cache-memory`.

### Event Bus Interface

```go
type EventBus interface {
    Publish(ctx context.Context, topic string, event any) error
    Subscribe(topic string, handler EventHandler) error
}
```

Implementations: `github.com/dreego-stack/dreego/plugins/eventbus-redis`, `github.com/dreego-stack/dreego/plugins/eventbus-nats`.

## Plugin Layout

Official plugins live in the same repository:

```
dreego/
├── core/
├── cmd/dreego/
└── plugins/
    ├── sample/              ← minimal example
    ├── auth/
    ├── db/
    └── ...
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

`dreego/plugins/cluster` — distributed state for multi-node deployments. Combines:

```
dreego/plugins/cluster
├── Node-Discovery       (memberlist / DNS / static)
├── Shared-Session-Store (dreego/plugins/session-redis)
├── PubSub-Sync          (dreego/plugins/eventbus-redis)
└── Distributed-Cache    (dreego/plugins/cache-redis)
```

A load balancer distributes requests across N Go instances. `dreego/plugins/cluster` ensures all instances see the same state. No Kubernetes needed — Valkey/Redis as backend is sufficient.
