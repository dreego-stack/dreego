
---
type: Decision
title: Form Actions — g-action + generierte Handler
description: Deklarative Form-Handler via g-action mit generierter CSRF+Validierungs-Pipeline
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Form Actions — g-action + generierte Handler

## Entscheidung

**Deklarative Form-Handler via `g-action`.** Der Entwickler definiert Struct+Funktion — Dreego generiert CSRF-Check, Form-Parsing, Validierung und Handler-Aufruf.

## Syntax

```html
<form g-action="Login">
    <input name="email" type="email" />
    <button>Login</button>
</form>

<go>
    type LoginForm struct {
        Email string `form:"email" validate:"required,email"`
    }

    func Login(c dreego.Context, form LoginForm) error {
        user, _ := db.Authenticate(form.Email)
        c.Session.Set("user_id", user.ID)
        return c.Redirect("/dashboard")
    }
</go>
```

## GLM-Korrekturen am Original-Konzept

### 1. Handler-Signatur: `dreego.Context` statt `*SSRContext`

Falsch: `func Login(c *SSRContext, form LoginForm) error`
Richtig: `func Login(c dreego.Context, form LoginForm) error`

Begründung: Target-Agnostik (AGENTS.md Garantie #2). SSG-Target übergibt `SSGContext`.

### 2. `errors`, `old`, `flash` über `c` beziehen

Nicht als Magic-Template-Variablen, sondern via Context:

```html
{#if c.Errors("email")}
    <p class="error">{c.Errors("email")}</p>
{/if}
<input name="email" value="{c.Old("email")}" />
```

Begründung: Funktioniert in Tests ohne HTTP-Server (AGENTS.md Garantie #7).

### 3. SSG: Actions werden nicht generiert

SSG-Target überspringt `g-action`-Codegen. Plain `<form action="/api/...">` bleibt als Escape-Hatch erhalten. Transpiler kapselt Action-Codegen hinter `TargetSSR`.

### 4. Form ohne `g-action` muss funktionieren

```html
<form method="post" action="/api/custom">
```
Kein Dreego-Handling. Der Entwickler ist selbst verantwortlich.

## Generierte Pipeline

```
CSRF-Check → r.ParseForm() → Struct-Mapping → validate.Struct(form)
  → c.SaveOld(form) + c.SaveErrors(err) → Handler-Aufruf
  → c.Redirect() | c.Render() | error
```

## Benennung

- Attribut: `g-action` (nicht `g-submit`)
- Action-Name = Go-Funktionsname, case-sensitive, exported
- `Login` → `<form g-action="Login">`

## Return-Semantik

Handler MUSS eines tun:
- `return c.Redirect(url)` — PRG-Pattern
- `return c.Render(template, data)` — Template rendern
- `return err` — Fehlerseite oder Flash+Redirect

`return nil` ohne Redirect/Render ist nicht erlaubt.

## Progressiv

| Modus         | Ohne JS           | HTMX                     | Alpine.js              |
|---------------|-------------------|--------------------------|------------------------|
| Form-Submit   | Full-Page-Reload  | `hx-post` Fragment-Tausch| `@submit` + fetch      |
| Validierung    | Server-seitig     | Server-seitig            | Client + Server        |
| CSRF          | Automatisch       | Automatisch              | Automatisch            |

## XSS-Schutz

Alle `{variable}` im Template werden automatisch HTML-escaped (Output Encoding).
Nur `{variable|raw}` umgeht das Escaping — explizit, selten, bewusst riskant.

## Konsequenzen

- `go-playground/validator` ist Core-Dependency
- CSRF-Check via Middleware, nicht hart im Action-Handler
- File-Uploads: Teil von `g-action` via `multipart.Form` (V1)
- `g-upload` für chunked/streaming erst in V2
