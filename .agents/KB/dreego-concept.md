
---
type: Reference
title: Dreego Concept — Gemini Chat (PDF Source)
description: Core concept of Dreego framework: transpiler, 5 sections, file-based routing, addon ecosystem
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Dreego Concept — Gemini Chat (PDF Source)

> Source: `Building a Go Framework like Svelte.pdf` — Gemini Chat from 2026-07-23
> The PDF was extracted via pdftotext and transferred to this Markdown file.

---

## Introduction

Building a custom framework that combines Svelte-like syntax, file-based routing, and server-side rendering (SSR) directly in Go.

Since Go is a compiled language, there are two paths:

| Feature              | Path A: Custom Transpiler Tool   | Path B: Pure-Go DSL            |
|----------------------|----------------------------------|-------------------------------|
| File extension       | `.gvelte`                        | `.go`                         |
| Generation step      | Yes (before `go run`)            | No (pure Go-nativer Code)     |
| Developer Experience  | Real HTML/Svelte syntax          | Native Go typing & autocompletion |
| Go Ecosystem          | Needs parser, AST & generator   | Works 100% out of the box     |

---

## Path A: Transpiler (Compiler Approach)

The chosen approach for Dreego. `.dreego` files are compiled into Go code via `dreego generate`.

### Benefits
- **Maximum Performance:** No disk IO, no runtime parsing
- **Single Binary:** Everything via `//go:embed` into the binary
- **Compile-Time Safety:** Faulty templates break `go build`
- **DevX through AI:** AI runs `dreego generate` in the background

### The Transpiler Workflow
```
.dreego File  →  [Lexer/Parser]  →  [AST]  →  [Go Code Generator]  →  .go File
```

---

## The Name: Dreego

- Pronunciation: Go-Ree (with long e)
- File extension: `.dreego`
- Package: `dreego`
- Reason: TTS-friendliness (previously "Gvelte" was distorted to "Gor-ee")

---

## The 5 Sections of a .dreego File

### 1. `<head>` — Component-specific Meta Tags & Assets
```html
<head>
    <script src="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.js"></script>
    <link href="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.css" rel="stylesheet" />
</head>
```
- Only loaded when the component is actually rendered
- Perfect for addons (dreego-map needs Mapbox only on the map page)

### 2. `<go>` — Server-side Go Code
```html
<go>
    userID := c.Param("id")
    user, err := db.GetUser(userID)
</go>
```
- Runs ONLY on the Go server before rendering
- Direct DB access, no API layer needed

### 3. Template (HTML) — The Markup
```html
<div class="user-card">
    <h1>Hello, {user.Name}!</h1>
    <button id="map-btn">Center map</button>
</div>
```

### 4. `<script>` — Client-side JavaScript
```html
<script lang="ts">
    document.getElementById("map-btn")?.addEventListener("click", () => {
        console.log("Centering map...");
    });
</script>
```
- V1: Vanilla JS (no TypeScript, no compiler)
- V2: TypeScript via esbuild integration

### 5. `<style>` — Scoped CSS
```html
<style>
    .user-card {
        padding: 1rem;
        background: #f4f4f4;
    }
</style>
```
- Classes are automatically hashed for scoping

---

## Template Logic (V1)

### {#if} / {#else if} / {#else}
```html
{#if user.IsAdmin}
    <a href="/admin">Dashboard</a>
{#else if user.IsLoggedIn}
    <a href="/profile">My Profile</a>
{#else}
    <a href="/login">Login</a>
{/if}
```
→ Compiles to `if condition { ... } else { ... }`

### {#switch} / {#case} / {#default}
```html
{#switch order.Status}
    {#case "pending"}
        <span class="badge yellow">Processing</span>
    {#case "shipped"}
        <span class="badge blue">En route</span>
    {#default}
        <span class="badge gray">Unknown</span>
{/switch}
```
→ Compiles to native Go `switch`

