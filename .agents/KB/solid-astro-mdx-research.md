# Solid.js, Astro & MDX — Deep-Dive Research

**Datum:** 23.07.2026
**Quellen:** docs.solidjs.com, docs.astro.build, mdxjs.com

---

## 1. Solid.js

### Was macht Solid anders als React/Svelte?

| Aspekt | React | Svelte | Solid |
|--------|-------|--------|-------|
| Rendering | Virtual DOM → Diff → DOM | Compiler → direktes DOM-Update | Fine-grained Reactivity → kein VDOM, kein Compiler |
| Component Model | Components re-rendern komplett | Compiler weiß, was sich ändert | Components laufen EINMAL, nur reaktive Bindings updaten |
| State | `useState` (Hook, immutable) | `$state()` (Rune, mutable via Compiler) | `createSignal()` (getter/setter, mutable) |
| Philosophie | "UI = f(state)" | "Compiler macht's effizient" | "Jede Zelle aktualisiert sich selbst" |

**Kern:** Solid-Komponenten sind **keine Render-Funktionen** — sie werden einmal ausgeführt. Der Return-Wert ist echtes DOM. Nur Signal-gesteuerte Attribute/TextNodes updaten sich bei Änderungen. React-Komponenten dagegen laufen komplett neu bei jedem State-Change.

### Wie funktionieren Signals in Solid (Fine-Grained Reactivity)?

**Mechanismus (Observer-Pattern):**

1. `createSignal(initialValue)` → `[getter, setter]`
2. `createEffect(fn)` registriert `fn` als Subscriber
3. Beim Aufruf von `getter()`: fügt aktuellen Subscriber zur Signal-Subscriber-Liste hinzu
4. `setter(newValue)`: notifiziert alle Subscriber (nur bei Wert-Änderung)
5. **Kein Dirty-Checking, kein VDOM-Diff** — Signal weiß exakt, welche DOM-Knoten es updaten muss

**Wichtige Eigenschaften:**
- **Synchron:** Tracking läuft synchron. `setTimeout` im Effect → kein Tracking
- **Memos:** `createMemo()` cached abgeleitete Werte (wie `$derived` in Svelte)
- **Resources:** `createResource(fetcher)` wandelt Async in Sync (wie Sveltes `{#await}`)
- **Stores:** `createStore()` erzeugt Proxy-basierte Signals für verschachtelte Objekte
- **Kein Stale Closure-Problem:** Weil Components nur einmal laufen und Signals direkt das DOM updaten

**Das ist relevant für Dreego, weil:**
- Dreego's Architecture (SSR + HTMX partials) ist im Prinzip dasselbe wie Signals, nur auf Server-Ebene
- State ändert sich → nur das betroffene HTML-Fragment wird neu gerendert und via HTMX ausgetauscht
- Kein Virtual DOM nötig — direktes DOM-Update durch HTML-Swap

### SolidStart Features

SolidStart ist das Metaframework (wie Next.js für React, SvelteKit für Svelte):

- **File-based Routing** — `routes/` Verzeichnis, nested layouts, dynamic routes `[id].tsx`, catch-all `[...slug].tsx`, route groups `(groupName)/`, escaping nested routes `users(details)/`
- **Multiple Rendering Modes:** CSR, SSR (Sync, Async, Streaming), SSG (Pre-rendering via `crawlLinks: true`)
- **Server Functions** — `"use server"` Directive, Datenbank-Zugriff ohne API-Endpoint. Integriert mit `query()` (Solid Router)
- **API Routes** — GET/POST/etc. Handler in `routes/`
- **Middleware** — Request/Response-Interception
- **Sessions, Auth, WebSocket Endpoints**
- **Vinxi/Nitro** — agnostischer Bundler + Server Runtime (kein Vendor Lock-in)
- **Deployment Presets:** Vercel, Netlify, Cloudflare, Node
- **Kein eigener Router** — verwendet `@solidjs/router` (trennt Meta vom Router)

### Solid's Killer Features

