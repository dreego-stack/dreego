# Wails v3

Wails v3 support is experimental while the upstream v3 release remains beta.
Dreego pins `github.com/wailsapp/wails/v3` to `v3.0.0-beta.20` for the first
implementation slice. Applications must not use `@latest` or enable Wails
server mode as a substitute for the desktop host.

## Phase 1 boundary

The first host renders target-neutral pages registered on `dreego.App` and
passes the resulting document to a native Wails window. It does not start a
Dreego HTTP listener. Routes that require `SSRContext`, form actions, request
headers, cookies, sessions, or other HTTP capabilities are not registered for
desktop rendering and return `dreego.ErrRenderRouteNotFound`.

The Linux integration suite observes both the render path and a running native
window. It checks the complete Wails process tree for TCP listener descriptors
and traces system calls so that even a short-lived `listen` call fails the
suite. It also exercises literal WebView navigation and native history. The
virtual display is configured with its own TCP transport disabled.

## Assets and literal routes

The desktop host gives Wails an in-process asset handler and loads the initial
literal route through the native Wails URL scheme. The handler serves only
pages registered for target-neutral rendering and files registered with
`App.RegisterStatic`. It never falls back to the SSR handler, so middleware,
cookies, form actions, redirects, and other HTTP behavior do not cross into the
desktop host.

Literal links perform full-document WebView navigation and use the native
back-forward history. Query strings, encoded paths, traversal segments,
backslashes, unknown routes, and non-GET requests are rejected. Dreego also
rejects `FRONTEND_DEVSERVER_URL` because an external development server would
violate the listener-free desktop contract.

Dynamic routes remain valid for SSR but are not supported by the Phase 1
desktop renderer. A matching desktop render attempt fails before window startup
with `dreego.ErrDynamicRenderRoute`, identifies the dynamic pattern, and directs
the developer to SSR or a literal desktop route.

```go
app := dreego.New()
if err := www.Register(app); err != nil {
    return err
}
host, err := wails.New(app)
if err != nil {
    return err
}
return host.Run(wails.Options{
    Name:   "My Dreego App",
    Title:  "My Dreego App",
    Path:   "/",
    Width:  800,
    Height: 600,
})
```

## Typed services and lifecycle

Desktop methods are available only when the application explicitly registers a
Wails service in `Options.Services`. Dreego validates every exported binding
method signature before creating the native application. Parameters and results
must use JSON-safe concrete types; channels, functions, interfaces, recursive
models, non-string map keys, and multiple non-error results fail with an
actionable startup diagnostic. Unexported struct fields are outside the client
contract and are ignored by JSON encoding.

```go
timer := NewTimerService(5 * 60)
return host.Run(wails.Options{
    Name: "Focus Timer",
    Services: []application.Service{
        application.NewService(timer),
    },
})
```

Dreego does not register filesystem, shell, network, clipboard, or window
services by default. Such capabilities exist only when the application adds a
service that exposes them. Services may implement Wails `ServiceStartup` and
`ServiceShutdown`; Wails cancels the application context before invoking
shutdown in reverse registration order. Dreego forwards the exact service list
without global registration.

Generate readable TypeScript declarations and browser-ready JavaScript from the
same explicit contract:

```sh
cd demo/demo-wailsv3
wails3 generate bindings -d app/bindings -ts -i -b ./app
wails3 generate bindings -d app/static/bindings -b -noevents ./app
```

The generated browser modules import `/wails/runtime.js`, which the native
Wails asset server supplies in process. Dreego embeds the remaining modules
from `static/`; no localhost listener or npm package is involved.

## Deterministic development cycle

Phase 1 intentionally uses restart-based reload instead of `wails3 dev` or
`FRONTEND_DEVSERVER_URL`. Stop the running application, run `dreego generate`,
regenerate bindings only when the Go service contract changed, and start the
application again. Every restart renders from registered source and embedded
assets, so it cannot depend on a stale external development server. Live bridge
updates remain deferred until Phase 2.

Wails v3 beta.20 can block both bridge calls and programmatic `Quit` during
headless Alpine GTK4 tests even after the application-started event. Phase 1
therefore verifies the native window, navigation, history, and no-listener
contract end to end, and verifies generated bindings plus Dreego's exact
service/lifecycle forwarding deterministically. Binding interaction remains a
manual supported-desktop check for the reference application; native headless
bridge and shutdown behavior are upstream-stable Phase 2 compatibility gates.

## Development toolchain

All repository commands run through `smd`. The image contains `pkgconf`, GTK 4,
and WebKitGTK 6.0 for the default Wails v3 Linux build. Install the pinned CLI
inside the container when needed:

```sh
smd go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.20
```

Do not run `wails3 setup` automatically. It may install host packages and is
not part of Dreego's reproducible build.

`wails3 doctor` in beta.20 panics while detecting the Alpine package manager.
The Dreego build therefore verifies the required packages directly. The panic
does not prevent compilation of the Wails host or timer demo.

See the [official Wails installation guide](https://v3.wails.io/getting-started/installation/),
[application API](https://v3.wails.io/reference/application/), and
[lifecycle documentation](https://v3.wails.io/concepts/lifecycle/). The
[binding method guide](https://v3.wails.io/features/bindings/methods/) describes
the generated client modules.
