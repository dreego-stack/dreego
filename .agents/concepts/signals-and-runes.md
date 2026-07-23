# Signals & Svelte Runes — Konzept für Dreego

## Was sind Signals?

Signals sind ein reaktives State-Primitive. Vereinfacht:

```js
const count = signal(0)
// count.value = 5 → alles was count liest, updated sich automatisch
```

Das Konzept ist framework-übergreifend:
- **Solid.js** hat Signals erfunden (fein-granulare Reaktivität ohne Virtual DOM)
- **Svelte 5 Runes** sind Signals als Compiler-Feature (`$state`, `$derived`, `$effect`)
- **Angular** hat Signals eingeführt
- **Preact** hat Signals
- **Vue** hat `ref()` / `reactive()` — funktional ähnlich
- **Signals sind #3 der meistgewünschten JS-Features** (State of JS 2025)

### Warum sind Signals so beliebt?

| Ohne Signals (React)           | Mit Signals (Svelte 5)           |
|--------------------------------|----------------------------------|
| `useState` + `useEffect`       | `let count = $state(0)`         |
| State-Änderung → ganzer Komponentenbaum re-rendert | State-Änderung → nur betroffene DOM-Knoten updaten |
| Virtual DOM Diffing            | Direktes DOM-Updating            |
| Boilerplate für Memoization    | Automatisch optimiert            |

## Svelte 5 Runes

```svelte
<script>
    let count = $state(0)         // Reaktive Variable
    let doubled = $derived(count * 2)  // Automatisch neu berechnet
    $effect(() => {
        console.log(`Count: ${count}`)  // Läuft bei jeder Änderung
    })
</script>

<button onclick={() => count++}>
    Count: {count} (×2 = {doubled})
</button>
```

**Runes sind Compiler-Direktiven.** Der Svelte-Compiler analysiert, welche Variablen reaktiv sind, und generiert minimalen Update-Code. Kein Runtime-Overhead.

## Dreego und Signals: Zwei Ebenen

Dreego ist SSR-First. Signals funktionieren bei uns auf zwei Ebenen:

### Ebene 1: Server-Seite (im `<go>`-Block)

Auf dem Server gibt es keine client-seitige Reaktivität. Der `<go>`-Block läuft einmal pro Request und erzeugt HTML. Das ist performant und einfach.

ABER: Das Reaktivitäts-Konzept von Svelte Runes entspricht funktional unserem Datenfluss:

```html
<go>
    count := 0                              // $state(0)
    doubled := count * 2                    // $derived(count * 2)
</go>

<div>
    Count: {count} (×2 = {doubled})         <!-- wie im Svelte-Template -->
    <button hx-post="/increment" ...>+1</button>
</div>
```

Bei einem Klick:
1. HTMX sendet POST an `/increment`
2. Go-Handler erhöht `count`, rendert HTML-Fragment neu
3. HTMX tauscht das Fragment im DOM aus

Das ist im Prinzip dasselbe wie ein Signal-Update — nur dass der State auf dem Server lebt und das Update als HTML-Fragment kommt.

### Ebene 2: Client-Seite (Alpine.js / Datastar)

Für lokale Interaktivität OHNE Server-Roundtrip nutzen wir Alpine.js:

```html
<div x-data="{ count: 0 }">
    <button @click="count++">
        Count: <span x-text="count"></span>  <!-- Reaktiv! -->
    </button>
</div>
```

`x-data` ist ein client-seitiges Signal. Alpine.js updated nur das betroffene `x-text`-Element — fine-grained Reactivity.

**Datastar** geht noch weiter — es bringt echte Signals mit SSE-Backend:

```html
<div data-signals="{count: 0}">
    <button data-on-click="$count++">
        Count: <span data-text="$count"></span>
    </button>
</div>
```

Der Unterschied: Alpine.js macht alles im Browser. Datastar kann State auch zum Server streamen und zurückbekommen (bidirektional).

## Was Dreego von Signals lernen kann

### Jetzt umsetzbar (V1)

| Signal-Konzept      | Dreego-Entsprechung                                  |
|---------------------|-----------------------------------------------------|
| `$state(x)`         | `<go>`-Block Variablen + HTML-Fragment-Updates via HTMX |
| `$derived(expr)`    | `{#let name = expr}` — im `<go>`-Block berechnen    |
| `$effect(fn)`       | HTMX-Events + Alpine `@click`, `@change`            |
| Reaktive DOM-Updates| HTMX partial swaps / Alpine `x-text`, `x-show`      |

### Was Dreego NICHT braucht

- **Virtual DOM** — wir haben keinen. SSR + Partial HTML Swaps sind schneller.
- **Client-seitigen Compiler** — Alpine.js ist 15 KB und macht alles nativ im Browser.
- **Reaktiven Template-Compiler** — das ist Sveltes Job. Dreego macht SSR.

### V2-Potenzial: Dreego Signals als First-Class Concept

```html
<go>
    count := 0
</go>

<!-- Dreego könnte ein "reactive fragment" Konzept einführen -->
<button hx-post="/inc" hx-target="#counter" hx-swap="outerHTML">
    <span id="counter">{count}</span>
</button>
```

Oder via Datastar-Integration:

```html
<go>
    count := dreego.Signal(0)
</go>

<button data-on-click="$count++">
    <span data-text="$count"></span>
</button>
```

## Zusammenfassung

Signals sind ein UX-Konzept (automatische DOM-Updates bei State-Änderung), kein Technologie-Konzept. Dreego erreicht dasselbe durch:

1. **SSR + HTMX** für Server-getriebene Updates (wie Phoenix LiveView)
2. **Alpine.js** für lokale, client-seitige Interaktivität
3. **Datastar** für SSE-basierte, bidirektionale Signals (optional)

Svelte-Runes sind ein Spezialfall von Signals, die durch den Compiler optimiert werden. Dreego braucht diesen Compiler nicht, weil es den State auf dem Server hält und Updates als HTML schickt — ein fundamental anderer, aber ebenso valider Ansatz.
