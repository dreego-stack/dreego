# Migrating from v0.1 to v0.2

This guide covers the public API changes between v0.1 and v0.2. The v0.2
release moves the HTTP host into an explicit target package and adds a typed,
render-neutral component contract.

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

`app.Listen`, `app.Shutdown`, and `app.SetServerConfig` were removed. `App`
continues to own SSR routes, middleware, sessions, and its compiled handler,
but no longer owns a running `http.Server` or lifecycle state. For programmatic
shutdown or custom timeouts, create a host explicitly:

```go
config := ssr.DefaultServerConfig()
config.ShutdownTimeout = 5 * time.Second
host := ssr.New(app, config)

if err := host.Start(addr); err != nil {
	return err
}

// Start has bound the address and now serves in the background.

if err := host.Shutdown(ctx); err != nil {
	log.Print(err)
}
```

`Start` binds synchronously and installs signal handling before returning. It
then serves in the background. `Wait` blocks until that lifecycle finishes and
returns its final serve or shutdown error. `Shutdown` drains active requests,
waits for cleanup, and returns the same lifecycle result.

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

- `dreego.Render(component)` renders an already typed generated component
  without an HTTP server and returns `dreego.Result`.
- `dreego.NewContext()` returns `dreego.RenderContext`, which exposes only
  render-safe context and data operations.
- Generated components implement
  `Render(dreego.RenderContext) (dreego.Result, error)`.
- Pure generated GET page renderers use `dreego.RenderContext`; routes with
  request-dependent server sections continue to use `dreego.SSRContext`.
- Pure pages expose a typed `Page<Name>() dreego.Component` constructor.

Typed component props are supplied through the generated constructor:

```go
result, err := dreego.Render(components.Badge("New"))
html := result.HTML
```

`Result.HTML` contains the complete byte-identical renderer output, including
head markup and inline scoped styles. Separate metadata fields are intentionally
deferred until the pipeline can compose them without loss.

`RenderContext` does not expose request, response, query, cookie, session, or
redirect operations. Code requiring those capabilities remains in generated
SSR handlers and fails to compile in a component render function.

These are new in v0.2 and enable future SSG and Wails targets.

## Release and versioning

v0.2.0 is tagged via a `stage/*` merge using `version: minor`. The stage merge
into main is the deliberate release act.

Routing and template syntax remain compatible. Regenerate application code
after upgrading so components adopt the new `RenderContext` and `Result`
signatures.
