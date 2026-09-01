# Migrating from v0.1 to v0.2

This guide covers the public API changes between v0.1 and v0.2. The v0.2
release moves the HTTP host into an explicit target package and adds
render-neutral rendering. Everything else stays compatible.

## HTTP host moved to core/ssr

The HTTP host now lives in `core/ssr`. Import it as:

```go
import ssr "github.com/dreego-stack/dreego/core/ssr"
```

Start applications with:

```go
ssr.Listen(app, addr)
ssr.DefaultAddr()
```

`app.Listen` still works but is no longer the documented entry point.

## Middleware, sessions, and forms moved

The following moved to `core/ssr`. Only the import path changed; signatures are
unchanged:

- `ssr.CSRF`
- `ssr.Compress`
- `ssr.Recovery`
- `ssr.RequestID`
- `ssr.RequestLogging`
- `ssr.MaxBodyReader`
- `ssr.WithStore`
- `ssr.BindForm`
- `ssr.ServerConfig`
- `ssr.DefaultServerConfig`
- `ssr.ErrServerRunning`

## Render-neutral rendering

v0.2 adds non-HTTP rendering:

- `dreego.Render(fn, props)` renders a component without an HTTP server.
- `dreego.NewContext()` creates a render-neutral context.

These are new in v0.2 and enable future SSG and Wails targets.

## Release and versioning

v0.2.0 is tagged via a `stage/*` merge using `version: minor`. The stage merge
into main is the deliberate release act.

## Nothing else changes

Generated code, the component API, routing, and templates stay compatible.
There are no other breaking changes in v0.2.