1. **True Reactivity ohne Compiler oder VDOM** — minimalster Runtime-Code
2. **Components laufen nur einmal** — kein re-rendering, nur punktuelle DOM-Updates
3. **Kein Stale Closure Problem** — weil Components nicht re-rendern
4. **Keine Dependency Arrays** — kein `useEffect([dep])`, kein `useMemo`, kein `useCallback`
5. **SSR + Hydration ohne doppeltes Rendering** — Solid hydratisiert nur Daten-Bindings, nicht ganze Komponenten
6. **Isomorphe Server Functions** — gleiche Funktion läuft auf Server und Client (mit `"use server"`)
7. **Extrem kleine Bundle-Größe** — <7 KB gzipped für Solid Core
8. **JSX ohne React** — JSX wird zu echten DOM-Operationen compiliert (nicht zu React.createElement)

---

## 2. Astro

### Islands Architecture

Astro hat die **Islands Architecture** popularisiert (ursprünglich von Etsy/Katie Sylor-Miller, 2019, dann Jason Miller/Preact, 2020):

**Konzept:**
- Die Seite ist ein **Meer aus statischem HTML**
- Interaktive Komponenten sind **Inseln** von JavaScript im statischen Meer
- Jede Insel wird **isoliert** geladen und hydratisiert
- Inseln können **parallel** laden — eine langsame Image-Carousel blockiert nicht den Header

**Zwei Typen von Inseln:**
1. **Client Islands** — interaktive UI-Komponenten (React/Svelte/Vue/Solid)
2. **Server Islands** (`server:defer`) — dynamische Server-Inhalte, die parallel zum Hauptinhalt streamen (z.B. User-Avatar, Produktbewertungen)

### Partial Hydration — Client Directives

Astro's Kern-Mechanismus: **JavaScript wird standardmäßig entfernt.** Nur explizit markierte Komponenten bekommen JS. Die Directives:

| Directive | Verhalten |
|-----------|-----------|
| `client:load` | Sofort laden (höchste Priorität) |
| `client:idle` | Wenn Browser idle ist (requestIdleCallback) |
| `client:visible` | Wenn Komponente in Viewport kommt (IntersectionObserver) |
| `client:media="(max-width: 50em)"` | Nur bei passender Media Query |
| `client:only="react"` | Nur client-seitig, kein SSR |

### Content Collections

Astro's Content-Management-System:

- **`src/content.config.ts`** — Zentrale Konfiguration aller Collections
- **Zod-Schema** pro Collection → TypeScript Typen automatisch generiert, Editor-Intellisense
- **Built-in Loaders:**
  - `glob()` — Verzeichnis von Markdown/MDX/Markdoc/JSON/YAML/TOML-Dateien
  - `file()` — Einzelne Datei mit Array von Einträgen
- **Custom Loader API** — CMS, Datenbank, API → alles integrierbar
- **Reference System** — Collections können aufeinander verweisen (`reference('authors')`)
- **`getCollection()`, `getEntry()`** — Typisierte Query-API
- **`render()`** — Rendert Markdown/MDX zu HTML + `<Content />` Component
- **Filter-API** — `getCollection('blog', ({data}) => data.draft !== true)`
- **Live Collections** — Für Echtzeit-Daten (Live-Updates ohne Rebuild)
- **Route Generation** — Aus Collection-Einträgen automatisch Seiten generieren

### Was macht Astro einzigartig (94% Satisfaction)?

1. **Zero-JS Default** — Es ist unmöglich, aus Versehen JavaScript zu schicken
2. **MPA statt SPA** — Multi-Page-Architektur statt Single-Page-App. Seitenwechsel laden neue HTML-Dokumente (schneller, einfacher, SEO)
3. **Content-Driven Design** — Für Content-Seiten optimiert (Blog, Marketing, Docs, E-Commerce), nicht für Web-Apps
4. **Server-First** — Rendering passiert auf dem Server, nicht im Browser (wie PHP/Laravel/Rails)
5. **"Opt in to Complexity"** — Starte mit HTML+CSS, füge bei Bedarf Frameworks/JS hinzu
6. **UI-Agnostisch** — React, Preact, Svelte, Vue, Solid, HTMX, Web Components — alles parallel nutzbar
7. **`.astro` Syntax** — Superset von HTML: jedes gültige HTML ist gültiges Astro-Template
8. **View Transitions Router** — SPA-ähnliche Animationen zwischen MPAs
9. **Dev Toolbar** — Integrierte Dev-Tools im Browser
10. **Adapter-System** — Trennung von Framework und Deployment-Ziel

