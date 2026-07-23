# Dreego Thinking List

**Stand:** 24.07.2026 — Phase 0.0.x aktiv, Transpiler+Layout+CSS-Scoping laufen.
**Hinweis:** Roadmap in [[../ROADMAP]]. Diese Liste enthält Detail-Fragen zu noch offenen Punkten.

## ✅ Implementiert (aus Think-List entfernt)

- Transpiler-Pipeline: Lexer → Parser → AST → CodeGen
- Alle 5 Sektionen: head, go, div, script, style
- Template-Logik: {var}, {#if}, {#each}, {#slot}, {#head}
- Layout-System: dreego/layouts/default.dreego
- CSS-Scoping: data-scope via Source-Hash
- File-based Routing: dreego/routes/*.dreego
- dreego generate: rekursiv, Hash-Cache, --force, Binary-Hash
- Docker: make up → localhost:8080
- net/http 1.22+ (kein Chi)

## 🔴 Noch offen in Phase 0

| Entscheidung | Dokument |
|---|---|
| Name dreego | [[decisions/name-dreego]] |
| Tech-Stack | [[decisions/technology-stack]] |
| Transpiler-Vorgehen | [[decisions/transpiler-pipeline]] (0.0.1: Single-Pass Scanner) |
| TypeScript in V2 | [[decisions/typescript-v2]] |
| 5 Sektionen | [[decisions/sections-in-dreego]] |
| SSR-First | [[decisions/ssr-first]] |
| Kein catch-Tag | [[decisions/no-catch-tag]] |
| File-based Routing | [[decisions/file-based-routing]] |
| SSG/Wails in V2 | [[decisions/ssg-wails-v2]] |
| Context-Design | [[decisions/context-design]] |
| Plugin-Interface | [[decisions/plugin-interface]] |
| Session-Management | [[decisions/session-management]] |
| Middleware-System | [[decisions/middleware-system]] |
| Form Actions | [[decisions/form-actions]] |

---

## 🔴 Architektur-Entscheidungen (müssen VOR Code-Start geklärt sein)

### Error-Handling-Strategie
Zu klären:
- Wie fließen Fehler von Addons ins Template?
- Stack Traces im Dev-Modus, generische Fehlerseite in Prod
- Error Boundary auf Komponenten-Ebene?
- Form-Validierungs-Feedback (Feld-Level-Errors, `dreego.Errors`)
- Flash-Messages (Erfolg/Fehler nach Redirect)
- Logging-Strategie: slog-Integration, Log-Level, strukturierte Logs

### Routing-Konventionen
Zu klären:
- `/routes/` vs `/pages/` — welcher Ordnername?
- `layout.dreego`-Konzept: Vererbung, Verschachtelung, Override
- Dynamische Segmente: `[id]`, `[...catchall]`, `[[optional]]`
- Route-Gruppen: `(marketing)/about.dreego`
- API-Routen: `/routes/api/` oder eigene Struktur?
- Middleware pro Route/Ordner
- Redirects, Rewrites in `dreego.config.json`

### Deployment-Strategie
Zu klären:
- Single Binary: Wie wird das gebaut? (`dreego build`)
- Docker-Image: Multi-Stage (test → build → deploy), FROM scratch
- Konfiguration: Environment-Variablen, `.env`, `dreego.config.json`
- Secrets-Management
- Graceful Shutdown
- Health-Check-Endpoint
- Zero-Downtime Deployments

---

## 🟡 Template-Engine & Developer Experience

### Template-Logik erweitern
- `{#each}` mit `$loop`-Variable: index, first, last, even, odd (wie Blade/Laravel)
- `{#once}` Block — Code nur einmal pro Render-Zyklus ausgeben
- `{#verbatim}` — Alpine.js/JS-Template-Syntax unangetastet lassen
- `{#fragment}` — Benannte Fragmente für AJAX-Teilladungen (wie Blade Fragments + HTMX)
- `{#md}` Block — Markdown-to-HTML im Template (MDX-like)
- `class:active={isActive}` — Svelte-Style conditional classes
- Template-Filter: `{name|upper}`, `{date|format:"2006-01-02"}` (wie Django/Jinja)
- `{#include "partial.dreego"}` — Partials
- Safe/Unsafe Output: `{html|raw}` für vertrauenswürdiges HTML
- Whitespace-Kontrolle: `{#if cond}` vs `{#if cond-}` (strip whitespace)

### Hot Reload
- Datei-Änderung → `dreego generate` → Browser Reload via SSE
- State-Erhalt über Reloads? (wie Dioxus Hot Reload)
- Nur betroffene Komponenten neu rendern, nicht ganze Seite
- Wie mit Tailwind JIT kombinieren?

### Dev-Server Features
- Error Overlay im Browser bei Compile-Fehlern
- Dev-Toolbar: Zeigt Route, `<go>`-Block Daten, HTMX-Requests, Performance
- Network-Tab: Zeigt alle HTMX-Requests und SSE-Events
- Accessibility-Checks im Dev-Mode
- Lighthouse-Integration?

### CLI-Spezifikation
```
dreego new <name>              # Projekt scaffolden
dreego dev                     # Dev-Server starten
dreego generate                # Transpiler ausführen
dreego build                   # Production-Binary
dreego build --static          # SSG (V2)
dreego routes                  # Alle Routen anzeigen
dreego add <addon>             # Addon installieren
dreego remove <addon>          # Addon entfernen
dreego version                 # Framework-Version
dreego upgrade                 # Framework upgraden
dreego generate page <name>    # Neue .dreego-Seite
dreego generate resource <name> # CRUD-Gerüst
dreego generate middleware <name>
dreego generate component <name>
dreego tinker                  # Go-REPL mit App-Kontext
```

### Konfigurations-System
Zu klären:
- `dreego.config.json` vs `dreego.toml` vs `dreego.yaml`
- Mehrere Quellen mit Priorität: CLI-Flags > Env-Vars > Config-Datei > Defaults
- Environment-spezifische Configs: `dreego.dev.json`, `dreego.prod.json`
- Reload on Change im Dev-Mode
- Secrets: `.env` (nie committed), User Secrets außerhalb Repo

### Scaffolding & Generatoren
- `dreego new` — Projekt-Template mit `routes/`, `assets/`, `dreego.config.json`
- Starter-Templates: `dreego new --template blog|saas|admin`
- Stub-Customization: `.dreego/templates/` für eigene Generator-Vorlagen
- `dreego generate page` — erzeugt `.dreego`-Datei mit `<go>` + Template + `<style>`
- `dreego generate resource` — CRUD mit Form, Validierung, DB-Query
- `dreego generate schema` — Go-Struct mit Validierungs-Tags

---

## 🟢 Features & Konzepte aus Framework-Research

### Von Leptos/Rust übernehmen
- Server-Functions-Muster: `g-submit`/`g-action` als ZENTRALES Interaktionsmodell
  - Die Go-Funktion IST die API — kein separater REST-Endpoint
  - Funktion läuft auf Server, Client ruft sie via Form/HTMX auf
  - Typ-Sicherheit über die Grenze durch Go-Struct-Tags
- Progressive Enhancement: Jedes `<form>` funktioniert ohne JS, HTMX upgraded es
- Context-Kaskade: `c.Provide("user", user)` macht Daten für Kind-Komponenten verfügbar
- Islands-Architektur explizit machen: `<script>` ist die Insel der Interaktivität
- `dreego:client` Directive: Seite explizit als client-seitig markieren (für SPAs)
- `dreego:defer` Directive: Komponente asynchron streamen (Server Islands)
- Streaming SSR: `http.Flusher` + HTMX für progressives Rendern

### Von Astro übernehmen
- "Zero JS by Default" als Design-Prinzip für Dreego übernehmen
- Content Collections: `dreego.collections.json` mit Schema → Typ-Generierung → Query-API
- Frontmatter in `.dreego`-Dateien (YAML zwischen `---` Blöcken, verfügbar als `Meta` im `<go>`)
- View Transitions: HTMX `hx-swap` + CSS View Transitions API
- Dev Toolbar: CLI-Dev-Server mit eingebautem Debug-Toolbar
- "Opt in to Complexity" als Framework-Philosophie
- Adapter-System: Core vs Deployment-Presets (Single Binary, Docker, Fly.io)
- Partial Hydration-Directives: `regeo:load`, `regeo:idle`, `regeo:visible`, `regeo:media`

### Von Phoenix übernehmen
- Ecto-Pattern: DB-Layer mit Trennung Schema ≠ Query ≠ Repo ≠ Changeset
- PubSub-System: Event-Bus für Echtzeit-Updates (Redis/NATS/In-Memory)
- Generator-Staffelung: embedded → schema → context → live → json → html
- Compile-Time Template-Validierung
- Flash-Messages: Built-in `c.Flash("success", "Gespeichert!")`
- File-Uploads: `dreego-storage` mit Progress-Tracking via SSE
- LiveView als Vorbild für SSE-Architektur (nicht WebSocket)
- Graceful Reconnection für SSE

### Von Laravel übernehmen
- Eloquent-Ebene: Naming-Conventions, Scopes, Eager Loading, Accessors/Mutators
- Validation-System: Form Requests, Field-Level Errors, `safe()`/`old()` API, 50+ Built-in Rules
- Blade Components & Slots → `<dreego:*>` Komponenten-Tags
- Queue/Job-System: Job-Middleware, Batching, Chaining, Delayed Dispatch
- Notification-System: Multi-Channel (Mail, DB, Slack)
- Ecosystem-First-Denken: Plugin-Interface muss Addon-Ökosystem ermöglichen
- Artisan-CLI als Vorbild für `dreego` CLI
- `$loop`-Variable in `{#each}`
- Stacks (`@push`/`@stack`) für Asset-Injection von Child- ins Parent-Layout

### Von Django übernehmen
- Admin-Panel: `dreego-admin` Addon (höchste Addon-Priorität)
  - Auto-Discovery, Model-Registrierung
  - `list_display`, `list_filter`, `search_fields` als Struct-Tags
  - Inline-Editing, Batch-Actions
  - Customizable Templates
- Auth-System: Password-Hashing, Reset-Flow, Permission-System
- DRF's ViewSet/Router-Pattern für `dreego-api` Addon
- Migrations: Auto-Detection aus Go-Struct-Änderungen
- Django's Dokumentations-Standard als Vorbild
- Stabilität: Keine Breaking Changes ohne lange Deprecation-Pfade

### Von C#/Blazor übernehmen
- Middleware-Pipeline mit klarer Ordnung und Dokumentation
- `UseWhen` für konditionale Middleware (Branch + Rejoin)
- Convention over Configuration: Ordnerstruktur = Routing, Namespace = Pfad
- `_Imports.dreego` für directory-scoped Defaults (Layout, Namespace, Imports)
- One-File Startup: Generiertes `main.go` minimal halten
- Dependency Injection: Funktionale Optionen oder Builder-Pattern
- Configuration-System: Multiple Sources mit Priority, Reload on Change
- `EditForm`/Form-Handling mit Validierung als First-Class
- `RenderMode`-Konzept: Static SSR, Interactive Server, Interactive Client
- Generics in Komponenten: `@typeparam T` → Go-Generics für Typed Components

---

## 🔵 Go's Single-Binary-Vorteil — In Features übersetzen

Go's Alleinstellungsmerkmal: Alles in einer Datei. Kein anderes Webframework kann das (außer Rust, das aber WASM braucht).

### Direkte Vorteile für Deployment
- `dreego deploy` — SCP die Binary auf Server, Systemd restart, fertig
- Docker-Image: `FROM scratch` — 0 MB Base-Image, nur die Binary
- Atomic Deployments: Neue Binary → `mv new old` → Systemd reload (kein Downtime)
- Cross-Compile: Von macOS Binary für Linux/ARM bauen
- Serverless-ready: Kein Cold-Start-Problem (Node braucht Module-Load, Go ist sofort da)
- Raspberry Pi / Edge: Binary läuft auf ARM, MIPS, RISC-V
- Self-contained Demos: "Hier ist die Binary, `./app` und los"
- Version-Management: `dreego upgrade` tauscht nur die Binary aus

### Features, die nur mit Single Binary möglich sind
- **Embedded Dev-Server:** `./app --dev` startet Dev-Server MIT Tailwind-Watcher
- **Self-Updating Apps:** Binary checkt GitHub Releases und updated sich selbst
- **Offline-First Deployments:** Binary enthält alles — kein CDN, kein npm, keine externen Fonts
- **Air-Gapped Installationen:** Kein Internetzugang auf dem Server nötig (Binary hat alles)
- **Portable Development:** `dreego` CLI ist selbst eine Binary — `curl ... | sh` Install
- **Snapshot-Testing:** Binary kann im Test-Modus HTML snapshots generieren und vergleichen
- **Embedded Admin Panel:** `dreego-admin` ist im Binary — `/admin` funktioniert immer, überall

### Performance-Vorteile, die kein JS-Framework hat
- Start in <10ms vs 500ms+ für Node.js/Python
- RAM: ~10MB idle vs 50-200MB für Node.js
- Kein Garbage-Collector-Stuttering (Go's GC ist sub-ms)
- Goroutines: Millionen paralleler Verbindungen auf einem Server

### Marketing-Punkte für die Website
- "Build a web app. Ship a single file."
- "No node_modules. No runtime. No excuses."
- "Deploy to a $5 VPS. Handle 10k concurrent users."
- "From zero to production in one binary."
- "Your entire app fits in a 15 MB file."

---

## 🟣 Addon-Liste (geplant für Dreego Ecosystem)

### Core Addons (offiziell maintained)

#### Auth & Security
- **dreego-auth** — OAuth2/OIDC, Passkeys, Sessions, Login/Register/Reset, 2FA, Passwordless
- **dreego-session** — Session-Store (Cookie, Redis, DB), Flash Messages

#### Daten & Storage
- **dreego-db** — DB-Wrapper mit Ecto-Pattern, Migrations, Query-Builder, Connection-Pool
- **dreego-storage** — File-Uploads (S3, R2, Local), Progress-Tracking, Image-Resize

#### UI & Komponenten
- **dreego-ui** — Komponenten-Bibliothek (Shadcn/ui-Prinzip: Copy-Paste), Tailwind, Accessible
- **dreego-admin** — Auto-generiertes Admin-Dashboard (Django-Admin als Vorbild)
- **dreego-map** — MapLibre/Leaflet Integration
- **dreego-charts** — Chart.js/Canvas-Integration
- **dreego-icons** — Icon-Bibliothek (Heroicons/Lucide)
- **dreego-markdown** — Markdown-Rendering, Frontmatter, GFM

#### Business
- **dreego-stripe** — Stripe Checkout, Webhooks, Subscription-Management
- **dreego-mail** — E-Mail-Versand (SMTP, Resend, Postmark), Templates
- **dreego-pdf** — PDF-Generierung aus HTML
- **dreego-seo** — Meta-Tags, OpenGraph, Twitter Cards, JSON-LD, Sitemap
- **dreego-analytics** — Privacy-friendly, server-seitig, kein Adblocker-Problem
- **dreego-jobs** — Hintergrund-Jobs, Cron, Queue mit Redis/DB-Backend
- **dreego-notify** — Multi-Channel Notifications (Mail, DB, Slack, Discord)
- **dreego-i18n** — Mehrsprachigkeit, Übersetzungsdateien, `t("welcome")`
- **dreego-search** — Volltextsuche (Bleve/Meilisearch/Typesense)
- **dreego-features** — Feature-Flags, A/B-Testing

#### DX & Operations
- **dreego-devtools** — Debug-Toolbar, Query-Log, Event-Log
- **dreego-pwa** — Service Worker, Offline-Caching, Push-Notifications
- **dreego-health** — Health-Check-Endpoint, Readiness-Probe, Metrics
- **dreego-cache** — Caching (In-Memory, Redis)
- **dreego-logging** — Strukturiertes Logging, Log-Drain zu externen Services

### Community-Addons (potentiell)

- **dreego-comments** — Kommentar-System (Disqus-like, self-hosted)
- **dreego-forum** — Forum/Discourse-like
- **dreego-ecommerce** — Shop-System
- **dreego-cms** — Headless CMS
- **dreego-newsletter** — Newsletter-Management
- **dreego-social** — Social-Login, Sharing
- **dreego-sitemap** — Automatische Sitemap-Generierung
- **dreego-rss** — RSS-Feed-Generierung
- **dreego-siteminder** — Uptime-Monitoring
- **dreego-abuse** — Spam/Abuse-Protection
- **dreego-embed** — oEmbed/Social-Embed
- **dreego-geo** — GeoIP, Standort-basierte Features

---

## 🟠 Zukunfts-Innovationen (Wo kann Dreego die Frontend-Welt voranbringen?)

### KI-native Entwicklung
Dreegos strukturiertes Format (5 Sektionen, Template-Logik, Go im `<go>`-Block) ist extrem AI-freundlich:
- KI kann `.dreego`-Dateien komplett generieren (Struktur ist klar definiert)
- `<go>`-Block ist pures Go — KI versteht Go besser als JSX/TS-Mischmasch
- Validierung durch Compiler: KI-generierter Code, der `dreego generate` besteht = funktioniert
- "Describe your app in natural language, get a .dreego file" — realistisch in 2026
- Dreego könnte ein `dreego ai` Subcommand haben, das LLMs einbindet

### WebAssembly Component Model
- Go kompiliert zu WASM — Komponenten könnten als WASM-Module ausgeliefert werden
- Addons als WASM: `dreego-auth` läuft in Sandbox, kann nicht auf Dateisystem zugreifen
- Plattform-unabhängige Komponenten: Gleiche `.dreego`-Datei auf Server und Client
- Wails V2: Desktop-Apps mit WASM-Komponenten

### Edge & Distributed
- Binary ist klein genug für Edge (Cloudflare Workers, Fly.io, Deno Deploy)
- Go kann als Shared Library eingebunden werden (Edge-Runtime)
- WebTransport/WebRTC für P2P-Kommunikation
- CRDTs für Offline-First mit automatischer Synchronisation
- Local-First Apps: Binary läuft lokal, sync mit Server

### Neue UI-Paradigmen
- HTML-over-WebSocket: Phoenix-LiveView-artige Architektur ohne JS
- Streaming HTML: Progressives Rendern von "above the fold" zu "below the fold"
- Adaptive Components: Gleiche Komponente rendert anders auf Mobile/Desktop/VR
- Voice-controlled Components: `<go>`-Block kann auf Voice-Events reagieren
- AR-Komponenten für WebXR

### Developer Experience Revolution
- **Time-Travel Debugging:** Go's Determinismus erlaubt Record/Replay von Requests
- **Visual State Explorer:** Alle `<go>`-Variablen im Dev-Toolbar inspizieren
- **Component Storybook:** Automatisch aus `.dreego`-Dateien generiert (kein extra Tool)
- **Zero-Config by Default:** Konventionen so gut, dass 90% der Projekte keine Config brauchen
- **Instant Previews:** Wie Vercel Previews, aber mit Single-Binary-Deploy

### Ecosystem-Innovationen
- **Component Registry:** `dreego.dev/components` — Durchsuchbare Komponenten-Bibliothek
- **One-Click Deploy:** `dreego deploy` → URL kommt zurück
- **dreego cloud:** Managed Hosting spezifisch für Dreego (wie Vercel für Next.js, aber für Go)
- **Template Marketplace:** Starter-Kits für SaaS, Blog, E-Commerce, Portfolio
- **dreego-native:** Mobile Apps via Gomobile/Wails

### Was Dreego NIE sein sollte
- Kein Electron-Ersatz (dafür gibt es Wails)
- Kein React-Killer (verschiedene Philosophien)
- Kein All-in-One-Monolith (Core klein, Addons für alles)
- Kein Vendor-Lock-in (kein "dreego cloud only")
- Kein Breaking-Change-Framework (Next.js' Schicksal vermeiden)

---

## ⚪ Testing-Strategie

### Unit-Tests
- Wie testet man eine `.dreego`-Seite?
- Request-Simulation für `<go>`-Blöcke: `dreegotest.Get("/users/1")`
- HTML-Output validieren: `assert.Contains(t, html, "<h1>Willkommen</h1>")`
- Template-Logik testen: `{#if}` mit verschiedenen Inputs
- Form-Validierung testen

### Integration-Tests
- `dreego test` — Führt `dreego generate` aus + `go test`
- Test-Helper: `dreegotest` Package mit Request-Builder, Session-Mock, DB-Mock
- Snapshot-Testing: HTML-Output snapshotten und vergleichen

### End-to-End-Tests
- Playwright-Integration
- HTMX-Interaktionen testen
- Alpine.js State-Änderungen testen
- SSE/WebSocket-Testing

---

## 🟤 Dokumentation & Community

### Dokumentation (von Tag 1)
- docs.dreego.dev: Getting Started, Guides, API Reference
- Tutorial: "Build a Blog with Dreego" (max 30 Minuten zum Durchlaufen)
- Migration Guides: Von Next.js/SvelteKit/Phoenix zu Dreego
- Examples Repo: `github.com/dreego-ecosystem/examples`
- Go-Doc für alle Core-Packages
- Recipe-Book: "How to do X in Dreego" (Auth, Upload, Search, etc.)

### Community
- GitHub Discussions
- Discord Server
- Showcase: dreego.dev/showcase
- Newsletter: Release-Notes & Tipps
- Contributing Guide
- Code of Conduct

---

## ⬜ Noch nicht kategorisiert / Brain-Dump

### Template-Engine Details
- `{#each}` mit Keys: `{#each items as item (item.ID)}` wie Svelte
- Template-Directives auch als Attribute: `<div dreego:if={cond}>` (Alternative zu `{#if}`)
- Whitespace-Kontrolle: `{#if cond-}` vs `{#if cond}`
- Self-closing Component-Tags: `<Button />`
- Spread-Props: `<Button {...props} />`
- Dynamic Components: `<dreego:component this={MyComponent} />`

### Form-Handling Details
- `g-submit` vs `g-action` — welcher Name?
- `g-submit="login"` → Handler im `<go>`-Block
- Auto-CSRF: Jedes `<form>` bekommt CSRF-Token (opt-out möglich)
- Form-Reset nach Erfolg
- Optimistische UI-Updates via Alpine.js

### Routing Details
- `dreego.config.json`: `redirects`, `rewrites`, `headers`
- Catch-All Routes: `routes/[...catchall].dreego`
- Optional Params: `routes/[[lang]]/about.dreego`
- Route-Priorität: Statisch > Dynamisch > Catch-All
- API-Routes: `routes/api/users.dreego` → `GET /api/users`
- HTTP-Method-Routing: `routes/users.post.dreego` oder via `g-action`

### Security
- CSP-Header-Generierung
- HSTS, X-Frame-Options, X-Content-Type-Options
- Subresource Integrity (SRI) für CDN-Scripts
- SQL-Injection-Prevention (Parameterized Queries by Default)
- Helmet-like Security-Header Middleware

### Observability
- Prometheus-Metrics-Endpoint
- OpenTelemetry-Tracing
- Strukturiertes Logging (slog)
- Request-ID-Tracking
- Health-Check: `/health` (liveness), `/ready` (readiness)

### Performance
- Template-Caching (gerenderte Templates cachen)
- ETag-Generierung
- Static-Asset-Caching mit `Cache-Control`
- DB-Query-Caching
- Partial-Page-Caching

### Accessibility
- ARIA-Attribute in Komponenten
- Focus-Management bei HTMX-Updates
- Skip-Links
- Screen-Reader-Testing
- Keyboard-Navigation

### Internationalisierung
- URL-basierte Lokalisierung: `/de/about`, `/en/about`
- `dreego-i18n` Addon: Übersetzungsdateien (JSON/TOML/YAML)
- `t("key")` im Template
- Pluralisierung, Datum/Zeit-Formatierung
- RTL-Support

### Build & CI/CD
- GitHub Action: `dreego build` + `dreego test`
- Docker-Image: Multi-Stage Build
- Versioning: Semantic Versioning für Framework
- Changelog-Generierung
- Release-Prozess

### Backend-Features
- Background-Job-System (dreego-jobs)
- Email-Sending (dreego-mail)
- File-Upload-Handling (dreego-storage)
- Webhook-Handling
- API-Rate-Limiting
- Request-Validation (dreego-validate)
- CORS-Handling
- Compression-Middleware (gzip/brotli)

### Monitoring & Alerting
- Error-Tracking (Sentry-like)
- Performance-Monitoring
- Uptime-Monitoring
- Usage-Analytics (opt-in, privacy-friendly)

---

## 📋 C# Blazor Spezifische Features (zur Evaluation)

- `ErrorBoundary` — Fehler in Sub-Komponenten abfangen
- `Virtualize` — Virtualisiertes Scrollen für große Listen
- `QuickGrid` — Data-Grid-Komponente
- `SectionContent`/`SectionOutlet` — Cross-Component Content Projection
- `data-permanent` — Element-Inhalt über Enhanced-Navigation erhalten
- `EnhancedNavigation` — SPA-ähnliche Page-Transitions (HTMX `hx-boost` macht das bereits)
- `PersistComponentState` — Komponenten-State über Navigation erhalten
- `@bind:culture` — Kulturabhängige Formatierung
- `@key` — Element-Identität über Updates erhalten (HTMX braucht das nicht)
- `Streaming rendering` — Async Content streamen (SSE/HTMX partials)

---

## 📋 Rust Framework Spezifische Features (zur Evaluation)

- `#[server]` Makro — Server Functions (entspricht Dreegos `g-action`)
- `view!` Makro — Compile-Time HTML-Validierung (Dreegos Transpiler)
- `#[island]` Makro — Explizite Interaktivitäts-Inseln
- `Memo`/Derived Signals — `{#let}` oder `<go>`-Berechnungen
- `provide_context`/`expect_context` — Context-API
- `<ActionForm>` — Progressiv Enhanced Forms
- `<Suspense/>` — Streaming SSR
- Multi-Thread-SSR — Go-Goroutines sind hier überlegen
- Dioxus LiveView-Mode — Server-seitiger VDOM + Diffs per WebSocket
- Dioxus Blitz Renderer — HTML/CSS nativ auf GPU gerendert (nicht relevant für Web)

---

## 📋 Solid.js Spezifische Features (zur Evaluation)

- `createResource` — Async Data als sync behandelbar (Dreego's `{#await}`)
- Keine Dependency-Arrays (Dreego braucht das Konzept gar nicht)
- Components laufen nur einmal (Dreego: `<go>`-Block läuft einmal pro Request)
- Kein Stale-Closure-Problem (Go's Scoping ist anders)
- SSR + Hydration ohne doppeltes Rendering (Dreego: HTMX partials brauchen keine Hydration)

---

## 📋 Astro Spezifische Features (zur Evaluation)

- Content Collections mit Zod-Schema → Go-Äquivalent: Struct-Tags + Validierung
- `server:defer` → Dreego: `dreego:defer` Directive
- `client:load/idle/visible/media` → Dreego: HTMX `hx-trigger` kann das
- UI-Agnostizität: Astro kann React/Vue/Svelte parallel → Dreego: Alpine.js/Datastar parallel
- View Transitions Router → Dreego: HTMX swaps + CSS View Transitions API
- Adapter-System → Dreego: Deployment-Presets

---

## 📋 Phoenix Spezifische Features (zur Evaluation)

- LiveView Hibernation nach Inaktivität → Dreego: Goroutine-Idle kann GC freigeben
- `assign_async` → Dreego: Async im `<go>`-Block
- `stream/4` → Dreego: Effiziente Listen-Updates via HTMX
- LiveComponents mit eigenem State → Dreego: `<go>` pro Komponente
- HEEx-Validierung zur Compile-Zeit → Dreego: Transpiler-Checks

---

## 📋 Laravel Spezifische Features (zur Evaluation)

- `@push`/`@stack` → Dreego: `<head>`-Block macht das bereits
- `@once` → Dreego: `{#once}` Block
- `@verbatim` → Dreego: `{#verbatim}` Block
- `$loop`-Variable → Dreego: `{#each}` erweitern
- Blade Fragments (`@fragment`) → Dreego: `{#fragment}` + HTMX `hx-select`
- Conditional Class Merging → Dreego: `class:active={isActive}`
- Service Injection (`@inject`) → Dreego: Dependency Injection
- Stub-Customization → Dreego: `.dreego/templates/`
- Tinker (REPL) → Dreego: `dreego tinker`
- Horizon (Queue-Dashboard) → Dreego: `dreego-jobs` Dashboard
- Telescope (Debug-Toolbar) → Dreego: `dreego-devtools`
- Pennant (Feature Flags) → Dreego: `dreego-features`

---

## 📋 Django Spezifische Features (zur Evaluation)

- Admin-Panel → `dreego-admin` (HÖCHSTE Addon-Priorität)
- `list_editable` (Inline-Editing in Admin-Liste)
- `autocomplete_fields` (Select2-Integration)
- `prepopulated_fields` (Slug aus Title)
- Facetten-Counts in Admin-Filtern
- Signals (`pre_save`, `post_delete`) → Model-Events
- F-Expressions (DB-Level-Operationen) → DB-Query-Builder
- Django's Migrations-Auto-Detection → `dreego-db` Migrations
- Browsable API (DRF) → `dreego-api` Addon
