# Form Actions (g-action)

Declarative server-side form handling. One struct, one function — Dreego generates parsing, validation, dispatch, and redirect logic.

## Quick Example

```html
<!-- dreego/routes/login/post.dreego -->
<go>
    type LoginForm struct {
        Email    string `form:"email" validate:"required,email"`
        Password string `form:"password" validate:"required,min=8"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        if form.Email == "admin@dreego.dev" {
            c.SetSessionVal("user", form.Email)
            return c.Redirect("/dashboard", 303)
        }
        c.SetSessionVal("error", "unknown user")
        return c.Redirect("/login", 303)
    }
</go>

<div>
    <h1>Login</h1>
    {#if c.Errors("email")}<p class="error">{{ c.Errors("email") }}</p>{/if}
    <form g-action="Login" method="post">
        <input name="email" type="email" value="{{ c.Old("email") }}">
        <input name="password" type="password">
        <button type="submit">Login</button>
    </form>
</div>
```

The `g-action="Login"` attribute on the `<form>` tells Dreego to generate a POST handler that:
1. Parses the form body via `r.ParseForm()`
2. Maps values to the `LoginForm` struct via `form:"..."` tags
3. Validates via `validate:"..."` tags
4. Calls `Login(c, form)` on success
5. Re-renders the page with errors on validation failure

## Generated Pipeline

```
POST /login
  → r.ParseForm()
  → Struct mapping (form:"email" → LoginForm.Email)
  → Validation (validate:"required,email" → errors map)
  → On validation fail: SaveErrors + SaveOld → re-render GET template
  → On success: call Login(c, form)
    → c.Redirect(url, 303) → ErrRedirect → handler returns
    → Any other error → 500
    → nil return → default redirect to GET
```

## Struct Tags

### `form:"name"`

Maps the form field `name` to the struct field. Falls back to lowercase field name if no tag.

```go
type Form struct {
    Email    string `form:"email"`    // maps form field "email"
    Remember bool                      // maps form field "remember" (lowercase)
}
```

### `validate:"rule,rule"`

Built-in validators (no external dependencies):

| Rule | Description |
|------|-------------|
| `required` | Value must not be empty |
| `email` | Must contain `@` and `.` |
| `min=N` | Minimum length N |
| `max=N` | Maximum length N |

Multiple rules are comma-separated: `validate:"required,email"`.

## Context Methods

### `c.Errors(field) string`

Returns the validation error message for a field. Only available after validation failure.

```html
{#if c.Errors("email")}
    <p class="error">{{ c.Errors("email") }}</p>
{/if}
```

### `c.Old(field) string`

Returns the previously submitted value after validation failure. Useful for re-populating form fields.

```html
<input name="email" value="{{ c.Old("email") }}">
```

### `c.Redirect(url string, code int) error`

Sends an HTTP redirect and returns `ErrRedirect` to signal the handler pipeline. Use for Post-Redirect-Get pattern.

```go
func Login(c dreego.Context, form LoginForm) error {
    // ... authenticate ...
    return c.Redirect("/dashboard", 303)
}
```

## CSRF Protection

CSRF is handled by the middleware, not the form handler. With `app.SetSessionStore(...)` before the app is built:
- CSRF token is set via cookie on any GET request
- POST/PUT/DELETE without valid token → 403
- Token sent via `X-CSRF-Token` header or `csrf_token` form field

HTMX does not send the CSRF token automatically. For HTMX requests, either keep
the hidden field in the form (HTMX serializes form fields, so the token is
submitted as `csrf_token`) or configure the header once via `hx-headers`:

```html
<body hx-headers='{"X-CSRF-Token": "{{ c.CSRFToken() }}"}'>
```

For plain HTML forms, include a hidden field:

```html
<input type="hidden" name="csrf_token" value="{{ c.CSRFToken() }}">
```

Disable CSRF for an API-only app before build with `app.SetCSRF(false)`. Route-specific exemptions require a separate explicit contract and must not weaken unrelated routes.

## Without g-action

Forms without `g-action` are plain HTML forms — no handler generation. Use `c.FormValue()` manually:

```html
<go>
    email := c.FormValue("email")
    c.Set("email", email)
</go>
<div>
<form method="post">
    <input name="email">
    <button>Submit</button>
</form>
</div>
```