---

## 3. MDX (und Astro's MDX-Integration)

### Wie funktioniert Markdown-to-Component Rendering?

MDX = Markdown + JSX:

```
# Hello, world!                           ← Markdown
<div className="note">                    ← JSX (Component!)
  > Some notable things in a block quote! ← Markdown IN JSX
</div>
```

**Verarbeitungskette:**
1. **Parse** MDX-Text → MDAST (Markdown AST) + JSX-Nodes
2. **Transform** Remark-Plugins modifizieren den MDAST
3. **Compile** → JavaScript (JSX wird zu `createElement`-Calls, Markdown zu HTML-Strings)
4. **Evaluate** → Ausführung im JS-Runtime (React/Preact/Vue)

**Wichtige Features:**
- **Import/Export** — `import` und `export` Statements im Markdown
- **Expressions** — `{Math.PI * 2}` im Markdown
- **Custom Components** — HTML-Elemente durch eigene Komponenten ersetzen: `export const components = {blockquote: CustomBlockquote}`
- **Frontmatter** — YAML/TOML am Anfang der Datei
- **ESM Support** — Volle JavaScript-Modul-Syntax

### Astro's MDX Integration

Astro erweitert MDX:

- **Content Collections + MDX** — `.mdx` als Collection-Einträge mit Zod-Schema und Typisierung
- **Astro Components in MDX** — `.astro` Components direkt in `.mdx` importieren und nutzen
- **Custom Components Mapping** — `<Content components={{h1: Heading}} />`
- **Eigener MDX-Compiler** — `@astrojs/mdx` mit `recmaPlugins`, `optimize` Option
- **Frontmatter als First-Class** — `{frontmatter.title}` direkt im MDX verwendbar
- **Separate Processor** — MDX kann anderen Markdown-Processor als `.md` Dateien nutzen
- **Hybrid-Modus** — Statische Seiten + MDX-Content: im Prinzip ein Headless CMS in Git

---

## Für Dreego: Adopt / Don't Adopt

### Solid.js

**ADOPTIEREN:**
- [x] **Signals als konzeptionelles Modell** — Dreego's SSR + HTMX ist funktional äquivalent: State-Änderung → nur betroffenes HTML-Fragment updaten. Das mentale Modell ist dasselbe.
- [x] **Kein Virtual DOM** — Solid beweist, dass VDOM nicht nötig ist. Dreego macht das auf Server-Ebene genauso (direktes HTML-Rendering).
- [x] **Server Functions Modell** — `"use server"` = Dreegos `<go>` Block. Der `<go>`-Block läuft ausschließlich serverseitig und interagiert mit DB/APIs — exakt dasselbe Konzept.
- [x] **Isomorpher Code** — SolidStart's Design-Prinzip: Code läuft auf Server und Client. Dreego könnte `dreego generate` nutzen um Go-Code zu generieren der sowohl Server- als auch Client-Logik enthält (via GopherJS/WASM in V2).
- [x] **Keine Dependency Arrays** — Dreego's Template-Engine braucht keine `useEffect`-Äquivalente. Das `<go>`-Block-Modell ist einfacher.
- [x] **Resource-Pattern** — `createResource` = Dreego's `{#await}` Tag. Async-Daten als synchron behandelbar machen.

**NICHT ADOPTIEREN:**
- [ ] **Client-seitiges Signal-System** — Dreego setzt auf HTMX + Alpine.js. Ein eigenes JS-Signal-System wäre Redundanz und unnötiger JS-Ship.
- [ ] **JSX-Syntax** — Dreego nutzt HTML-Template-Syntax (wie Svelte), nicht JSX. JSX ist zu stark an JavaScript/React gebunden.
- [ ] **Vinxi/Nitro als Build-Tool** — Dreego hat Go's `go build`. Braucht kein JS-Build-Tool.
- [ ] **SolidStart als Architektur-Vorbild** — SolidStart ist weniger etabliert als SvelteKit/Next.js/Astro. Für File-based Routing lieber SvelteKit/Astro als Vorbild nehmen.

