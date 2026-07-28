---
type: Concept
title: "Template-Logik in Dreego"
description: "Template-Syntax und deren Codegen zu nativem Go-Code (Stand v0.0.10)"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Template-Logik in Dreego

## Design-Philosophie

- **Kein echtes Go im Template** — komplexe Logik gehört in `<go>`-Block
- **Alle Template-Blöcke kompilieren zu nativem Go** — null Laufzeit-Verlust
- **Auto-Escaping**: alle `{var}` Ausdrücke werden via `html.EscapeString` escaped

## Template-Blöcke (Stand v0.0.10)

| Block | Status | Go-Entsprechung |
|---|---|---|
| `{#if cond}` / `{#else}` / `{/if}` | ✅ | `if cond { } else { }` |
| `{#each items as item}` / `{#each else}` / `{/each}` | ✅ | `for i, item := range items { }` + `len == 0` check |
| `{#slot}` | ✅ | `c.Get("slot")` |
| `{#slot name}...{/slot}` | ✅ | `c.Get("slot_name")` |
| `{#verbatim}...{/verbatim}` | ✅ | Raw `b.WriteString()` |
| `$loop.Index / .First / .Last / .Even / .Odd` | ✅ | Codegen-generiertes `core.EachLoop` struct |
| `{var\|raw}` / `{var\|upper}` | ✅ | Filter-Chain im Codegen |
| `{#switch}` / `{#case}` | ❌ V2 | — |
| `{#await}` | ❌ V2 | — |
| `{#let}` | ❌ V2 | — |
| `{#each else if}` | ❌ | Noch nicht implementiert |

## {#if} / {#else}

```
{#if show}
    <p>sichtbar</p>
{#else}
    <p>versteckt</p>
{/if}
```

**Regeln:**
- Unterstützt beliebige Go-Bedingungen (Variablen aus `<go>`)
- `{#else}` optional
- `{#else if}` noch nicht implementiert

## {#each} / {#each else} / $loop

```
{#each users as user}
    <li>{$loop.Index}: {user.Name}</li>
{#each else}
    <li>Keine Daten</li>
{/each}
```

**Regeln:**
- Iteriert über Slice/Array (Variablen aus `<go>`)
- `{#each else}` rendert bei leerer Slice
- `$loop.Index` (0-basiert), `.First`, `.Last`, `.Even`, `.Odd`
- Codegen: `var loop := core.EachLoop{Index: i, ...}` mit String-Replacement `$loop.` → `loop.`

## {#verbatim}

```
{#verbatim}
    <script>var x = {a: 1};</script>
{/verbatim}
```

**Regeln:**
- Alles zwischen `{#verbatim}` und `{/verbatim}` wird 1:1 ausgegeben
- Kein Parsing, kein Escaping — perfekt für JS-Templates
- Lexer scanned als einzelnes `TokenVerbatim` mit Raw-Content

## Filter

```
<p>{html|raw}</p>
<div>{name|upper}</div>
```

**Regeln:**
- `|raw` — kein HTML-Escaping (für vertrauenswürdiges HTML)
- `|upper` — `strings.ToUpper()` für Text
- `parseExpression()` splittet am `|` im Parser
- Codegen baut Filter-Chain: `strings.ToUpper(fmt.Sprintf("%v", name))`

## {#slot} / Named Slots

**Default-Slot (kein Name, kein `{/slot}`):**

```
<div>{#slot}</div>
```

**Named Slots (mit `{/slot}` closing):**

Component:
```
Component Card (title string)
<div>
  {#slot header}{/slot}
  <h2>{title}</h2>
  {#slot}
</div>
```

Route:
```
<@Card title="Hi">
  {#slot header}<nav>menu</nav>{/slot}
  <p>body content</p>
</@Card>
```

**Codegen:**
- Route: `c.Set("slot_header", sb.String())`, `c.Set("slot", cb.String())`
- Component: `b.WriteString(ctx.Get("slot_header"))`, `b.WriteString(ctx.Get("slot"))`
- Components mit Kindern: `cb.WriteString` / `sb.WriteString` Puffer-Umleitung

## Fehlerbehandlung

Kein spezielles Error-Tag. Fehler via `<go>`-Block und `{#if}`:

```
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

{#if hasError}
    <p>Fehler beim Laden.</p>
{#else}
    <h1>{user.Name}</h1>
{/if}
```
