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
suite. The virtual display is configured with its own TCP transport disabled.

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
[lifecycle documentation](https://v3.wails.io/concepts/lifecycle/).