### Astro

**ADOPTIEREN:**
- [x] **Islands-Architektur** — Das ist DER entscheidende Insight für Dreego. `.dreego` Seiten sind standardmäßig statisch (Null JS). Interaktive "Inseln" werden via HTMX/Alpine.js deklariert. Kein Framework-JS auf der Seite außer dem, was explizit als interaktiv markiert wurde.
- [x] **"Zero JS by Default"** — Dieses Prinzip ist perfekt für Dreego. Kein Framework-JS Code im Output, nur das was der Entwickler explizit will.
- [x] **Partial Hydration-Konzept → "Partial HTML Swap"** — Astro's `client:visible` etc. lassen sich auf Dreego übertragen: HTMX-Partials können lazy, on-visible, on-interaction geladen werden.
- [x] **Content Collections** — Dreego braucht dringend ein ähnliches Konzept für `.dreego` Seiten. Eine `dreego.collections.toml` mit Schema-Definition, automatischer Typ-Generierung, Query-API. Das wäre ein Killer-Feature für Content-Seiten.
- [x] **MPA-Ansatz (Multi-Page-App)** — Dreego ist per Definition MPA. Astro beweist, dass MPA der richtige Ansatz für Content-Seiten ist und SPAs overkill sind. Das validiert Dreegos Grund-Architektur.
- [x] **Adapter-System** — Astro's Trennung von Framework-Core und Deployment-Adapter ist architektonisch sauber. Dreego könnte dasselbe machen: Core-Framework + Deployment-Presets (Single Binary, Docker, Fly.io, VPS).
- [x] **".astro Syntax als HTML-Superset" → Dreegos Template-Syntax** — `.dreego` Template sollte sich an HTML anlehnen (wie es bereits der Fall ist). Jedes gültige HTML ist gültiges `.dreego` Template — das senkt die Einstiegs-Hürde massiv.
- [x] **View Transitions** — Astro's MPA+SPA-Animationen. Dreego könnte ähnliches mit HTMX's `hx-swap` + CSS View Transitions API erreichen.
- [x] **"Opt in to Complexity"** — Design-Phrase für Dreego übernehmen. Starte einfach, füge Komplexität nur bei Bedarf hinzu.
- [x] **Dev Toolbar** — Ein CLI-Dev-Server mit eingebautem Debug-Toolbar (zeigt aktuelle Route, `<go>`-Block Daten, HTMX-Requests) wäre ein Differenzierungsmerkmal.
- [x] **Server Islands (`server:defer`)** — Für Dreego hieße das: Teile einer Seite können asynchron streamen. `<div dreego:defer>` → wird separat gerendert und per SSE/HTMX nachgeladen. Ideal für langsame DB-Queries.
- [x] **UI-Agnostizität** — Astro unterstützt parallel React, Vue, Svelte. Dreego könnte im `<script>`-Block mehrere Client-Frameworks erlauben (Alpine, Datastar, Petite-Vue). Für V2 sogar WASM-Komponenten.
- [x] **CSR-Fallback (`client:only`)** — Manche Seiten brauchen volle Client-Interaktivität. Dreego könnte `regeo:client="true"` als Page-Level-Directive haben, um eine Seite komplett client-seitig zu machen.

**NICHT ADOPTIEREN:**
- [ ] **Server-First Rendering in JS** — Astro rendert das auf Node/Deno. Dreego macht das in Go. Der Mechanismus ist fundamental anders (Compiler vs. Runtime-Renderer).
- [ ] **Vite-basierter Dev-Server** — Dreego braucht keinen JS-Bundler. Der Dev-Server ist Go-nativ.
- [ ] **npm/Node-Ökosystem** — Dreego hat Go-Module. Astro's Stärke ist die Integration mit dem JS-Ökosystem. Dreego kann das nicht replizieren (und will es nicht).
- [ ] **Astro Components auf Client und Server** — Dreego's `<go>`-Block und Template sind strikt getrennt (Server vs. Client). Astro's `.astro` Components sind eine Mischform.
- [ ] **Integrations-API im JS-Stil** — Dreego braucht ein Plugin-System, aber in Go-Idiomatik (Interfaces, nicht JS-Funktionen).

