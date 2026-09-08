
---
type: Decision
title: Form Actions — g-action + Generated Handlers
description: Declarative form handlers via g-action with generated CSRF+validation pipeline
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Form Actions — g-action + Generated Handlers

**Date:** 2026-07-28
**Status:** Accepted — generated pipeline is current; see note below

> **Historical note:** The generated `g-action` pipeline is current. The
> dependency on `go-playground/validator` is superseded — validation uses
> built-in validators in Core (see [_docs/forms.md](../../_docs/forms.md)). The
> `dreego.Context` interface and target-agnostic handler signature remain
> current (see [context-design](context-design.md)).

## Decision

**Declarative form handlers via `g-action`.** The developer defines a struct + function — Dreego generates CSRF check, form parsing, validation, and handler call.

## Syntax

```html
<form g-action="Login">
    <input name="email" type="email" />
    <button>Login</button>
</form>

<server>
    type LoginForm struct {
        Email string `form:"email" validate:"required,email"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        user, _ := db.Authenticate(form.Email)
        c.Session.Set("user_id", user.ID)
        return c.Redirect("/dashboard")
    }
</server>
```

## GLM Corrections to the Original Concept

### 1. Handler Signature: `dreego.Context` instead of `*SSRContext`

Wrong: `func Login(c *SSRContext, form LoginForm) error`
Correct: `func Login(c dreego.Context, form LoginForm) error`

Rationale: Rendering and tests can run without exposing a raw HTTP request in
application code.

### 2. Access `errors`, `old`, `flash` via `c`

Not as magic template variables, but via context:

```html
{#if c.Errors("email")}
    <p class="error">{{ c.Errors("email") }}</p>
{/if}
<input name="email" value="{{ c.Old("email") }}" />
```

Rationale: Works in tests without an HTTP server (AGENTS.md guarantee #7).

### 3. Non-HTTP hosts do not inherit form-action semantics

Wails does not silently receive HTTP form-action semantics. A plain
`<form action="/api/...">` remains an explicit HTTP escape hatch when the
application provides that endpoint.

### 4. Forms without `g-action` must work

```html
<form method="post" action="/api/custom">
```
No Dreego handling. The developer is responsible themselves.

## Generated Pipeline

```
CSRF Check → r.ParseForm() → Struct Mapping → validate.Struct(form)
  → c.SaveOld(form) + c.SaveErrors(err) → Handler Call
  → c.Redirect() | error
```

## Naming

- Attribute: `g-action` (not `g-submit`)
- Action name = Go function name, case-sensitive, exported
- `Login` → `<form g-action="Login">`

## Return Semantics

Handler MUST do one of:
- `return c.Redirect(url)` — PRG pattern
- `return err` — Error page or flash+redirect

`return nil` without redirect is not allowed.

## Progressive Enhancement

| Mode         | Without JS        | HTMX                     | Alpine.js              |
|--------------|-------------------|--------------------------|------------------------|
| Form Submit  | Full page reload  | `hx-post` fragment swap  | `@submit` + fetch      |
| Validation   | Server-side       | Server-side              | Client + Server        |
| CSRF         | Automatic         | Automatic                | Automatic              |

## XSS Protection

All `{{ variable }}` expressions use context-aware output escaping.
Only `{{ variable|raw }}` bypasses escaping — explicit, rare, consciously risky.

## Consequences

- Validation uses built-in validators in Core; no external validator dependency
- CSRF check via middleware, not hardcoded in the action handler
- File uploads: Part of `g-action` via `multipart.Form` (V1)
- `g-upload` for chunked/streaming only in V2
