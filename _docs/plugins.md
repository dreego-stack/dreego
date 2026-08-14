# Plugins

Dreego plugins are ordinary Go packages in separate repositories. Each plugin
owns its `go.mod`, dependencies, releases, tests, CI, and typed configuration.
Core never imports a plugin package.

## Registration

Before v1 there is no central `Plugin` interface. A plugin exposes an explicit
registration function bound to the owning App:

```go
package auth

type Options struct {
    LoginPath  string
    CookieName string
}

func Register(app *dreego.App, options Options) error {
    if err := app.Use(sessionMiddleware(options)); err != nil {
        return err
    }
    if err := app.Register(http.MethodGet, options.LoginPath, loginHandler); err != nil {
        return err
    }
    return app.Register(http.MethodPost, options.LoginPath, authenticateHandler)
}
```

The application calls it before the App is built:

```go
app := dreego.New()
if err := auth.Register(app, auth.Options{
    LoginPath:  "/login",
    CookieName: "session",
}); err != nil {
    log.Fatal(err)
}
```

Registration order is source order. Duplicate routes fail instead of silently
overriding another handler. Registration after `Build`, `Handler`, `ServeHTTP`,
or `Listen` returns `dreego.ErrAppBuilt`.

## Core boundary

Core contains the SSR capabilities required by a normal Dreego application.
Optional capabilities, provider integrations, SSE, and WebSockets live in
separate plugin repositories, even when an implementation currently needs only
the standard library.

A provider-neutral interface is added to Core only after at least two real
implementations demonstrate the same small contract. Assets and lifecycle hooks
remain plugin-owned until real plugins prove that an App-level contract is
necessary. Plugin contracts remain provisional until v1.

## Repository layout

```text
github.com/dreego-stack/
├── dreego/
├── plugin-auth/
├── plugin-sse/
└── plugin-websocket/
```

## Rules

1. Core never imports a plugin.
2. Plugins import Core and register against a specific App.
3. Every optional plugin has its own repository and module.
4. Plugins use typed options rather than unstructured configuration maps.
5. A plugin must not weaken defaults for unrelated routes.

See the non-binding roadmap and `_todo/plugins/` for possible future plugins.
