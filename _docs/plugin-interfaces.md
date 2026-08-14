# Plugin Interfaces

Dreego's plugin system is based on Go interfaces. Every interface is a contract: the core defines, plugins implement.

Official plugins live in separate repos under `github.com/dreego-stack/`. Each plugin has its own `go.mod` and requires `github.com/dreego-stack/dreego`. Core never imports any plugin package.

## Core Interfaces (in `github.com/dreego-stack/dreego/core`)

The session Store is the only remaining core infrastructure contract. EventBus, Queue, KVStore, and Storage were removed from core before v0.1; optional infrastructure contracts must first be proven by real plugins, which own their APIs in their own repositories.

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

### Current Plugin Interface

The following fat interface exists in the released pre-v0.1 implementation.
It is deprecated design and will be removed by the App migration; it is not a
compatibility promise for plugin authors.

```go
type Plugin interface {
    Name() string
    RegisterRoutes()
    Middlewares() []func(http.Handler) http.Handler
    Assets() fs.FS
    OnStart(ctx context.Context) error
    OnShutdown(ctx context.Context) error
}
```

The accepted v0.1 direction has no required central plugin interface. A plugin
instead exposes `Register(app, typedOptions) error`. Only real plugins may later
justify small shared capability interfaces, and compatibility begins at v1.

### Proposed Middleware Hooks (superseded)

```go
type MiddlewareProvider interface {
    Middlewares() []func(http.Handler) http.Handler
}
```

This interface was exploratory and is not the accepted v0.1 contract. Plugins
register middleware directly on their owning App.

### Proposed Route Hooks (superseded)

```go
type RouteProvider interface {
    Routes() []Route
}
```

This interface was exploratory and is not the accepted v0.1 contract. Plugins
register routes directly on their owning App.

## Plugin Interfaces (not yet implemented)

### Email Interface

```go
type Mailer interface {
    Send(ctx context.Context, msg Message) error
}
```

Implementations: `github.com/dreego-stack/dreego/plugins/mail-smtp`, `github.com/dreego-stack/dreego/plugins/mail-resend`.

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
