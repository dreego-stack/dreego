
---
type: Concept
title: "Dreego Gap Analysis: What's missing, what can be improved?"
description: "Systematic analysis of missing features compared to SvelteKit, Next.js, Phoenix"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Gap Analysis: What's missing, what can be improved?

**Date:** 2026-07-28
**Purpose:** Systematic analysis of what Dreego still lacks compared to SvelteKit, Next.js, Phoenix, or where the design can be improved.

## Svelte Runes — do we have this covered?

See [Signals & Runes](signals-and-runes.md) for the full analysis.

**Summary:** Svelte Runes are a compiler feature for client-side reactivity. Dreego achieves the same goal (reactive UI updates) through SSR + HTMX partials + Alpine.js. We don't map runes 1:1, but the result (state change → DOM update) is the same — just with a different architecture.

| Svelte Rune    | Dreego Equivalent                          |
|----------------|-------------------------------------------|
| `$state(x)`    | `<go>` variables + HTMX fragment updates  |
| `$derived`     | `{#let}` in template / computation in `<go>` |
| `$effect`      | Alpine `@click`, HTMX `hx-trigger`        |

## What's currently missing in Dreego's design?

### 1. Form Actions (Server-side form handling)

**Status:** Concept exists (`g-submit`), but not detailed enough.

**SvelteKit model:**
```svelte
<form method="POST" action="?/login">
    <input name="email" />
    <button>Login</button>
</form>
```
```ts
export const actions = {
    login: async ({ request }) => {
        const data = await request.formData()
        // ...
    }
}
```

**For Dreego:**
```html
<form g-submit="login">
    <input name="email" />
    <button>Login</button>
</form>
```

```html
<go>
    type LoginForm struct {
        Email    string `form:"email" validate:"required,email"`
        Password string `form:"password" validate:"min=8"`
    }

    func login(form LoginForm) error {
        return db.Authenticate(form.Email, form.Password)
    }
</go>
```

**What needs clarifying:**
- How does `g-submit="login"` map to the Go function?
- Validation: automatically evaluate `go-playground/validator` struct tags
- Return errors to the template
- Progressive Enhancement: `<form>` also works without JS

→ [ ] **ToDo:** Write detailed Form Actions concept

### 2. Progressive Enhancement

**Status:** Mentioned, but not concrete.

HTMX + Alpine are by definition Progressive Enhancement:
- `<form>` works without JS (normal HTML)
- HTMX upgrades the experience (no full page reload)
- Alpine.js upgrades local interactions (dropdowns)

**What needs clarifying:**
- What happens when JS is disabled? Does the app work?
- CSRF token also without JS?
- Form validation without JS?

→ [ ] **ToDo:** Document Progressive Enhancement strategy

### 3. Error Handling & Validation (more detailed)

**Status:** Basic decision made (no `<catch>`, `{#if hasError}`).

**What's missing:**
- Form validation feedback (field-level errors)
- Flash messages (success/error after redirect)
- Error boundaries at component level
- Stack traces in dev mode, generic error page in prod

→ [ ] **ToDo:** Detail Error Handling concept

### 4. Middleware System

**Status:** Mentioned in the addon concept, but no core concept.

Every framework needs a middleware pipeline. Chi already has middleware, but Dreego should provide its own convenience wrappers:

```go
app := dreego.New()
app.Use(regeo.Logger())
app.Use(regeo.Session(sessionStore))
app.Use(regeo.CSRF())
app.Use(regeo.Auth())
```

**What needs clarifying:**
- Standard middleware that's always included (Logger, Recovery, CSRF)
- Order of middleware
- How addons inject their own middleware

→ [ ] **ToDo:** Middleware concept + built-in middleware list

### 5. Session & Auth (Core, not Addon)

**Status:** Only listed as addon idea (`regeo-auth`).

Question: Should session management be part of the core or an addon?

**Argument for Core:**
- Sessions are fundamental (almost every app needs them)
- CSRF protection needs sessions
- `c.User()` in the `<go>` block requires auth

