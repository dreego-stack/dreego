# Dreego Gap-Analyse: Was fehlt, was kann besser werden?

**Datum:** 23.07.2026
**Zweck:** Systematische Analyse, was Dreego im Vergleich zu SvelteKit, Next.js, Phoenix noch fehlt oder wo das Design verbessert werden kann.

## Svelte Runes — haben wir das abgedeckt?

Siehe [[signals-and-runes]] für die vollständige Analyse.

**Kurzfassung:** Svelte Runes sind ein Compiler-Feature für client-seitige Reaktivität. Dreego erreicht dasselbe Ziel (reaktive UI-Updates) durch SSR + HTMX partials + Alpine.js. Wir bilden Runes nicht 1:1 ab, aber das Ergebnis (State-Änderung → DOM-Update) ist dasselbe — nur mit anderer Architektur.

| Svelte Rune    | Dreego-Entsprechung                        |
|----------------|-------------------------------------------|
| `$state(x)`    | `<go>`-Variablen + HTMX Fragment-Updates  |
| `$derived`     | `{#let}` im Template / Berechnung im `<go>` |
| `$effect`      | Alpine `@click`, HTMX `hx-trigger`        |

## Was fehlt Dreego aktuell im Design?

### 1. Form Actions (Server-seitiges Form-Handling)

**Status:** Konzept existiert (`g-submit`), aber nicht detailliert genug.

**Vorbild SvelteKit:**
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

**Für Dreego:**
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

**Was zu klären ist:**
- Wie mapped `g-submit="login"` auf die Go-Funktion?
- Validierung: `go-playground/validator` Struct-Tags automatisch auswerten
- Fehler zurück ins Template geben
- Progressive Enhancement: `<form>` funktioniert auch ohne JS

→ [ ] **ToDo:** Detailliertes Form-Actions-Konzept schreiben

### 2. Progressive Enhancement

**Status:** Erwähnt, aber nicht konkret.

HTMX + Alpine sind per Definition Progressive Enhancement:
- `<form>` funktioniert ohne JS (normales HTML)
- HTMX upgraded das Erlebnis (kein Full-Page-Reload)
- Alpine.js upgraded lokale Interaktionen (Dropdowns)

**Was zu klären ist:**
- Was passiert, wenn JS deaktiviert ist? Funktioniert die App?
- CSRF-Token auch ohne JS?
- Form-Validierung ohne JS?

→ [ ] **ToDo:** Progressive-Enhancement-Strategie dokumentieren

### 3. Error Handling & Validation (detallierter)

**Status:** Grundlegend entschieden (kein `<catch>`, `{#if hasError}`).

**Was fehlt:**
- Form-Validierungs-Feedback (Feld-Level-Errors)
- Flash-Messages (Erfolg/Fehler nach Redirect)
- Error Boundaries auf Komponenten-Ebene
- Stack Traces im Dev-Modus, generische Fehlerseite in Prod

→ [ ] **ToDo:** Error-Handling-Konzept detaillieren

### 4. Middleware-System

**Status:** Im Addon-Konzept erwähnt, aber kein Core-Konzept.

Jedes Framework braucht Middleware-Pipeline. Chi hat bereits Middleware, aber Dreego sollte eigene Convenience-Wrapper bieten:

```go
app := dreego.New()
app.Use(regeo.Logger())
app.Use(regeo.Session(sessionStore))
app.Use(regeo.CSRF())
app.Use(regeo.Auth())
```

**Was zu klären ist:**
- Standard-Middleware, die immer dabei ist (Logger, Recovery, CSRF)
- Reihenfolge der Middleware
- Wie Addons eigene Middleware injizieren

→ [ ] **ToDo:** Middleware-Konzept + Built-in Middleware Liste

### 5. Session & Auth (Core, nicht Addon)

**Status:** Nur als Addon-Idee gelistet (`regeo-auth`).

Frage: Sollte Session-Management Teil des Core sein oder ein Addon?

**Argument für Core:**
- Sessions sind fundamental (fast jede App braucht sie)
- CSRF-Schutz braucht Sessions
- `c.User()` im `<go>`-Block setzt Auth voraus

**Argument für Addon:**
- Nicht jede App braucht Auth (Landing Pages, Blogs)
- Core sollte minimal bleiben
- Addons können via Plugin-Interface alles injecten

→ [ ] **Entscheidung:** Session-Management: Core oder Addon?

### 6. Asset-Pipeline (CSS/JS Bundling)

