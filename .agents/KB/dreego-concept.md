
---
type: Reference
title: Dreego Concept — Gemini Chat (PDF-Quelle)
description: Core concept of Dreego framework: transpiler, 5 sections, file-based routing, addon ecosystem
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Dreego Concept — Gemini Chat (PDF-Quelle)

> Quelle: `Go-Framework wie Svelte bauen.pdf` — Gemini Chat vom 23.07.2026
> Die PDF wurde via pdftotext extrahiert und in diese Markdown-Datei überführt.

---

## Einleitung

Ein eigenes Framework zu bauen, das Svelte-artige Syntax, File-based Routing und Server-Side Rendering (SSR) direkt in Go vereint.

Da Go eine kompilierte Sprache ist, gibt es zwei Wege:

| Eigenschaft            | Weg A: Eigenes Transpiler-Tool   | Weg B: Pure-Go DSL            |
|------------------------|----------------------------------|-------------------------------|
| Dateiendung            | `.gvelte`                        | `.go`                         |
| Generierungsschritt    | Ja (vor `go run`)                | Nein (rein Go-nativer Code)   |
| Developer Experience   | Echte HTML/Svelte-Syntax         | Native Go-Typisierung & AC    |
| Go-Ökosystem           | Benötigt Parser, AST & Generator | Funktioniert 100% Out of Box  |

---

## Weg A: Transpiler (Compiler-Ansatz)

Der gewählte Ansatz für Dreego. `.dreego`-Dateien werden via `dreego generate` in Go-Code kompiliert.

### Vorteile
- **Maximale Performance:** Keine Disk-IO, kein Laufzeit-Parsing
- **Single Binary:** Alles via `//go:embed` ins Binary
- **Compile-Time Safety:** Fehlerhafte Templates brechen `go build` ab
- **DevX durch KI:** KI führt `dreego generate` im Hintergrund aus

### Der Transpiler-Workflow
```
.dreego Datei  →  [Lexer/Parser]  →  [AST]  →  [Go-Code Generator]  →  .go Datei
```

---

## Der Name: Dreego

- Aussprache: Go-Ree (mit langem e)
- Dateiendung: `.dreego`
- Package: `dreego`
- Grund: TTS-Freundlichkeit (vorher "Gvelte" wurde zu "Gor-ee" verzerrt)

---

## Die 5 Sektionen einer .dreego-Datei

### 1. `<head>` — Komponenten-spezifische Meta-Tags & Assets
```html
<head>
    <script src="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.js"></script>
    <link href="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.css" rel="stylesheet" />
</head>
```
- Wird nur geladen, wenn die Komponente tatsächlich gerendert wird
- Perfekt für Addons (dreego-map braucht Mapbox nur auf der Map-Seite)

### 2. `<go>` — Server-seitiger Go-Code
```html
<go>
    userID := c.Param("id")
    user, err := db.GetUser(userID)
</go>
```
- Läuft NUR auf dem Go-Server vor dem Rendern
- Direkter DB-Zugriff, keine API-Schicht nötig

### 3. Template (HTML) — Das Markup
```html
<div class="user-card">
    <h1>Hallo, {user.Name}!</h1>
    <button id="map-btn">Karte zentrieren</button>
</div>
```

### 4. `<script>` — Client-seitiges JavaScript
```html
<script lang="ts">
    document.getElementById("map-btn")?.addEventListener("click", () => {
        console.log("Karte wird zentriert...");
    });
</script>
```
- V1: Vanilla JS (kein TypeScript, kein Compiler)
- V2: TypeScript via esbuild-Integration

### 5. `<style>` — Scoped CSS
```html
<style>
    .user-card {
        padding: 1rem;
        background: #f4f4f4;
    }
</style>
```
- Klassen werden automatisch mit Hashes versehen (Scoping)

---

## Template-Logik (V1)

### {#if} / {#else if} / {#else}
```html
{#if user.IsAdmin}
    <a href="/admin">Dashboard</a>
{#else if user.IsLoggedIn}
    <a href="/profile">Mein Profil</a>
{#else}
    <a href="/login">Anmelden</a>
{/if}
```
→ Kompiliert zu `if condition { ... } else { ... }`

