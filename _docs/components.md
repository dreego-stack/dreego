# Components

Component = reusable `.dreego` file with props and its own scope.

## Declaration

```html
Component Name (prop Type, prop Type = default)
```

The first line of a component file **must** be the `Component` declaration.

| Element | Description |
|---------|-------------|
| `Name` | Component name (PascalCase). Becomes `<@Name>`. |
| `(prop Type)` | Props with Go type. `= default` optional. |
| `{#slot}` | Default slot (always available). |

**Example:**

```
Component Card (title string)

<div>
    <article class="card">
        <h2>{title}</h2>
        <div>{#slot}</div>
    </article>
</div>
<style>
.card { border: 1px solid var(--border); padding: 1rem; }
</style>
```

## Usage

Components are automatically discovered from `dreego/components/`. No `import` needed.

**Self-closing:**

```html
<div><@Card title="Hello"/></div>
```

**With children (default slot):**

```html
<div>
    <@Card title="Welcome">
        <p>Slot content goes here</p>
    </@Card>
</div>
```

## Rules

1. **`Component` line** — Always line 1 of the file.
2. **`<@Name>`** — Component call. `@` prefix distinguishes from HTML tags.
3. **File-based Discovery** — `dreego/components/Card.dreego` → `<@Card>`.
4. **Scoped Styles** — `data-scope` per component. No leak to parent.
5. **Self-closing** — `<@Icon name="star"/>` when no body.
6. **Slots** — `{#slot}` in component template = child content.

## Named Slots (v0.0.8)

**Component:**
```
Component Card (title string)

<div>
    <article>
        {#slot header}{/slot}
        <h2>{title}</h2>
        {#slot}
    </article>
</div>
```

**Route:**
```html
<@Card title="Hi">
    {#slot header}<strong>HEADER</strong>{/slot}
    <p>Default content here</p>
</@Card>
```

- `{#slot header}{/slot}` — Placeholder in component (empty body)
- `{#slot header}content{/slot}` — Definition in the route (with body)
- `{#slot}` — Default slot (no `{/slot}` needed)
- Multiple named slots per component possible

## Generated Go Code

```go
func Card(title string) core.Component {
    return core.ComponentFunc(func(ctx *core.SSRContext) (string, error) {
        var b strings.Builder
        b.WriteString("<div data-scope=\"abc123\">")
        b.WriteString(`<article class="card">`)
        b.WriteString(`<h2>`)
        b.WriteString(html.EscapeString(fmt.Sprintf("%v", title)))
        b.WriteString(`</h2>`)
        b.WriteString(ctx.Get("slot"))
        b.WriteString(`</article>`)
        b.WriteString("</div>")
        return b.String(), nil
    })
}
```

Call `<@Card title="x"/>` → `Card("x").Render(c)`.