**Status:** Tailwind CLI erwähnt, aber kein klares Konzept.

**Fragen:**
- Wie wird Tailwind im Dev-Server ausgeführt?
- Minification für Production?
- Wie werden mehrere `<style>`-Sektionen zu einer CSS-Datei zusammengeführt?
- Wie werden `<script>`-Blöcke gebundled?
- Cache-Busting (File-Hashes für Assets)

→ [ ] **ToDo:** Asset-Pipeline-Konzept

### 7. Meta-Framework Features (SEO, Sitemap)

**Status:** Teilweise als Addon-Ideen (`regeo-seo`).

**Was oft im Core erwartet wird:**
- `<title>` pro Seite setzen
- Meta-Description, OpenGraph-Tags
- Automatische Sitemap-Generierung (wenn SSG in V2)
- robots.txt

→ [ ] **ToDo:** SEO-Konzept (Core vs Addon)

### 8. CLI & Dev-Experience im Detail

**Status:** Erwähnt, aber nicht detailliert.

**Was die CLI können sollte:**
```
dreego new my-app           # Projekt scaffolden
dreego dev                   # Dev-Server mit Hot Reload
dreego generate              # Transpiler ausführen
dreego build                 # Production Binary
dreego routes                # Alle Routen anzeigen (Debug)
dreego add <addon>           # Addon installieren
```

**Dev-Server Features:**
- File-Watcher (`.dreego` → `dreego generate` → Browser Reload)
- Error Overlay (Compiler-Fehler direkt im Browser)
- Tailwind JIT Watch
- Proxy zu Backend-API (falls vorhanden)

→ [ ] **ToDo:** CLI-Spezifikation

### 9. Testing-Strategie

**Status:** Kein Konzept vorhanden.

**Fragen:**
- Wie testet man eine `.dreego`-Seite?
- Request-Simulation für `<go>`-Blöcke?
- HTML-Output validieren?
- End-to-End mit Playwright?

```go
func TestIndexPage(t *testing.T) {
    rec := dreegotest.Get("/", nil)
    assert.Contains(t, rec.Body.String(), "<h1>Willkommen</h1>")
}
```

→ [ ] **ToDo:** Testing-Konzept

### 10. Performance & Caching

**Status:** Nicht thematisiert.

- Template-Caching (gerenderte Templates wiederverwenden)
- ETag / Last-Modified Header
- Static-Asset-Caching
- Datenbank-Query-Caching

→ [ ] **ToDo:** Performance-Konzept (V2)

### 11. Internationalisierung (i18n)

**Status:** Als Addon-Idee gelistet.

Frage: Routing für Sprachen? (`/de/about`, `/en/about`)

→ [ ] **ToDo:** i18n-Konzept (V2, aber Routing jetzt bedenken)

### 12. Security (über CSRF/XSS hinaus)

**Status:** CSRF + XSS erwähnt.

**Was noch fehlt:**
- Content Security Policy (CSP) Header
- Rate Limiting
- SQL-Injection (weniger relevant bei Go + ORM, aber trotzdem)
- Secure Cookie Flags (HttpOnly, SameSite)
- CORS-Konfiguration

→ [ ] **ToDo:** Security-Konzept detaillieren

## Architektur-Entscheidungen, die wir JETZT treffen sollten

Um uns nicht die Zukunft zu verbauen:

| Entscheidung                    | Warum jetzt wichtig                              |
|---------------------------------|--------------------------------------------------|
| Plugin-Interface finalisieren   | Änderungen später = Breaking Changes für Addons  |
| Transpiler-Pipeline definieren  | Erweiterungspunkte für TS, SSG, etc.             |
| Context-Objekt designen         | `c.User()`, `c.Session()`, `c.Param()` — API     |
| Error-Handling-Strategie        | Wie fließen Fehler von Addons ins Template?      |
| Routing-Konventionen            | `/routes/` vs `/pages/`, `layout.dreego`-Konzept   |
| Asset-Pipeline-Architektur      | Wie kommen CSS/JS ins Binary?                    |

## Prioritäten für die nächste Planungs-Runde

1. **Form Actions detailliert ausarbeiten** (höchste Priorität, Kern-Feature)
2. **Plugin-Interface finalisieren** (damit Addons gebaut werden können)
3. **Middleware-Konzept** (Session/CSRF/Auth Basis)
4. **CLI-Spezifikation** (Developer Experience von Tag 1)
5. **Testing-Strategie** (weil wir testgetrieben entwickeln wollen)
