# Plugin Interfaces

Dreego's plugin system is based on Go interfaces. Every interface is a contract: the core defines, plugins implement.

Official plugins live in separate repos under `github.com/dreego-stack/`. Each plugin has its own `go.mod` and requires `github.com/dreego-stack/dreego`. Core never imports any plugin package.

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

### Event Bus Interface

Typed pub/sub contract. Implementations may back it with in-memory storage,
Redis, NATS or similar; core code stays transport-agnostic.

```go
type EventBus[T any] interface {
    Publish(ctx context.Context, event T) error
    Subscribe(ctx context.Context, handler func(T)) (Subscription, error)
    Unsubscribe(sub Subscription)
}
```

`Subscription` is an opaque handle identifying a registered handler:

```go
type Subscription interface {
    ID() uint64
}
```

Built-in: `NewInMemoryBus[T]()`. Plugins: `github.com/dreego-stack/dreego/plugins/eventbus-redis`, `github.com/dreego-stack/dreego/plugins/eventbus-nats` (planned).

### Key-Value Store Interface

Like `database/sql`: core defines the contract, plugins provide the implementation (Redis, Ristretto, in-memory). Distinct from `Storage` (blobs) — KV holds small values with an optional TTL.

```go
type KVStore interface {
    Get(ctx context.Context, key string) ([]byte, error)
    Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Expire(ctx context.Context, key string, ttl time.Duration) error
}
```

Semantics:
- Get returns the value stored under key; an error if key does not exist or ttl expired.
- Set stores val under key with ttl; ttl <= 0 means no expiry (keep forever).
- Delete removes key; idempotent (no error for missing key).
- Expire sets/adjusts the ttl on an existing key; error if key does not exist.
- All methods respect ctx cancellation.

Plugins: `github.com/dreego-stack/dreego/plugins/kv-redis`, `github.com/dreego-stack/dreego/plugins/kv-memory` (planned).

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

## Plugin Layout

Official plugins live in separate repos under `github.com/dreego-stack/`:

```
github.com/dreego-stack/
├── dreego/              ← main repo (core + CLI, single module)
├── plugin-example/      ← minimal example
├── plugin-auth/
├── plugin-db/
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

`github.com/dreego-stack/plugin-cluster` — distributed state for multi-node deployments. Combines:

```
plugin-cluster
├── Node-Discovery       (memberlist / DNS / static)
├── Shared-Session-Store (plugin-session-redis)
├── PubSub-Sync          (plugin-eventbus-redis)
└── Distributed-Cache    (plugin-cache-redis)
```

A load balancer distributes requests across N Go instances. `plugin-cluster` ensures all instances see the same state. No Kubernetes needed — Valkey/Redis as backend is sufficient.
