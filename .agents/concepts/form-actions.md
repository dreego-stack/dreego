
---
type: Concept
title: "Form Actions — Konzept"
description: "Deklaratives Form-Handling mit automatischem Parsing, Validierung und CSRF-Check"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Form Actions — Konzept

## Übersicht

Form Actions ersetzen manuelles `r.ParseForm()` + `r.FormValue()` durch deklarative Form-Handler. Der Entwickler definiert einen Struct und eine Go-Funktion — Dreego generiert Parsing, Validierung und CSRF-Check.

## Syntax

```html
<!-- routes/login.dreego -->
<form g-action="Login">
    <input name="email" type="email" />
    <input name="password" type="password" />
    <button>Login</button>

    {#if errors.general}
        <div class="error">{errors.general}</div>
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
            c.Flash("error", "Ungültige Anmeldedaten")
            return c.Redirect("/login")
        }
        c.Session.Set("user_id", user.ID)
        return c.Redirect("/dashboard")
    }
</go>
```

## Generierter Code (was Dreego daraus macht)

```go
func handleLoginForm(w http.ResponseWriter, r *http.Request) {
    ctx := dreego.NewSSRContext(r, w)

    // 1. CSRF-Check (automatisch)
    if err := ctx.VerifyCSRF(r); err != nil {
        ctx.RenderError(403, "CSRF validation failed")
        return
    }

    // 2. Form-Parsing
    r.ParseForm()
    var form LoginForm
    form.Email = r.FormValue("email")
    form.Password = r.FormValue("password")

    // 3. Validierung (go-playground/validator)
    if err := validate.Struct(form); err != nil {
        ctx.SaveOld(form)           // Felder für old()-Wiederherstellung
        ctx.SaveErrors(err)          // Feld-Level-Errors für Template
        ctx.Redirect(r.Referer())    // Zurück zum Form
        return
    }

    // 4. Handler aufrufen
    if err := Login(ctx, form); err != nil {
        ctx.HandleActionError(err)   // Flash, Redirect, etc.
        return
    }
}
```

## Progressiv — funktioniert mit und ohne JS

```html
<!-- Ohne JS: normales HTML-Form, POST, Full-Page-Reload -->
<form g-action="Login" method="POST">
</form>

<!-- Mit HTMX: Fragment-Tausch, kein Reload -->
<form g-action="Login" method="POST" hx-post="/api/action/Login" hx-target="#result">
</form>

<!-- Mit Alpine.js: Client-seitige Validierung + optimistisches UI -->
<form g-action="Login" x-data="formValidation()" @submit.prevent="submit">
</form>
```

Alle drei Modi funktionieren mit DEMSELBEN Form-Tag. HTMX und Alpine upgraden nur das Erlebnis.

## Template-Helfer

```html
{#if errors.email}
    <p class="error">{errors.email}</p>
{/if}

<input name="email" value="{old.email}" />

{#if flash.success}
    <div class="success">{flash.success}</div>
{/if}
```

- `errors` — Feld-Level-Validierungsfehler (auto-gesetzt nach Validierung)
- `old` — Vorherige Eingaben (nach fehlgeschlagener Validierung)
- `flash` — Session-Flash-Messages (nach Redirect)

## Sicherheit

### XSS (Output Encoding)
Alle Template-Variablen `{variable}` werden HTML-escaped:
- `{user.Name}` → `&lt;script&gt;` wird zu `&amp;lt;script&amp;gt;`
- Verhindert Stored XSS: Schadcode in der DB wird beim Anzeigen neutralisiert
- Primeagens Ansatz: beim DISPLAY escapen, nicht beim Speichern filtern

Nur `{variable|raw}` erlaubt unescaped HTML — explizit, selten, bewusst riskant.

### Auto-Escaping im Core
Kein Addon — MUSS im Template-Renderer eingebaut sein. Sonst ist jede Dreego-App unsicher.

## Offene Design-Fragen (für GLM)

1. `g-action` vs `g-submit` — Name des Attributs?
2. Mehrere Actions pro Seite — `g-action="Login"` vs `g-action="Register"` im selben Template?
3. File-Uploads — eigenes `g-upload` oder Teil von `g-action`?
4. Non-SSR Actions — funktionieren Actions im SSG-Modus? (nein — Build-Zeit, kein Form-Handling)
5. Action-Naming — `Login` ist der Go-Funktionsname — Case-sensitiv?
6. Redirect vs Render — `return c.Redirect()` vs `return c.Render()` vs `return nil` (bleibt auf Seite)?
