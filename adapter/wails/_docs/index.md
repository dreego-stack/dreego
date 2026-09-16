# Wails adapter

`adapter/wails` exposes a listener-free `http.Handler` for Dreego render routes
and static assets. It does not create or run a Wails application.

The application owns Wails options, windows, services, bindings, lifecycle,
and `Run`. Pass the handler returned by `wails.New(app)` to Wails' asset
configuration.

The adapter accepts `GET` and `HEAD`. It rejects query strings, unsafe paths,
dynamic routes, SSR-only handlers, sessions, actions, and middleware semantics.