### {#switch} / {#case} / {#default}
```html
{#switch order.Status}
    {#case "pending"}
        <span class="badge yellow">In Bearbeitung</span>
    {#case "shipped"}
        <span class="badge blue">Unterwegs</span>
    {#default}
        <span class="badge gray">Unbekannt</span>
{/switch}
```
→ Kompiliert zu nativem Go `switch`

### {#each} / {#else}
```html
<ul>
    {#each users as user, index}
        <li>#{index + 1}: {user.Name}</li>
    {#else}
        <li>Keine Benutzer gefunden.</li>
    {/each}
</ul>
```
→ {#else} bei leerer Liste — kein `if len(users) == 0` nötig

### {#await} (Async/Streams)
```html
{#await fetchUserData()}
    <div class="skeleton-loader">Lädt...</div>
{#then user}
    <p>Willkommen, {user.Name}!</p>
{#catch err}
    <p class="error">{err.Error()}</p>
{/await}
```
→ Nutzt Go-Channels oder JS Promises

### {#slot} / {#fill} (Komponenten-Slots)
```html
<!-- In Card.dreego -->
<div class="card">
    <div class="card-header">
        {#slot header}
            <h3>Standard Titel</h3>
        {/slot}
    </div>
    <div class="card-body">
        {#slot}
        {/slot}
    </div>
</div>
```

```html
<!-- Verwendung -->
<Card>
    {#fill header}
        <h3 class="gold">Premium Benutzer</h3>
    {/fill}
    <p>Das ist der Inhalt.</p>
</Card>
```

### {#let} (Hilfsvariablen)
```html
{#let fullName = user.FirstName + " " + user.LastName}
{#let isGoldCustomer = user.OrdersCount > 50}
```

---

## File-based Routing mit Chi

### Ordnerstruktur
```
src/
└── routes/
     ├── get.dreego           →  /
     ├── about.dreego           →  /about
     └── users/
         └── [id].dreego        →  /users/:id
```

### Routing-Engine
- `dreego generate` scannt den `routes/` Ordner
- Generiert eine zentrale `dreego_router.go`:
```go
// CODE GENERATED BY DREEGO. DO NOT EDIT.
package main

import "net/http"

func RegisterDreegoRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/", RenderIndexPage)
    mux.HandleFunc("/users/{id}", RenderUsersPage)
}
```
→ Kein Laufzeit-Scanning, keine Disk-IO, alles im Binary

---

## Der Core-Stack von Dreego

| Bereich            | Wahl                      | Begründung                                    |
|--------------------|---------------------------|-----------------------------------------------|
| HTTP-Router        | go-chi/chi                | Schnell, 100% Go-Stdlib kompatibel            |
| Template Engine    | Dreego Transpiler          | Herzstück — muss eigen sein                   |
| Interaktivität     | HTMX + Alpine.js          | LiveView-Gefühl ohne SPA-Komplexität          |
| SSE/Real-time      | Datastar (optional)       | Streamt DOM-Updates via SSE                   |
| CSS                | Tailwind CLI (embedded)   | Kein eigener CSS-Parser nötig                 |
| Validierung        | go-playground/validator   | Struct-Tag Validierung                        |
| Binary Packaging   | `embed` (Go Stdlib)       | Native Go, kein Extra-Tool                    |
| JS Bundling (V2)   | esbuild                   | Als Go-Bibliothek einbindbar                  |

---

## Wann ist Dreego besser als die Konkurrenz?

### Besser als React / Vue / Svelte (SPAs)
- Keine JS-Build-Hölle (kein node_modules)
- 0 MB JS-State-Synchronisierung (State nur auf Server)
- Perfektes SEO + FCP auch auf schwachen Geräten

### Besser als Phoenix (Elixir)
- Single-Binary Deployment (keine BEAM VM, kein Node.js)
- Kompilier- und Typensicherheit (Go-Compiler verhindert Runtime-Errors)
- Bruchteil des RAM-Verbrauchs, Start in Millisekunden

### Besser als reine Go-Templates / Templ
- Echte Framework-Erfahrung (Routing, Actions, Layouts, Dev-Server)
- Kein JS-Glue-Code (HTMX/Datastar direkt integriert)

---

## Addon/Plugin-Ökosystem

### Architektur
Ein Goree-Addon ist ein Go-Package, das das `dreego.Plugin` Interface erfüllt:

```go
type Plugin interface {
    Name() string
    RegisterRoutes(mux http.Handler)
    Middlewares() []func(http.Handler) http.Handler
    Assets() *embed.FS
}
```

### Nutzung in main.go
```go
app := dreego.New()
app.UsePlugin(auth.New("super-secret-key"))
app.Listen(":8080")
```

### Addon-Ideen (vollständige Liste)

| Addon               | Zweck                                                   |
|---------------------|---------------------------------------------------------|
| dreego-auth          | Sessions, OAuth, Passkeys, Login/Register               |
| dreego-map           | MapLibre/Leaflet Integration                            |
| dreego-admin         | Auto-generiertes Admin-Dashboard                        |
| dreego-seo           | Meta-Tags, OpenGraph, Sitemap                           |
| dreego-db            | DB-Integration (SQLite, Turso, PostgreSQL)              |
| dreego-analytics     | Privacy-friendly, serverseitig                          |
| dreego-i18n          | Mehrsprachigkeit                                        |
| dreego-mail          | E-Mail-Versand mit .dreego-Templates                     |
| dreego-stripe        | Stripe Checkout & Webhooks                              |
| dreego-storage       | Dateiuploads (S3, R2, local)                            |
| dreego-jobs          | Hintergrundaufgaben & Cronjobs                          |
| dreego-ui            | Shadcn-ähnliche UI-Komponenten                          |
| dreego-search        | Volltextsuche (Bleve/Meilisearch)                       |
| dreego-pdf           | PDF-Generierung aus .dreego-Templates                    |
| dreego-pwa           | Progressive Web App                                     |

### Vorteile des Go-Addon-Systems
- Keine Dependency-Hölle (go.mod löst sauber auf)
- Compile-Time Safety: Build bricht bei Fehlern ab
- Tree-Shaking: Nicht genutzter Code wird rauskompiliert

---

## V1 Blueprint (MVP)

1. **Transpiler:** Liest `.dreego` → generiert `.go`
2. **3 Sektionen:** `<go>`, HTML-Template, `<style>`
3. **Template-Logik:** `{#if}` und `{#each}`
4. **Chi-Router Wrapper:** Ordnerpfade → HTTP-Routen
5. **Single Binary:** `//go:embed` für Assets

## V2 Ausblick

- TypeScript Support via esbuild
- SSG (Static Site Generation)
- Wails Integration (Desktop-Apps)
- `{#await}`, `{#slot}`, `{#switch}`
- Auto-Form Binding
- Inline API Endpoints

---

## Fehlerbehandlung in Dreego

**Kein `<catch>`-Tag!** Fehler werden in Go-Idiomen behandelt:

```html
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

<div class="profile">
    {#if hasError}
        <p class="error">Benutzer konnte nicht geladen werden.</p>
    {#else}
        <h1>Hallo, {user.Name}!</h1>
    {/if}
</div>
```

---

## Tailwind Class Merging für dreego-ui

```html
<go>
    finalClasses := dreego.MergeClasses(
        "bg-blue-500 text-white py-2 px-4 rounded",
        props.Class,
    )
</go>

<button class="{finalClasses}" {attrs}>
    {#slot}Button Text{/slot}
</button>
```

Nutzung:
```html
<Button class="bg-red-600 mt-4" disabled>Löschen</Button>
```

---

## A/B Testing via Go-Logik

```html
<go>
    userIP := c.Request.RemoteAddr
    showPriceVariantA := isOddIP(userIP)
</go>

<div>
    {#if showPriceVariantA}
        <div class="price-tag">Sonderpreis: 19 Euro</div>
    {#else}
        <div class="price-tag">Regulärer Preis: 29 Euro</div>
    {/if}
</div>
```

---

## Dreego vs Svelte vs React — Zusammenfassung

- **Fühlt sich an wie Svelte:** `.dreego`-Dateien, HTML-first Syntax
- **Strukturiert wie React:** Komponenten, Props, Children-Slots
- **Robuster als beide:** Kein Node.js, Single Binary, Go-Typensicherheit

---

## Quellen aus dem Chat

1. Cross-Platform Go: Runtime Checks vs. Build Tags
2. Customizing Go Binaries with Build Tags - DigitalOcean
3. Gomacro: code generation made easy and fun
4. How to Use go:embed in Go - JetBrains
5. Chi vs Encore.go - Which Go Backend Framework to Choose in 2026
6. Datastar: Web Framework for the Future?
7. irgo - Native Apps with Go + Datastar
