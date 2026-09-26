# Recipe: Machine Endpoints (Webhooks)

A **machine endpoint** is a route that a third party calls without a browser:
a webhook, a signed callback, a service-to-service POST. It has no CSRF token,
no session cookie, and no HTML form.

The problem is that Dreego's default middleware stack protects state-changing
requests with CSRF once a session store is configured (globally via
`app.SetSessionStore` or per profile). A `POST /hooks/github` from a webhook
sender arrives with no token and is rejected with `403 invalid csrf token`.

Two patterns solve this. Pattern B is the accepted approach from v0.10.9; keep
Pattern A only for older releases.

## Why CSRF fires

Core middleware order is:

```
[Recovery → SecurityHeaders → Compression → RequestLogging]
  → Session → CSRF
    → Redirect/Rewrite
      → Router (mux)
        → your handlers
```

Middleware registered with `app.Use` wraps the session/CSRF stack, so it runs
**before** CSRF (this is also why an app-wide `MaxBodyReader` protects form
posts before CSRF parses them). That ordering is the basis of Pattern A.

## Pattern A: intercept before CSRF (pre-0.10.9)

Register a middleware that matches the webhook path and handles it directly.
Because user middleware runs before session and CSRF, a matching request never
reaches the CSRF check.

```go
app.Use(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/hooks/") {
			handleWebhook(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
})
```

The route is then **not** declared in `www/routes`; the middleware owns it.
This works, but it splits the endpoint away from the route tree: there is no
`+page.dreego`, no generated handler, and the path is matched by hand.

## Pattern B: `PROFILE` with CSRF disabled (v0.10.9+)

A **profile** scopes session and CSRF middleware to a set of routes. Declare the
profile once in `main.go`, then bind the route folder to it with the `PROFILE`
header directive.

A profile with no `Session` and `CSRF: false` adds no session/CSRF handling when
no global store is configured — exactly what a webhook needs. `CSRF: false`
overrides the App default for the profile's routes. A profile `Session` overrides
the global `SetSessionStore` for its routes; when the profile leaves `Session`
nil and a global store exists, the profile's CSRF setting still applies but the
global store is used.

```go
app := dreego.New()

disabled := false
if err := app.Profile("hooks", dreego.Profile{CSRF: &disabled}); err != nil {
	log.Fatal(err)
}
```

```text
GOIMPORT { io, myapp/internal/webhooks }

PROFILE "hooks"

<server method="post">
	body, err := io.ReadAll(io.LimitReader(c.R.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if err := webhooks.VerifySignature(c.R, body); err != nil {
		c.Write(http.StatusUnauthorized, "text/plain", []byte("bad signature"))
		return "", nil
	}
	c.JSON(http.StatusOK, map[string]string{"status": "accepted"})
	return "", nil
</server>
```

Place the file in its own folder, for example
`www/routes/(hooks)/github/+page.dreego`. The `(hooks)` group directory does not
appear in the URL but carries the `PROFILE` for every descendant route folder.
`dreego generate` emits `app.ApplyProfile("/github", "hooks")`, and the
compiler's profile resolution finds the nearest ancestor `PROFILE`.

Rules for profiles:

- The `PROFILE` name must match a profile registered with `app.Profile`. An
  unknown name fails at build time with a diagnostic naming the profile.
- A profile applies to its route folder and all descendant route folders;
  the nearest ancestor declaration wins.
- A route folder must not declare two conflicting `PROFILE` values.
- When at least one `PROFILE` is bound, session and CSRF become
  **per-profile**: only profiled routes receive the session/CSRF stack, and
  unprofiled routes (for example marketing pages) get no session cookie at all.
  A profile that leaves `Session` nil falls back to a global
  `app.SetSessionStore`; a profile-level `CSRF` value overrides the App default
  for that profile's routes.
- With **no** `PROFILE` at all, the App behaves exactly as before: a global
  `app.SetSessionStore` applies to every route and CSRF stays global.

## Authenticate machine requests yourself

CSRF is a browser-origin defense; it does not authenticate a webhook. Verify
the sender explicitly:

- **HMAC signature** over the raw body with a shared secret, compared with
  `hmac.Equal` (constant time).
- **Bearer token** or a signed header from the provider.
- **Source allowlist** only as an additional check, never as the sole one.

Read the body exactly once, verify it, then parse it. Cap the size with
`io.LimitReader` or a `core.MaxBodyReader` middleware so a webhook cannot stream
an unbounded body.

## Custom 403 responses

If a profile keeps CSRF enabled and you want an HTML or JSON error page instead
of the plaintext default, configure the App error handler:

```go
app.SetErrorHandler(http.StatusForbidden, func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"error":"forbidden"}`))
})
```

The CSRF middleware routes its `403` through this handler when it is set, and
falling back to the plaintext `invalid csrf token` body when it is not.
Configure it before the App is built.

## Responding from a machine endpoint

A machine endpoint needs explicit control over the response, because the caller
does not negotiate content. Prefer `type="custom"` with the response helpers and
an explicit `return "", nil`:

```text
<server type="custom" method="post">
	if err := handleDelivery(c); err != nil {
		c.Write(http.StatusBadRequest, "text/plain", []byte("rejected"))
		return "", nil
	}
	c.JSON(http.StatusAccepted, map[string]any{"queued": true})
	return "", nil
</server>
```

- `type="custom"` runs unconditionally; there is no `Accept` negotiation.
- `return "", nil` stops the render function before it emits a default body.
- `c.JSON`, `c.XML`, and `c.Write` set the response and status directly.

In contrast, a `<server type="json">` block is **content-negotiated**: it runs
only when the request `Accept` header matches `application/json`. A webhook that
sends no `Accept` header or `Accept: */*` will not match, so do not rely on a
`type="json"` block for a machine endpoint unless the sender sets that header.

If you keep CSRF enabled for a profile and want a structured error, prefer the
`app.SetErrorHandler(403, ...)` route shown above over a typed `<server>`
block.

## See Also

- [Middleware](https://github.com/dreego-stack/dreego/blob/main/_docs/middleware.md) — fixed order and CSRF
- [Forms](https://github.com/dreego-stack/dreego/blob/main/_docs/forms.md) — CSRF token handling for browsers
- [Routing](https://github.com/dreego-stack/dreego/blob/main/_docs/routing.md) — method sections and typed responses
- [Context Methods](context-methods.md) — response helpers on the route context
