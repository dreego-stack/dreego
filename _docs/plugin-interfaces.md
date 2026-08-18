# Plugin Interfaces

Dreego plugins use explicit App-bound registration functions. Core interfaces
exist only for proven framework contracts such as sessions; plugins do not
implement a central Plugin interface.

> **Provisional until v1.** The plugin contract is explicitly excluded from the
> v0.1 stability promise. It is validated by real external plugins between v0.1
> and v1 and may change until then. See
> [Compatibility](https://github.com/dreego-stack/dreego/blob/main/_docs/compatibility.md).

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

### Plugin registration

There is no central Plugin interface before v1. Plugin packages expose a typed
`Register(app, options) error` function and use `app.Register` and `app.Use`
directly. Shared capability interfaces are introduced only after multiple real
plugins prove the same small contract.

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
├── main.go          (calls auth.Register)
├── plugins/
│   └── auth/
│       └── auth.go  (exports Register and Options)
└── dreego/
    └── routes/
```

Then register the feature explicitly on the owning app:
```go
app := dreego.New()
if err := auth.Register(app, auth.Options{
	LoginPath: "/login",
}); err != nil {
	log.Fatal(err)
}
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
