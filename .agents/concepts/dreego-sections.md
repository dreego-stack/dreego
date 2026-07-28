
---
type: Concept
title: "The 5 Sections of a .dreego File"
description: "Structure and behavior of the 5 sections: head, go, Template, script, style"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# The 5 Sections of a .dreego File

## Overview

Every `.dreego` file consists of up to 5 sections. The order is:

1. `<head>` — Component-specific assets
2. `<go>` — Server-side Go code
3. Template (implicit, the HTML part) — The markup
4. `<script>` — Client-side JavaScript
5. `<style>` — Scoped CSS

## Section 1: `<head>`

**Purpose:** Declare assets that only need to be loaded for this component.

```html
<head>
    <script src="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.js"></script>
    <link href="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.css" rel="stylesheet" />
    <meta name="description" content="Interactive map" />
</head>
```

**Behavior:**
- Collected by the transpiler
- Dynamically injected into the HTML head when the component is rendered
- Global scripts are only loaded when the component is actually rendered

## Section 2: `<go>`

**Purpose:** Server-side logic — database queries, request processing, business logic.

```html
<go>
    userID := c.Param("id")
    user, err := db.GetUser(userID)
    hasError := err != nil

    type UserLocation struct {
        Lat float64 `json:"lat"`
        Lng float64 `json:"lng"`
    }
    loc := UserLocation{Lat: 52.52, Lng: 13.40}
</go>
```

**Behavior:**
- Compiled 1:1 into a Go function
- Has access to `*http.Request`, `http.ResponseWriter`, DB pool, session, etc.
- Runs EXCLUSIVELY on the server
- Never visible to the client

## Section 3: Template (HTML)

**Purpose:** The markup. Data from `<go>` is rendered here.

```html
<div class="user-card">
    {#if hasError}
        <p class="error">User could not be loaded.</p>
    {#else}
        <h1>Hello, {user.Name}!</h1>
        <p>Age: {user.Age}</p>
    {/if}
</div>
```

**Supported Template Logic:**
- `{#if}`, `{#else if}`, `{#else}`, `{/if}`
- `{#each items as item, index}`, `{#else}`, `{/each}`
- `{#switch expr}`, `{#case val}`, `{#default}`, `{/switch}`
- `{#let name = expr}`
- `{variable}` — escaped HTML output
- `{variable|raw}` — unescaped (only when safe!)

## Section 4: `<script>`

**Purpose:** Client-side JavaScript. V1: Vanilla JS. V2: TypeScript.

```html
<script>
    document.getElementById("map-btn")?.addEventListener("click", () => {
        console.log("Centering map...");
    });
</script>
```

**Behavior:**
- V1: Extracted 1:1 and embedded as `<script>` tag in the HTML
- V2: TypeScript compiled via esbuild
- Runs EXCLUSIVELY in the browser
- Has no access to Go variables (explicit separation!)

## Section 5: `<style>`

**Purpose:** Component-specific CSS. Automatically scoped.

```html
<style>
    .user-card {
        padding: 1rem;
        border-radius: 8px;
        background: #f4f4f4;
    }
    .user-card h1 {
        font-size: 1.5rem;
    }
</style>
```

**Behavior:**
- Classes automatically receive a hash (`.user-card-a8f3d`)
- Hashed classes are used in the template
- No global CSS leaking

## Transpiler Output

The transpiler processes a `.dreego` file into 3 outputs:

1. **Go Code:** `<go>` + Template → `_dreego.go` (server rendering)
2. **CSS:** `<style>` → collected into final CSS file
3. **JS:** `<script>` + `<head>` → embedded in HTML or as separate file
