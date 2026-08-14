---
type: Spec
title: Form Actions (g-action)
description: Declarative server-side form handling with auto-CSRF, struct mapping, validation, and error binding.
phase: v0.0.16
block: form-actions.1
depends: csrf.1, session.1, context-refactoring.1, routing.1
testcount: 30
timestamp: 2026-07-29T01:16:00+02:00
---

# Form Actions — Spec

## Goal

`g-action="Login"` on `<form>` — Dreego generates CSRF check, `r.ParseForm()`, struct mapping (`form:"email"`), validation (`validate:"required,email"`), error binding (`c.Errors()`, `c.Old()`), and handler dispatch (`func Login(c, form) error`). Developer writes ONE struct + ONE function — the rest is generated boilerplate elimination.

## Non-Goals

- Client-side validation (that's Alpine.js)
- AJAX form submission (that's HTMX)
- File uploads (`multipart/form-data` → later, `g-upload` V2)
- SSG target (skipped, plain `<form action>` passthrough)

## Syntax

```html
<!-- dreego/routes/login/post.dreego -->
<form g-action="Login">
    <input name="email" type="email" />
    {#if c.Errors("email")}
        <p class="error">{{ c.Errors("email") }}</p>
    {/if}
    <button>Login</button>
</form>

<div>Welcome {{ name }}</div>

<go>
    type LoginForm struct {
        Email string `form:"email" validate:"required,email"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        if form.Email == "admin@dreego.dev" {
            c.SetSessionVal("user", form.Email)
            return c.Redirect("/dashboard", 303)
        }
        c.Set("error", "unknown user")
        return c.Redirect("/login", 303)
    }
</go>
```

### Rules

1. **`g-action="Name"`** — maps to exported Go function `Name`. Case-sensitive.
2. **Handler signature**: `func Name(c dreego.Context, form T) error` where T is a struct with `form:""` and optional `validate:""` tags.
3. **Return semantics**:
   - `return c.Redirect(url, code)` → sends redirect response
   - `return err` → 500, recovery middleware catches, or set `c.Set("error", ...)` + re-render
4. **Forms without `g-action`**: Plain `<form method="post" action="/x">` — no Dreego handling, developer calls `c.FormValue()` manually.
5. **No g-action on GET**: `g-action` is POST/PUT/DELETE only. GET uses `c.Query()`.
6. **File-based method wins**: `post.dreego` sets method=POST for all handlers in that file. `<go method="post">` attribute overrides.

## API Additions

### Context Interface (new methods)

```go
type Context interface {
    // existing
    context.Context
    Param(name string) string
    Data(key string) any

    // new for form-actions
    Errors(field string) string   // validation error for field
    Old(field string) string      // previous form value after failed validation
    Redirect(url string, code int) error
}
```

### SSRContext (new methods)

```go
func (c *SSRContext) Errors(field string) string { ... }
func (c *SSRContext) Old(field string) string { ... }
func (c *SSRContext) Redirect(url string, code int) error { ... }
```

### Form validation (new file: `dreego-core/validate.go`)

```go
func BindForm(r *http.Request, target any) error { ... }    // r.ParseForm() + struct mapping
func ValidateForm(form any) map[string]string { ... }       // uses go-playground/validator
```

## AST Changes

```go
type GoSection struct {
    Code        string
    Method      string
    ContentType string
    Action      string   // NEW: "Login" from g-action="Login"
}

type FormAction struct {
    Name       string   // "Login"
    Struct     string   // "LoginForm"
    Handler    string   // "func Login(c, form) error"
    Fields     []FormField
}

type FormField struct {
    Name     string   // "email"
    Validate string   // "required,email" (empty = no validation)
}
```

## Generated Code

### File: `gen/routes.go` (section: POST login handler)

```go
func HandleLoginPost(w http.ResponseWriter, r *http.Request) {
    c := core.NewSSR(w, r)

    // 1. Parse form
    var form LoginForm
    if err := core.BindForm(r, &form); err != nil {
        renderLoginForm(c, form, map[string]string{"_form": err.Error()})
        return
    }

    // 2. Validate
    if errs := core.ValidateForm(form); len(errs) > 0 {
        renderLoginForm(c, form, errs)
        return
    }

    // 3. Call handler
    if err := Login(c, form); err != nil {
        http.Error(w, err.Error(), 500)
    }
}
```

### Generated helpers

```go
func renderLoginForm(c *core.SSRContext, form LoginForm, errs map[string]string) {
    for k, v := range errs {
        c.Set("error_"+k, v)    // → c.Errors("email") 
    }
    c.Set("old_email", form.Email)  // → c.Old("email")
    // ... render same template with errors
}
```

## CSRF Integration

CSRF token already in `<form>` via middleware `csrf_token` hidden input. Form-actions:
- Does NOT generate CSRF token (middleware handles it)
- Does NOT generate hidden `<input>` (developer places `{{ c.CSRFToken() }}` or middleware injects)
- CSRF check: middleware already intercepts POST/PUT/DELETE → 403. Generated code never runs if CSRF fails.
- HTMX auto-appends CSRF header → no manual `<input>` needed for HTMX users

**Decision**: CSRF stays in middleware. Form-actions focuses on struct mapping + validation + dispatch.

## Dependencies

- `go-playground/validator/v10` — new core dependency for struct tag validation
- No other new dependencies. All parsing/routing/rendering via existing standard library code.

## Test Plan

### Parser (12)
| Test | What |
|------|------|
| `g-action-basic` | `<form g-action="Login">` + handler builds |
| `g-action-no-handler` | g-action without matching Go function → error |
| `g-action-unexported` | `g-action="login"` (lowercase) → error |
| `g-action-wrong-arity` | Handler with wrong signature → error |
| `form-no-g-action` | Plain `<form>` without g-action, no handler gen |
| `form-multiple-actions` | Two g-actions on same page |
| `form-struct-tags` | `form:"email"` → ParseForm code generated |
| `form-validate-tags` | `validate:"required"` → validator call generated |
| `form-no-validate` | No validate tag → no validator import |
| `form-plain-form` | Struct without form tags → manual `r.PostForm` |
| `form-handler-signature` | Return type not error → error |
| `form-method-post-file` | `post.dreego` filename → POST method auto-set |

### Codegen (8)
| Test | What |
|------|------|
| `form-gen-csrf` | CSRF check in generated handler |
| `form-gen-parseform` | `r.ParseForm()` + struct mapping generated |
| `form-gen-validate` | `validate.Struct(form)` when validate tags exist |
| `form-gen-saveold` | `c.Set("old_...")` after validation error |
| `form-gen-saveerrors` | `c.Set("error_...")` after validation error |
| `form-gen-redirect` | Handler returning `c.Redirect()` compiles |
| `form-gen-autoimports` | Auto-import `validator` when needed |

### Runtime HTTP (10)
| Test | What |
|------|------|
| `form-submit-valid` | POST valid data → redirect |
| `form-submit-invalid` | POST invalid data → re-render with errors |
| `form-submit-csrf-pass` | POST with CSRF token → OK |
| `form-submit-csrf-fail` | POST without CSRF → 403 |
| `form-submit-csrf-wrong` | POST with wrong CSRF → 403 |
| `form-submit-old-values` | Validation error → old values in HTML |
| `form-submit-errors` | Validation error → error messages in HTML |
| `form-plain-post` | Plain POST without g-action, manual handling |
| `form-multiple-submit` | Two forms, both actions work |
| `form-session-login` | Login handler → session set → redirect |

### Implementation Order

1. **Step 1**: Parser: `g-action` discovery in template, `GoSection.Action`, action→handler matching
2. **Step 2**: Context: `Errors()`, `Old()`, `Redirect()` → real implementation
3. **Step 3**: Validation: `BindForm()`, `ValidateForm()` with `go-playground/validator`
4. **Step 4**: Codegen: Form handler generation (parse → validate → dispatch → render)
5. **Step 5**: Tests: 30 integration tests
6. **Step 6**: Docs: `_docs/forms.md`, CHANGELOG, KB update