### {#each} / {#else}
```html
<ul>
    {#each users as user, index}
        <li>#{index + 1}: {user.Name}</li>
    {#else}
        <li>No users found.</li>
    {/each}
</ul>
```
→ {#else} on empty list — no `if len(users) == 0` needed

### {#await} (Async/Streams)
```html
{#await fetchUserData()}
    <div class="skeleton-loader">Loading...</div>
{#then user}
    <p>Welcome, {user.Name}!</p>
{#catch err}
    <p class="error">{err.Error()}</p>
{/await}
```
→ Uses Go channels or JS promises

### {#slot} / {#fill} (Component Slots)
```html
<!-- In Card.dreego -->
<div class="card">
    <div class="card-header">
        {#slot header}
            <h3>Default Title</h3>
        {/slot}
    </div>
    <div class="card-body">
        {#slot}
        {/slot}
    </div>
</div>
```

```html
<!-- Usage -->
<Card>
    {#fill header}
        <h3 class="gold">Premium User</h3>
    {/fill}
    <p>This is the content.</p>
</Card>
```

### {#let} (Helper Variables)
```html
{#let fullName = user.FirstName + " " + user.LastName}
{#let isGoldCustomer = user.OrdersCount > 50}
```

---

## File-based Routing with Chi

### Directory Structure
```
src/
└── routes/
     ├── get.dreego           →  /
     ├── about.dreego           →  /about
     └── users/
          └── [id].dreego        →  /users/:id
```

### Routing Engine
- `dreego generate` scans the `routes/` directory
- Generates a central `dreego_router.go`:
```go
// CODE GENERATED BY DREEGO. DO NOT EDIT.
package main

import "net/http"

func RegisterDreegoRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/", RenderIndexPage)
    mux.HandleFunc("/users/{id}", RenderUsersPage)
}
```
→ No runtime scanning, no disk IO, everything in the binary

---

## Dreego's Core Stack

| Area               | Choice                    | Reasoning                                    |
|--------------------|---------------------------|----------------------------------------------|
| HTTP Router        | go-chi/chi                | Fast, 100% Go stdlib compatible              |
| Template Engine    | Dreego Transpiler         | Centerpiece — must be custom                 |
| Interactivity      | HTMX + Alpine.js          | LiveView feel without SPA complexity         |
| SSE/Real-time      | Datastar (optional)       | Streams DOM updates via SSE                  |
| CSS                | Tailwind CLI (embedded)   | No custom CSS parser needed                  |
| Validation         | go-playground/validator   | Struct tag validation                        |
| Binary Packaging   | `embed` (Go Stdlib)       | Native Go, no extra tool                     |
| JS Bundling (V2)   | esbuild                   | Embeddable as Go library                     |

---

## When is Dreego Better than the Competition?

### Better than React / Vue / Svelte (SPAs)
- No JS build hell (no node_modules)
- 0 MB JS state synchronization (state only on server)
- Perfect SEO + FCP even on weak devices

### Better than Phoenix (Elixir)
- Single binary deployment (no BEAM VM, no Node.js)
- Compile-time and type safety (Go compiler prevents runtime errors)
- Fraction of RAM usage, starts in milliseconds

### Better than plain Go templates / Templ
- Real framework experience (Routing, Actions, Layouts, Dev Server)
- No JS glue code (HTMX/Datastar directly integrated)

---

## Addon/Plugin Ecosystem

### Architecture
A Goree addon is a Go package that fulfills the `dreego.Plugin` interface:

```go
type Plugin interface {
    Name() string
    RegisterRoutes(mux http.Handler)
    Middlewares() []func(http.Handler) http.Handler
    Assets() *embed.FS
}
```

### Usage in main.go
```go
app := dreego.New()
app.UsePlugin(auth.New("super-secret-key"))
app.Listen(":8080")
```

### Addon Ideas (full list)

| Addon               | Purpose                                                  |
|---------------------|----------------------------------------------------------|
| dreego-auth          | Sessions, OAuth, Passkeys, Login/Register                |
| dreego-map           | MapLibre/Leaflet Integration                             |
| dreego-admin         | Auto-generated admin dashboard                           |
| dreego-seo           | Meta tags, OpenGraph, Sitemap                            |
| dreego-db            | DB integration (SQLite, Turso, PostgreSQL)               |
| dreego-analytics     | Privacy-friendly, server-side                            |
| dreego-i18n          | Multi-language support                                   |
| dreego-mail          | Email delivery with .dreego templates                    |
| dreego-stripe        | Stripe Checkout & Webhooks                               |
| dreego-storage       | File uploads (S3, R2, local)                             |
| dreego-jobs          | Background tasks & cron jobs                             |
| dreego-ui            | Shadcn-like UI components                                |
| dreego-search        | Full-text search (Bleve/Meilisearch)                     |
| dreego-pdf           | PDF generation from .dreego templates                    |
| dreego-pwa           | Progressive Web App                                      |

### Benefits of the Go Addon System
- No dependency hell (go.mod resolves cleanly)
- Compile-Time Safety: Build breaks on errors
- Tree-Shaking: Unused code is compiled out

---

## V1 Blueprint (MVP)

1. **Transpiler:** Reads `.dreego` → generates `.go`
2. **3 Sections:** `<go>`, HTML Template, `<style>`
3. **Template Logic:** `{#if}` and `{#each}`
4. **Chi Router Wrapper:** Directory paths → HTTP routes
5. **Single Binary:** `//go:embed` for assets

## V2 Outlook

- TypeScript Support via esbuild
- SSG (Static Site Generation)
- Wails Integration (Desktop Apps)
- `{#await}`, `{#slot}`, `{#switch}`
- Auto-Form Binding
- Inline API Endpoints

---

## Error Handling in Dreego

**No `<catch>` tag!** Errors are handled in Go idioms:

```html
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

<div class="profile">
    {#if hasError}
        <p class="error">User could not be loaded.</p>
    {#else}
        <h1>Hello, {user.Name}!</h1>
    {/if}
</div>
```

---

## Tailwind Class Merging for dreego-ui

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

Usage:
```html
<Button class="bg-red-600 mt-4" disabled>Delete</Button>
```

---

## A/B Testing via Go Logic

```html
<go>
    userIP := c.Request.RemoteAddr
    showPriceVariantA := isOddIP(userIP)
</go>

<div>
    {#if showPriceVariantA}
        <div class="price-tag">Special price: 19 EUR</div>
    {#else}
        <div class="price-tag">Regular price: 29 EUR</div>
    {/if}
</div>
```

---

## Dreego vs Svelte vs React — Summary

- **Feels like Svelte:** `.dreego` files, HTML-first syntax
- **Structured like React:** Components, Props, Children-Slots
- **More robust than both:** No Node.js, Single Binary, Go type safety

---

## Sources from the Chat

1. Cross-Platform Go: Runtime Checks vs. Build Tags
2. Customizing Go Binaries with Build Tags - DigitalOcean
3. Gomacro: code generation made easy and fun
4. How to Use go:embed in Go - JetBrains
5. Chi vs Encore.go - Which Go Backend Framework to Choose in 2026
6. Datastar: Web Framework for the Future?
7. irgo - Native Apps with Go + Datastar
