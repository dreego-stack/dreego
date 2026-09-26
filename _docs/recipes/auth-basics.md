# Recipe: Authentication Basics

Dreego Core has no built-in user model. Authentication is an application or
plugin concern. What Core provides is the session store, CSRF protection, and
the context methods you build login and logout on. This recipe shows that
foundation: session-store integration, a login flow, `DestroySession` logout,
and the flash pattern.

> **Boundary:** Core ships the session `Store` and the context methods, not a
> user database or an auth provider. A full auth system — password hashing,
> OAuth, passkeys — lives in application code or a plugin repository.

## 1. Configure the session store

Create a cookie store and register it before the App is built. The signing
secret must be at least 32 bytes; `NewCookieStore` panics on a shorter one and
`App.Build` rejects an unsafe store with an error.

```go
app := dreego.New()

store := dreego.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
if err := app.SetSessionStore(store); err != nil {
	log.Fatal(err)
}
```

Registering a store enables:
- session reads and writes (`c.SessionVal`, `c.SetSessionVal`),
- the session cookie with secure defaults (`HttpOnly`, `SameSite=Lax`,
  TLS-aware `Secure`, `Path=/`),
- CSRF, which stores its token in the session and checks state-changing
  requests.

For AES-256-GCM encrypted cookies or trusted TLS proxies, see
[Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md).

If you want session and CSRF only for part of the site, scope the store and CSRF
with `PROFILE` instead of `SetSessionStore`. See
[Machine Endpoints](machine-endpoints.md) for that pattern.

## 2. Login

Never store a raw password in the session. Store an opaque identifier and look
up the user on each request.

```dreego
<!-- www/routes/login/+page.dreego -->
<server>
	type LoginForm struct {
		Email    string `form:"email" validate:"required,email"`
		Password string `form:"password" validate:"required"`
	}

	func Login(c dreego.Context, form LoginForm) error {
		user, err := authenticate(form.Email, form.Password)
		if err != nil {
			c.Flash("error", "Invalid email or password")
			return c.Redirect("/login", 303)
		}
		c.SetSessionVal("user_id", user.ID)
		c.Flash("notice", "Welcome back")
		return c.Redirect("/dashboard", 303)
	}

	loginError := c.FlashGet("error")
</server>

<body>
	<h1>Sign in</h1>
	{#if loginError}<p role="alert">{{ loginError }}</p>{/if}
	<form g-action="Login" method="post">
		{{ c.CSRFInput()|raw }}
		<label for="email">Email</label>
		<input id="email" name="email" type="email" value="{{ c.Old("email") }}">
		{#if c.Errors("email")}<p class="error">{{ c.Errors("email") }}</p>{/if}
		<label for="password">Password</label>
		<input id="password" name="password" type="password">
		{#if c.Errors("password")}<p class="error">{{ c.Errors("password") }}</p>{/if}
		<button type="submit">Sign in</button>
	</form>
</body>
```

Notes:

- `c.CSRFInput()` renders the complete hidden `csrf_token` field. Use `{{ ...
  |raw }}` so the markup is not escaped; the token value is escaped inside the
  helper.
- The `g-action` pipeline re-renders on validation or bind failure, but a
  handler that returns `nil` falls through to a default `303` redirect back to
  the route. A wrong password therefore cannot re-render inline from the
  handler: store a flash message and redirect (`Post-Redirect-Get`), then read it
  with `c.FlashGet` in the GET section.
- `authenticate` must compare hashed passwords with a constant-time comparison.
  It is application code (or a plugin call), declared in a glue file in the same
  route package or imported through `GOIMPORT`. See
  [Calling App Code](app-code.md).
- On success, store only an identifier and post a one-shot notice.

## 3. Read the session in a protected route

```html
<server>
	userID := c.SessionVal("user_id")
	user, err := loadUser(userID)
	if err != nil {
		return "", err
	}
</server>

<body>
	{#if user == nil}
		<p>Please <a href="/login">sign in</a> to continue.</p>
	{#else}
		<h1>Hello {{ user.Name }}</h1>
	{/if}
</body>
```

A plain render section returns `(string, error)` and renders HTML; it cannot abort
with a redirect. Use this shape when a guarded page can render a signed-out
state. When a guard must redirect before rendering, handle it in a `g-action`
handler (whose pipeline understands `c.Redirect` returning
`dreego.ErrRedirect`) or in an explicitly registered route handler — not in a
plain `<server>` section.

