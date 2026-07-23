# Die 5 Sektionen einer .dreego-Datei

## Übersicht

Jede `.dreego`-Datei besteht aus bis zu 5 Sektionen. Die Reihenfolge ist:

1. `<head>` — Komponenten-spezifische Assets
2. `<go>` — Server-seitiger Go-Code
3. Template (implizit, der HTML-Teil) — Das Markup
4. `<script>` — Client-seitiges JavaScript
5. `<style>` — Scoped CSS

## Sektion 1: `<head>`

**Zweck:** Assets deklarieren, die nur für diese Komponente geladen werden müssen.

```html
<head>
    <script src="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.js"></script>
    <link href="https://api.mapbox.com/mapbox-gl-js/v2.14.1/mapbox-gl.css" rel="stylesheet" />
    <meta name="description" content="Interaktive Karte" />
</head>
```

**Verhalten:**
- Wird vom Transpiler gesammelt
- Beim Rendern der Komponente dynamisch in den HTML-Head injiziert
- Globale Skripte werden nur geladen, wenn die Komponente tatsächlich gerendert wird

## Sektion 2: `<go>`

**Zweck:** Server-seitige Logik — Datenbank-Abfragen, Request-Verarbeitung, Business-Logik.

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

**Verhalten:**
- Wird 1:1 in eine Go-Funktion kompiliert
- Hat Zugriff auf `*http.Request`, `http.ResponseWriter`, DB-Pool, Session, etc.
- Läuft AUSSCHLIESSLICH auf dem Server
- Nie im Client sichtbar

## Sektion 3: Template (HTML)

**Zweck:** Das Markup. Hier werden Daten aus `<go>` gerendert.

```html
<div class="user-card">
    {#if hasError}
        <p class="error">Benutzer konnte nicht geladen werden.</p>
    {#else}
        <h1>Hallo, {user.Name}!</h1>
        <p>Alter: {user.Age}</p>
    {/if}
</div>
```

**Unterstützte Template-Logik:**
- `{#if}`, `{#else if}`, `{#else}`, `{/if}`
- `{#each items as item, index}`, `{#else}`, `{/each}`
- `{#switch expr}`, `{#case val}`, `{#default}`, `{/switch}`
- `{#let name = expr}`
- `{variable}` — escaped HTML-Output
- `{variable|raw}` — unescaped (nur wenn sicher!)

## Sektion 4: `<script>`

**Zweck:** Client-seitiges JavaScript. V1: Vanilla JS. V2: TypeScript.

```html
<script>
    document.getElementById("map-btn")?.addEventListener("click", () => {
        console.log("Karte wird zentriert...");
    });
</script>
```

**Verhalten:**
- V1: Wird 1:1 extrahiert und als `<script>`-Tag ins HTML eingebettet
- V2: TypeScript via esbuild kompiliert
- Läuft AUSSCHLIESSLICH im Browser
- Hat keinen Zugriff auf Go-Variablen (explizite Trennung!)

## Sektion 5: `<style>`

**Zweck:** Komponenten-spezifisches CSS. Wird automatisch gescoped.

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

**Verhalten:**
- Klassen erhalten automatisch einen Hash (`.user-card-a8f3d`)
- Im Template werden die gehashten Klassen verwendet
- Kein globales CSS-Leaking

## Transpiler-Output

Der Transpiler verarbeitet eine `.dreego`-Datei in 3 Ausgaben:

1. **Go-Code:** `<go>` + Template → `_dreego.go` (Server-Rendering)
2. **CSS:** `<style>` → gesammelt in finale CSS-Datei
3. **JS:** `<script>` + `<head>` → in HTML eingebettet oder als separate Datei
