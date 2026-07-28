
---
type: Concept
title: "Template Logic in Dreego"
description: "Template syntax and its codegen to native Go code (as of v0.0.10)"
tags: [v0.0.10]
timestamp: 2026-07-28T00:00:00Z
---

# Template Logic in Dreego

## Design Philosophy

- **No real Go in the template** — complex logic belongs in `<go>` block
- **All template blocks compile to native Go** — zero runtime cost
- **Auto-escaping**: all `{var}` expressions are escaped via `html.EscapeString`

## Template Blocks (as of v0.0.10)

| Block | Status | Go Equivalent |
|---|---|---|
| `{#if cond}` / `{#else}` / `{/if}` | ✅ | `if cond { } else { }` |
| `{#each items as item}` / `{#each else}` / `{/each}` | ✅ | `for i, item := range items { }` + `len == 0` check |
| `{#slot}` | ✅ | `c.Get("slot")` |
| `{#slot name}...{/slot}` | ✅ | `c.Get("slot_name")` |
| `{#verbatim}...{/verbatim}` | ✅ | Raw `b.WriteString()` |
| `$loop.Index / .First / .Last / .Even / .Odd` | ✅ | Codegen-generated `core.EachLoop` struct |
| `{var\|raw}` / `{var\|upper}` | ✅ | Filter chain in codegen |
| `{#switch}` / `{#case}` | ❌ V2 | — |
| `{#await}` | ❌ V2 | — |
| `{#let}` | ❌ V2 | — |
| `{#each else if}` | ❌ | Not yet implemented |

## {#if} / {#else}

```
{#if show}
    <p>visible</p>
{#else}
    <p>hidden</p>
{/if}
```

**Rules:**
- Supports arbitrary Go conditions (variables from `<go>`)
- `{#else}` optional
- `{#else if}` not yet implemented

## {#each} / {#each else} / $loop

```
{#each users as user}
    <li>{$loop.Index}: {user.Name}</li>
{#each else}
    <li>No data</li>
{/each}
```

**Rules:**
- Iterates over slice/array (variables from `<go>`)
- `{#each else}` renders on empty slice
- `$loop.Index` (0-based), `.First`, `.Last`, `.Even`, `.Odd`
- Codegen: `var loop := core.EachLoop{Index: i, ...}` with string replacement `$loop.` → `loop.`

## {#verbatim}

```
{#verbatim}
    <script>var x = {a: 1};</script>
{/verbatim}
```

**Rules:**
- Everything between `{#verbatim}` and `{/verbatim}` is output 1:1
- No parsing, no escaping — perfect for JS templates
- Lexer scans as a single `TokenVerbatim` with raw content

## Filters

```
<p>{html|raw}</p>
<div>{name|upper}</div>
```

**Rules:**
- `|raw` — no HTML escaping (for trusted HTML)
- `|upper` — `strings.ToUpper()` for text
- `parseExpression()` splits at `|` in the parser
- Codegen builds filter chain: `strings.ToUpper(fmt.Sprintf("%v", name))`

## {#slot} / Named Slots

**Default Slot (no name, no `{/slot}`):**

```
<div>{#slot}</div>
```

**Named Slots (with `{/slot}` closing):**

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
- Components with children: `cb.WriteString` / `sb.WriteString` buffer redirection

## Error Handling

No special error tag. Errors via `<go>` block and `{#if}`:

```
<go>
    user, err := db.GetUser(id)
    hasError := err != nil
</go>

{#if hasError}
    <p>Error loading.</p>
{#else}
    <h1>{user.Name}</h1>
{/if}
```