Fail closed: a missing or unknown `user_id` means "not signed in". Never treat
an empty session as a valid user.

## 4. Logout with `DestroySession`

Logout must clear the whole session, not one key. `DestroySession` clears every
key and sends one expired cookie that preserves `Secure`, `HttpOnly`,
`SameSite`, and `Path` from the active policy. Use a POST form so logout is
CSRF-protected and cannot be triggered by a link:

```dreego
<!-- www/routes/logout/+page.dreego -->
<server>
	type LogoutForm struct{}

	func Logout(c dreego.Context, form LogoutForm) error {
		c.DestroySession()
		if err := c.SessionError(); err != nil {
			return err
		}
		return c.Redirect("/", 303)
	}
</server>

<body>
	<form g-action="Logout" method="post">
		{{ c.CSRFInput()|raw }}
		<button type="submit">Sign out</button>
	</form>
</body>
```

`c.DelSessionVal("user_id")` deletes a single key and is **not** a logout: it
leaves the rest of the session (including CSRF state) in place. Use
`DestroySession` to invalidate the login.

A plain `<server>` render section does not interpret `c.Redirect`'s
`ErrRedirect`; only the `g-action` pipeline does. Put logout in a `g-action`
handler (as above) or an explicitly registered handler, not in a plain
`<server>` block with only a `return c.Redirect(...)`.

The built-in cookie store is stateless: it does not track a server-side session
identifier, so a stolen cookie can be replayed until it expires or the signing
secret rotates. Immediate invalidation requires a server-side store plugin or a
secret rotation. See
[Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md).

## 5. Flash messages

A flash message survives exactly one redirect: it is written to the session,
read once on the next request, and consumed. Use it after a redirect when you
cannot render the target page inline. A `g-action` handler writes the flash
before returning `c.Redirect` (see the login handler in section 2); the pipeline
understands the redirect and completes the request.

```dreego
<server method="post">
	type SaveForm struct {
		Name string `form:"name" validate:"required"`
	}

	func Save(c dreego.Context, form SaveForm) error {
		if err := saveProfile(c.SessionVal("user_id"), form.Name); err != nil {
			return err
		}
		c.Flash("notice", "Profile saved")
		return c.Redirect("/profile", 303)
	}
</server>
```

```dreego
<!-- reads once on the next GET; the value is then gone -->
<server>
	notice := c.FlashGet("notice")
</server>

<body>
	{#if notice}<p role="status">{{ notice }}</p>{/if}
</body>
```

In a plain (non-`g-action`) POST section, write the flash, issue the redirect
directly, and return an empty body instead of returning the redirect error:

```dreego
<server method="post">
	c.Flash("notice", "Profile saved")
	c.Redirect("/profile", 303)
	return "", nil
</server>
```

Method set:

| Method | Behavior |
|--------|----------|
| `c.Flash(key, message)` | Store a one-shot message under an internal `flash_` key |
| `c.FlashGet(key)` | Read the message and consume it in one call |
| `c.FlashPeek(key)` | Read the message without consuming it |

Flash values live in the session store, so a flash without a configured store is
a no-op and `c.FlashGet` returns `""`. One value per key; a second `Flash` for
the same key overwrites the first.

## Rules

1. Configure the store before the App is built; never rotate it per request.
2. Store an opaque identifier, never credentials or secrets, in the session.
3. Authenticate with hashed passwords and constant-time comparison.
4. Logout calls `DestroySession`, not `DelSessionVal`.
5. Use `Flash`/`FlashGet` for post-redirect messages, not for persistent state.
6. Fail closed: a missing session value means "not logged in".

## See Also

- [Session Encryption](https://github.com/dreego-stack/dreego/blob/main/_docs/session-encryption.md) — encryption, policy, rotation
- [Forms](https://github.com/dreego-stack/dreego/blob/main/_docs/forms.md) — `g-action`, validation, CSRF
- [Context Methods](context-methods.md) — session and flash methods
- [Machine Endpoints](machine-endpoints.md) — `PROFILE` session/CSRF scoping
- [Plugins](https://github.com/dreego-stack/dreego/blob/main/_docs/plugins.md) — auth as a plugin
