---
type: Concept
title: "Form Actions — Concept"
description: "Declarative form handling with automatic parsing, validation, and CSRF check"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Form Actions — Concept

## Overview

Form Actions replace manual `r.ParseForm()` + `r.FormValue()` with declarative form handlers. The developer defines a struct and a Go function — Dreego generates parsing, validation, and CSRF check.

## Syntax

```html
<!-- routes/login/post.dreego -->
<form g-action="Login">
    <input name="email" type="email" />
    <input name="password" type="password" />
    <button>Login</button>

    {#if errors.general}
        <div class="error">{{ errors.general }}</div>
    {/if}
</form>

<go>
    type LoginForm struct {
        Email    string `form:"email" validate:"required,email"`
        Password string `form:"password" validate:"min=8"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        user, err := db.Authenticate(form.Email, form.Password)
        if err != nil {
            c.Flash("error", "Invalid credentials")
            return c.Redirect("/login")
        }
        c.Session.Set("user_id", user.ID)
        return c.Redirect("/dashboard")
    }
</go>
```

## Generated Code (what Dreego makes out of this)

```go
func handleLoginForm(w http.ResponseWriter, r *http.Request) {
    ctx := dreego.NewSSRContext(r, w)

    // 1. CSRF check (automatic)
    if err := ctx.VerifyCSRF(r); err != nil {
        ctx.RenderError(403, "CSRF validation failed")
        return
    }

    // 2. Form parsing
    r.ParseForm()
    var form LoginForm
    form.Email = r.FormValue("email")
    form.Password = r.FormValue("password")

    // 3. Validation (go-playground/validator)
    if err := validate.Struct(form); err != nil {
        ctx.SaveOld(form)           // Fields for old() restore
        ctx.SaveErrors(err)          // Field-level errors for template
        ctx.Redirect(r.Referer())    // Back to form
        return
    }

    // 4. Call handler
    if err := Login(ctx, form); err != nil {
        ctx.HandleActionError(err)   // Flash, Redirect, etc.
        return
    }
}
```

## Progressive — works with and without JS

```html
<!-- Without JS: normal HTML form, POST, full page reload -->
<form g-action="Login" method="POST">
</form>

<!-- With HTMX: fragment swap, no reload -->
<form g-action="Login" method="POST" hx-post="/api/action/Login" hx-target="#result">
</form>

<!-- With Alpine.js: client-side validation + optimistic UI -->
<form g-action="Login" x-data="formValidation()" @submit.prevent="submit">
</form>
```

All three modes work with the SAME form tag. HTMX and Alpine only upgrade the experience.

## Template Helpers

```html
{#if errors.email}
    <p class="error">{{ errors.email }}</p>
{/if}

<input name="email" value="{{ old.email }}" />

{#if flash.success}
    <div class="success">{{ flash.success }}</div>
{/if}
```

- `errors` — Field-level validation errors (auto-set after validation)
- `old` — Previous inputs (after failed validation)
- `flash` — Session flash messages (after redirect)

## Security

### XSS (Output Encoding)
All template expressions `{{ expression }}` are HTML-escaped:
- `{{ user.Name }}` renders `<script>` as escaped text rather than markup
- Prevents Stored XSS: malicious code in the DB is neutralized on display
- Primeagen's approach: escape on DISPLAY, don't filter on storage

Only `{{ expression|raw }}` allows unescaped HTML — explicit, rare, deliberately risky.

### Auto-Escaping in Core
Not a plugin — MUST be built into the template renderer. Otherwise every Dreego app is insecure.

## Open Design Questions (for GLM)

1. `g-action` vs `g-submit` — name of the attribute?
2. Multiple actions per page — `g-action="Login"` vs `g-action="Register"` in the same template?
3. File uploads — separate `g-upload` or part of `g-action`?
4. Non-SSR Actions — do actions work in SSG mode? (no — build time, no form handling)
5. Action Naming — `Login` is the Go function name — case-sensitive?
6. Redirect vs Stay — `return c.Redirect()` vs `return nil` (stay on page)?