### MDX

**ADOPTIEREN:**
- [x] **Markdown + Components** — Dreego könnte `.dreego` Dateien ein `{#md}` Block-Konzept geben, in dem Markdown mit Template-Komponenten gemischt wird.
- [x] **Custom Component Mapping** — Wie MDX's `export const components = {h1: MyHeading}` → Dreego könnte Templates definieren, die HTML-Tags in `.dreego` Components umwandeln.
- [x] **Frontmatter als First-Class** — Jede `.dreego` Datei könnte optional Frontmatter (YAML/TOML) haben, das im `<go>`-Block als `Meta`-Variable verfügbar ist.
- [x] **Content Collections + Schema** — Wie MDX in Astro Collections: `.dreego` Seiten in einer Collection mit Zod-ähnlichem Schema (Go struct tags?), automatische Typ-Generierung.
- [x] **Markdown in Template-Blöcken** — Template-Code (`{#if}`, `{#each}`) innerhalb von Markdown-Content. "Literate Programming" für Webseiten.
- [x] **ESM/Import-Pattern → Dreego Addon-System** — MDX's `import` Statement → Dreego's `{#use addon}` im Template. Komponenten aus Addons importieren.
- [x] **Remark/Rehype Plugin-System → Dreego Markdown Pipeline** — Eine Pipeline von Markdown-Transformern (Syntax Highlighting, Table of Contents, Link-Rewriting). Als Go-Pipeline implementierbar.

**NICHT ADOPTIEREN:**
- [ ] **JSX in Markdown** — Funktioniert nur mit JS-Runtime. Dreego müsste eine eigene Component-Syntax für Markdown erfinden.
- [ ] **Evaluate-Schritt** — MDX evaluated JS-Code zur Laufzeit. Dreego compiliert alles vorher — kein Runtime-Evaluate. Sicherer, aber weniger dynamisch.
- [ ] **npm-Plugin-Ökosystem** — MDX's Stärke ist das Remark/Rehype-Ökosystem. Dreego kann das nicht 1:1 replizieren. Eigenes Markdown-Plugin-System in Go wäre nötig.
- [ ] **MDX als Dateiformat** — Dreego hat `.dreego` als Format. MDX wäre ein zweites, konkurrierendes Format. Besser: MDX-ähnliche Features in `.dreego` integrieren.

---

## Konkrete Vorschläge für Dreego

### 1. Islands-Architektur für `.dreego`
```
<!-- Seite ist standardmäßig statisch -->
<div>
  <h1>{title}</h1>
  
  <!-- Interaktive Insel via HTMX -->
  <div hx-get="/api/counter" hx-trigger="load">
    <!-- lädt nach -->
  </div>
  
  <!-- Client-seitige Insel via Alpine -->
  <div x-data="{ open: false }">
    <button @click="open = !open">Toggle</button>
  </div>
</div>
```

### 2. Content Collections via `dreego.collections.toml`
```toml
[collections.blog]
source = "content/blog"
pattern = "**/*.dreego"
[schema]
title = "string"
date = "time"
draft = "bool"
```

### 3. Server Islands via `dreego:defer`
```html
<!-- Paralleles Server-Rendering, via SSE gestreamt -->
<UserAvatar dreego:defer user="{userId}" />
<ProductReviews dreego:defer product="{productId}" />
```

### 4. Frontmatter in `.dreego`
```html
---
title: My Blog Post
date: 2026-07-23
draft: false
---

<go>
    // Meta.title, Meta.date, Meta.draft verfügbar
</go>

<h1>{Meta.title}</h1>
```

### 5. Markdown in Templates (MDX-like)
```html
<go>
    posts := query.Posts()
</go>

{#md}
# Willkommen auf meinem Blog

Hier sind die neuesten Posts:

{#each posts as post}
- [{post.Title}](/post/{post.Slug})
{#end}
{/md}
```