**Argument for Addon:**
- Not every app needs auth (landing pages, blogs)
- Core should remain minimal
- Addons can inject everything via plugin interface

→ [ ] **Decision:** Session Management: Core or Addon?

### 6. Asset Pipeline (CSS/JS Bundling)

**Status:** Tailwind CLI mentioned, but no clear concept.

**Questions:**
- How is Tailwind run in the dev server?
- Minification for production?
- How are multiple `<style>` sections merged into one CSS file?
- How are `<script>` blocks bundled?
- Cache busting (file hashes for assets)

→ [ ] **ToDo:** Asset Pipeline concept

### 7. Meta-Framework Features (SEO, Sitemap)

**Status:** Partially as addon ideas (`regeo-seo`).

**What is often expected in core:**
- Set `<title>` per page
- Meta description, OpenGraph tags
- Automatic sitemap generation (when SSG in V2)
- robots.txt

→ [ ] **ToDo:** SEO concept (Core vs Addon)

### 8. CLI & Dev Experience in Detail

**Status:** Mentioned, but not detailed.

**What the CLI should be able to do:**
```
dreego new my-app           # Scaffold project
dreego dev                   # Dev server with hot reload
dreego generate              # Run transpiler
dreego build                 # Production binary
dreego routes                # Show all routes (debug)
dreego add <addon>           # Install addon
```

**Dev Server Features:**
- File watcher (`.dreego` → `dreego generate` → Browser Reload)
- Error overlay (compiler errors directly in browser)
- Tailwind JIT Watch
- Proxy to backend API (if present)

→ [ ] **ToDo:** CLI specification

### 9. Testing Strategy

**Status:** No concept exists.

**Questions:**
- How do you test a `.dreego` page?
- Request simulation for `<go>` blocks?
- Validate HTML output?
- End-to-end with Playwright?

```go
func TestIndexPage(t *testing.T) {
    rec := dreegotest.Get("/", nil)
    assert.Contains(t, rec.Body.String(), "<h1>Welcome</h1>")
}
```

→ [ ] **ToDo:** Testing concept

### 10. Performance & Caching

**Status:** Not addressed.

- Template caching (reuse rendered templates)
- ETag / Last-Modified headers
- Static asset caching
- Database query caching

→ [ ] **ToDo:** Performance concept (V2)

### 11. Internationalization (i18n)

**Status:** Listed as addon idea.

Question: Routing for languages? (`/de/about`, `/en/about`)

→ [ ] **ToDo:** i18n concept (V2, but consider routing now)

### 12. Security (beyond CSRF/XSS)

**Status:** CSRF + XSS mentioned.

**What's still missing:**
- Content Security Policy (CSP) headers
- Rate Limiting
- SQL Injection (less relevant with Go + ORM, but still)
- Secure Cookie Flags (HttpOnly, SameSite)
- CORS configuration

→ [ ] **ToDo:** Detail Security concept

## Architecture Decisions we should make NOW

So we don't close doors for the future:

| Decision                         | Why important now                                   |
|----------------------------------|-----------------------------------------------------|
| Finalize Plugin Interface        | Changes later = breaking changes for addons         |
| Define Transpiler Pipeline       | Extension points for TS, SSG, etc.                  |
| Design Context Object            | `c.User()`, `c.Session()`, `c.Param()` — API        |
| Error Handling Strategy          | How do errors flow from addons into the template?   |
| Routing Conventions              | `/routes/` vs `/pages/`, `layout.dreego` concept    |
| Asset Pipeline Architecture      | How do CSS/JS get into the binary?                  |

## Priorities for the Next Planning Round

1. **Detail Form Actions** (highest priority, core feature)
2. **Finalize Plugin Interface** (so addons can be built)
3. **Middleware Concept** (Session/CSRF/Auth basics)
4. **CLI Specification** (Developer Experience from Day 1)
5. **Testing Strategy** (because we want to develop test-driven)
