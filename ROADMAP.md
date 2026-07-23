# Dreego Roadmap

**Stand:** 23.07.2026 — Transpiler-Prototyp funktioniert.
**Prinzip:** Jedes Feature startet als `0.x.0` (quick & dirty, aber funktioniert) und wird in `0.x.1...n` sauber integriert. Breaking Changes bis v1.0 erlaubt.

---

## Phase 0 — 0.0.x: Solider erster Prototyp

Ziel: Ein lauffähiges Mini-Projekt (Blog, Landingpage) kann damit gebaut werden.

### ✅ Erledigt (Prototyp in code/transpiler/)

- Formale Pipeline: Lexer → Parser → AST → CodeGen
  - `lexer/` — Tokenizer
  - `parser/` — 5 Sektions-Parser (head, go, div, script, style)
  - `ast/` — AST-Typen
  - `codegen/` — Go-Code-Generator + SSR-Target
- Alle 5 Sektionen: `<head>`, `<go>`, `<div>`, `<script>`, `<style>`
- Template-Logik: `{var}`, `{#if}`, `{#each}`
- File-based Routing: `routes/index.dreego` → `/`, `routes/about.dreego` → `/about`
- `dreego generate` CLI
- Lauffähiger Server mit net/http 1.22+

### Offen für Phase 0

- `layout.dreego` — Wrapper-Layout mit `{#slot}` für Seiteninhalt
  - DOCTYPE, `<html>`, `<head>` aus Layout
  - `<head>`-Inhalte pro Seite mergen (Seiten-Head + Layout-Head)
  - Navigation, Footer im Layout
- Dynamische Routen-Segmente: `routes/users/[id].dreego`
- `dreego init <name>` — Projekt-Scaffolding
  - `routes/`, `assets/`, `dreego.config.json`
- `dreego dev` — Dev-Server via `air` (Hot Reload)
- Recovery-Middleware (Panic → 500)
- 404-Seite (automatisch oder `routes/404.dreego`)
- CSS-Scoping korrekt: Scope-Hash an Klassen anhängen, nicht nur extrahieren
- `<script>`-Blöcke in `<head>` oder vor `</body>` (konfigurierbar)
- Error-Messages mit `.dreego`-Zeilennummern bei Parse-Fehlern
- Snapshot-Tests für generierten Code
- `dreego build` — `go build` Wrapper mit `//go:embed` für Assets
- Statische Assets ausliefern (`/static/` → `assets/` Ordner)

---

## Phase 1 — 0.x: Auf dem Weg zu v1.0

Jeder Punkt = eine `0.x.0` Version. Cleanup in `0.x.1...n`.

### Routing erweitern
- `[...catchall]` Catch-All-Segmente
- `[[optional]]` Optionale Segmente
- Route-Gruppen `(group)/` für Layout-Gruppierung
- HTTP-Methoden-Routing: `routes/api/users.post.dreego`

### Template-Logik erweitern
- `{#switch}` / `{#case}` / `{#default}`
- `{#let name = expr}` — Hilfsvariablen
- `{$loop}`-Variable in `{#each}`: index, first, last, even, odd
- `{#verbatim}` — Alpine/JS-Syntax unberührt lassen
- `{#fragment}` — Benannte Fragmente für HTMX-Teilladungen
- `{#slot}` / `{#fill}` — Komponenten-Slots (V1.x)

### Context
- `dreego.Context` Interface aus SSR-Context extrahieren
- SSG-Context und Wails-Context als Stubs

### Plugin-System (experimental)
> Plugin-API als *experimental* markieren. Erst nach N=3+ echten Addons als frozen deklarieren.

- `dreego.Plugin` Base-Interface
- `MiddlewareProvider` — Middleware injizieren
- `RouteRegistrar` — Routen registrieren
- `AssetProvider` — Assets via `fs.FS`
- `ContextExtender` — Context erweitern (`c.User()`)
- `Lifecycle` — OnStart / OnShutdown
- `TranspilerHook` — Custom-Tags (`<dreego:map />`)

### Security
- CSRF-Middleware (Core, opt-out via Config)
- CORS-Middleware (Core, default restriktiv)
- Session-Interface + Cookie-Store (Core)

