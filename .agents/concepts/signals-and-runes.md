
---
type: Concept
title: "Signals & Svelte Runes — Concept for Dreego"
description: "Reactive state primitives and their equivalent in Dreego"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---
# Signals & Svelte Runes — Concept for Dreego

## What are Signals?

Signals are a reactive state primitive. Simplified:

```js
const count = signal(0)
// count.value = 5 → everything that reads count updates automatically
```

The concept is framework-agnostic:
- **Solid.js** invented signals (fine-grained reactivity without Virtual DOM)
- **Svelte 5 Runes** are signals as a compiler feature (`$state`, `$derived`, `$effect`)
- **Angular** introduced signals
- **Preact** has signals
- **Vue** has `ref()` / `reactive()` — functionally similar
- **Signals are #3 of the most wanted JS features** (State of JS 2025)

### Why are signals so popular?

| Without Signals (React)        | With Signals (Svelte 5)          |
|--------------------------------|----------------------------------|
| `useState` + `useEffect`       | `let count = $state(0)`         |
| State change → entire component tree re-renders | State change → only affected DOM nodes update |
| Virtual DOM Diffing            | Direct DOM updating              |
| Boilerplate for memoization    | Automatically optimized          |

## Svelte 5 Runes

```svelte
<script>
    let count = $state(0)              // Reactive variable
    let doubled = $derived(count * 2)  // Automatically recalculated
    $effect(() => {
        console.log(`Count: ${count}`) // Runs on every change
    })
</script>

<button onclick={() => count++}>
    Count: {count} (×2 = {doubled})
</button>
```

**Runes are compiler directives.** The Svelte compiler analyzes which variables are reactive and generates minimal update code. No runtime overhead.

## Dreego and Signals: Two Levels

Dreego is SSR-First. Signals work at two levels for us:

### Level 1: Server-Side (in the `<go>` block)

On the server, there is no client-side reactivity. The `<go>` block runs once per request and produces HTML. This is performant and simple.

BUT: The reactivity concept of Svelte Runes functionally corresponds to our data flow:

```html
<go>
    count := 0                              // $state(0)
    doubled := count * 2                    // $derived(count * 2)
</go>

<div>
    Count: {count} (×2 = {doubled})         <!-- like in Svelte template -->
    <button hx-post="/increment" ...>+1</button>
</div>
```

On click:
1. HTMX sends POST to `/increment`
2. Go handler increments `count`, re-renders HTML fragment
3. HTMX swaps the fragment in the DOM

This is essentially the same as a signal update — except the state lives on the server and the update comes as an HTML fragment.

### Level 2: Client-Side (Alpine.js / Datastar)

For local interactivity WITHOUT a server round-trip, we use Alpine.js:

```html
<div x-data="{ count: 0 }">
    <button @click="count++">
        Count: <span x-text="count"></span>  <!-- Reactive! -->
    </button>
</div>
```

`x-data` is a client-side signal. Alpine.js only updates the affected `x-text` element — fine-grained reactivity.

**Datastar** goes even further — it brings real signals with an SSE backend:

```html
<div data-signals="{count: 0}">
    <button data-on-click="$count++">
        Count: <span data-text="$count"></span>
    </button>
</div>
```

The difference: Alpine.js does everything in the browser. Datastar can also stream state to the server and receive it back (bidirectional).

## What Dreego Can Learn from Signals

### Implementable Now (V1)

| Signal Concept       | Dreego Equivalent                                     |
|----------------------|------------------------------------------------------|
| `$state(x)`          | `<go>` block variables + HTML fragment updates via HTMX |
| `$derived(expr)`     | `{#let name = expr}` — calculate in `<go>` block     |
| `$effect(fn)`        | HTMX events + Alpine `@click`, `@change`             |
| Reactive DOM Updates | HTMX partial swaps / Alpine `x-text`, `x-show`       |

### What Dreego Does NOT Need

- **Virtual DOM** — we don't have one. SSR + Partial HTML Swaps are faster.
- **Client-side Compiler** — Alpine.js is 15 KB and does everything natively in the browser.
- **Reactive Template Compiler** — that's Svelte's job. Dreego does SSR.

### V2 Potential: Dreego Signals as a First-Class Concept

```html
<go>
    count := 0
</go>

<!-- Dreego could introduce a "reactive fragment" concept -->
<button hx-post="/inc" hx-target="#counter" hx-swap="outerHTML">
    <span id="counter">{count}</span>
</button>
```

Or via Datastar integration:

```html
<go>
    count := dreego.Signal(0)
</go>

<button data-on-click="$count++">
    <span data-text="$count"></span>
</button>
```

## Summary

Signals are a UX concept (automatic DOM updates on state change), not a technology concept. Dreego achieves the same via:

1. **SSR + HTMX** for server-driven updates (like Phoenix LiveView)
2. **Alpine.js** for local, client-side interactivity
3. **Datastar** for SSE-based, bidirectional signals (optional)

Svelte Runes are a special case of signals optimized by the compiler. Dreego doesn't need this compiler because it keeps state on the server and sends updates as HTML — a fundamentally different but equally valid approach.
