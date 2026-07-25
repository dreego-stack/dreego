
---
type: Concept
title: "Template-Logik in Dreego"
description: "Template-Blöcke und deren Kompilierung zu nativem Go-Code"
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Template-Logik in Dreego

## Design-Philosophie

- **Kein echtes Go im Template** — das wäre unleserlich (Razor/JSX-Falle)
- **Komplexe Logik** → `<go>`-Block
- **Visuelle Struktur-Logik** → Template-Syntax
- **Alle Template-Blöcke kompilieren zu nativem Go** — null Laufzeit-Verlust

## Template-Blöcke (Übersicht)

| Block                    | Einsatzzweck              | Go-Entsprechung                  |
|--------------------------|---------------------------|----------------------------------|
| `{#if}` / `{#else}`     | Bedingte Anzeige          | `if ... { } else { }`           |
| `{#switch}` / `{#case}` | Mehrfach-Unterscheidung   | `switch val { case ...: }`      |
| `{#each}` / `{#else}`   | Slices/Arrays iterieren   | `for i, item := range slice`    |
| `{#await}` (V2)         | Async Data / SSE Streams  | Go Channels / JS Promises        |
| `{#slot}` / `{#fill}`   | Layout-Verschachtelung    | Function Callbacks               |
| `{#let}`                | Hilfsvariablen im Markup  | `var := expr`                    |

## {#if} / {#else if} / {#else}

```html
{#if user.IsAdmin}
    <a href="/admin">Dashboard</a>
{#else if user.IsLoggedIn}
    <a href="/profile">Mein Profil</a>
{#else}
    <a href="/login">Anmelden</a>
{/if}
```

**Regeln:**
- Unterstützt `&&`, `||`, `!`, `==`, `!=`, `<`, `>`, `<=`, `>=`
- Variablen müssen im `<go>`-Block deklariert sein
- Kein beliebiger Go-Code (keine Funktionsaufrufe, kein `range`)

## {#switch} / {#case} / {#default}

```html
{#switch order.Status}
    {#case "pending"}
        <span class="badge yellow">In Bearbeitung</span>
    {#case "shipped"}
        <span class="badge blue">Unterwegs</span>
    {#case "delivered"}
        <span class="badge green">Zugestellt</span>
    {#default}
        <span class="badge gray">Unbekannt</span>
{/switch}
```

**Regeln:**
- `{#case}` unterstützt String, Int, Bool-Literale
- Kein Fallthrough (wie in Go)
- `{#default}` ist optional

## {#each} / {#else}

```html
<ul>
    {#each users as user, index}
        <li>#{index + 1}: {user.Name} ({user.Email})</li>
    {#else}
        <li>Keine Benutzer gefunden.</li>
    {/each}
</ul>
```

**Regeln:**
- `{#else}` wird angezeigt, wenn die Slice leer oder nil ist
- `index` ist optional (`{#each users as user}`)
- Der Index startet bei 0

## {#let}

```html
{#let fullName = user.FirstName + " " + user.LastName}
{#let isGoldCustomer = user.OrdersCount > 50}

<div class="profile">
    <h2>{fullName}</h2>
    {#if isGoldCustomer}
        <span class="badge">VIP Kunde</span>
    {/if}
</div>
```

**Regeln:**
- Nur String-Konkatenation und einfache Vergleiche
- Für komplexe Berechnungen → `<go>`-Block verwenden

## {#slot} / {#fill} (V2)

```html
<!-- Card.dreego -->
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
    <p>Das ist der Inhalt der Karte.</p>
</Card>
```

## {#await} (V2)

```html
{#await fetchUserData()}
    <div class="skeleton-loader">Daten werden geladen...</div>
{#then user}
    <p>Willkommen zurück, {user.Name}!</p>
{#catch err}
    <p class="error">Fehler: {err.Error()}</p>
{/await}
```

## Tailwind Class Merging

```html
<!-- Button.dreego -->
<go>
    finalClasses := dreego.MergeClasses(
        "bg-blue-500 text-white font-bold py-2 px-4 rounded",
        props.Class,
    )
</go>

<button class="{finalClasses}" {attrs}>
    {#slot}Button Text{/slot}
</button>
```

`MergeClasses()` überschreibt konfliktierende Tailwind-Klassen intelligent (z.B. `bg-red-600` überschreibt `bg-blue-500`).

## Fehlerbehandlung

Kein spezielles Error-Tag. Fehler werden im `<go>`-Block behandelt:

```html
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

{#if hasError}
    <p class="error">Benutzer konnte nicht geladen werden.</p>
{#else}
    <h1>Hallo, {user.Name}!</h1>
{/if}
```