### Form Actions
- `<form g-action="Name">` Syntax
- Auto-Parsing: Form-Fields → Go-Struct
- Validierung via Struct-Tags (`go-playground/validator`)
- Flash-Messages: `c.Flash()`, `c.Errors()`, `c.Old()`
- File-Uploads via `multipart/form-data`
- Progressiv: ohne JS, mit HTMX, mit Alpine

### Logging & Observability
- slog-Integration (strukturiertes Logging)
- Request-ID-Middleware
- Health-Check: `/health`

### CLI
- `dreego routes` — Alle Routen anzeigen
- `dreego tinker` — Go-REPL mit App-Kontext
- `dreego add <plugin>` — Plugin installieren
- Reservierte Flags: `--static`, `--wails`, `--mobile` (Coming in V2)

### Asset-Pipeline
- `<style>`-Sektionen zu einer CSS-Datei mergen
- `<script>`-Blöcke konkatenieren
- Tailwind CLI Integration (Dev JIT + Production Build)
- Cache-Busting mit Content-Hash im Dateinamen
- Compress-Middleware (gzip)

### Deployment
- `dreego build` — Single Binary
- Docker: Multi-Stage → `FROM scratch`
- Konfiguration: Env-Vars + `dreego.config.json`
- Graceful Shutdown

### Erste Addons
- `dreego-admin` — Auto-generiertes Admin-Dashboard
- `dreego-auth` — OAuth2, Passkeys, Login
- `dreego-ui` — Komponenten-Bibliothek (Shadcn-Prinzip)
- `dreego-db` — Ecto-inspirierter DB-Wrapper

---

## Phase 2 — v1.x: Stabile Erweiterungen

Keine Breaking Changes mehr. Alle Addons aus Phase 1 sind stabil.

### Developer Experience
- Built-in Hot Reload (eigenes, nicht air)
- Error Overlay im Browser
- Dev-Toolbar: Routen, `<go>`-Daten, HTMX-Requests
- Source-Maps

### Content
- Content Collections (Astro-inspiriert)
- Frontmatter in `.dreego` (YAML zwischen `---`)
- `{#md}` Block — Markdown im Template
- SEO: Meta-Tags, OpenGraph, Sitemap, robots.txt

### Performance
- Template-Caching
- ETag + Cache-Control
- Static-Asset-Caching

### Weitere Addons
- `dreego-jobs` — Queue/Job-System
- `dreego-mail` — E-Mail mit .dreego-Templates
- `dreego-storage` — S3/R2/Local File-Uploads
- `dreego-stripe` — Payments & Webhooks
- `dreego-i18n` — Mehrsprachigkeit
- `dreego-search` — Volltextsuche
- `dreego-analytics` — Privacy-friendly
- `dreego-pdf` — PDF-Generierung
- `dreego-pwa` — Service Worker, Offline
- `dreego-features` — Feature-Flags
- `dreego-notify` — Multi-Channel Notifications
- `dreego-cache` — Redis/In-Memory
- `dreego-map` — MapLibre/Leaflet
- `dreego-charts` — Chart.js
- `dreego-icons` — Icon-Bibliothek
- `dreego-markdown` — Markdown-Rendering
- `dreego-ratelimit` — DDoS-Schutz
- `dreego-devtools` — Debug-Toolbar

### Monitoring
- Prometheus-Metrics
- OpenTelemetry-Tracing

---

## Phase 3 — v2+: Neue Targets

### SSG (Static Site Generation)
- `dreego build --static` → `dist/` mit HTML
- Cloudflare Pages / GitHub Pages

### Wails v3 (Desktop + Mobile)
- `dreego build --wails` → Native App
- Gleiche `.dreego`-Komponenten für Desktop + Mobile
- iOS + Android via Wails

### TypeScript
- `<script lang="ts">` via esbuild
- Types-Sharing: Go-Struct → TS-Interface

### Weitere Ziele
- PostCSS + Tailwind JIT Built-in
- Bild-Optimierung (WebP/AVIF)
- WASM-Komponenten
- Component Registry: `dreego.dev/components`
