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

## Attribute Props

Props can be passed inside any HTML attribute, both in the component call and in the component body. A `{prop}` (or `{expr}`) placeholder inside a quoted attribute value is resolved to the Go expression.

**In the call:**

```html
<div><@Card link="mailto:{email}"/></div>
```

**In the component body:**

```html
Component Link (url string)
<a href="{url}">{#slot}</a>
```

```html
<div><@Link url="https://dreego.dev">Home</@Link></div>
```

The attribute expression is HTML-escaped before emission, so `href="{url}"` with untrusted `url` input is safe.

## Rules

1. **`Component` line** — Always line 1 of the file.
2. **`<@Name>`** — Component call. `@` prefix distinguishes from HTML tags.
3. **File-based Discovery** — `dreego/components/Card.dreego` → `<@Card>`.
4. **Scoped Styles** — `data-scope` per component. No leak to parent.
5. **Self-closing** — `<@Icon name="star"/>` when no body.
6. **Slots** — `{#slot}` in component template = child content.
7. **Attribute Props** — `{prop}` inside a quoted attribute value is resolved and escaped.

## Scoped CSS

Every component's `<style>` block is scoped via a `data-scope` attribute (a 12-char hash of the component source). Selectors are rewritten to `[data-scope=hash] <selector>` so styles never leak to parent or sibling components.

- **Declarations preserved**: the body between `{` and `}` is copied verbatim, so values like `radial-gradient(circle, #ccfbf1 1px, transparent 1px)`, `calc(100% - 20px)` and `rgb(1, 2, 3)` survive intact.
- **Nested selectors**: comma-separated selectors are each scoped on their own line; pseudo-selectors like `:hover` scope correctly.
- **`@media`**: the `@media` header stays unscoped, its inner selectors are scoped recursively.
- **`@keyframes`**: header and full body are copied verbatim (unscoped) so animation steps (`from`, `to`, percent) are preserved; only referencing selectors are scoped.

```html
Component Spinner ()
<div class="spinner"></div>
<style>
@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }
.spinner { animation: spin 1s linear infinite; }
</style>
```

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
