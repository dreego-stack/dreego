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
        <h2>{{ title }}</h2>
        <div>{#slot}</div>
    </article>
</div>
<style>
.card { border: 1px solid var(--border); padding: 1rem; }
</style>
```

## Usage

Component declarations and imports are header directives. They are the only
content allowed outside the five root sections: `<go>`, `<head>`, `<div>`,
`<style>`, and `<script>`. Free text, HTML, and component calls at the root are
generation errors.

```dreego
import Card "components/Card.dreego"

<div><@Card title="Hello" /></div>
```

**Self-closing:**

```html
<div><@Card title="Hello"/></div>
```

A self-closing call renders empty default and named slots. Content after the
tag belongs to the parent and remains a normal sibling.

**With children (default slot):**

```html
<div>
    <@Card title="Welcome">
        <p>Slot content goes here</p>
    </@Card>
</div>
```

## Self-closing Calls and Slot Fallback

`<@Card/>` is allowed and renders an empty default slot.

The following paragraph is a sibling of the component:

```html
<div><@Card/><p>Sibling content</p></div>
```

Use `<@Card>...</@Card>` when content should populate the component's slots.
`<@Card></@Card>` is valid and also renders with empty slots.

## Attribute Props

Dynamic text and HTML attribute values use `{{ expression }}`. Component props
use `{expression}` without quotes so the generated Go value keeps its type.
`prop={expr}` passes the Go expression `expr` directly to the generated
component call.

**In the call:**

```html
<div><@Card count={count} /></div>
```

**In the component body:**

```html
Component Link (url string)
<div><a href="{{ url }}">{#slot}</a></div>
```

```html
<div><@Link url="https://dreego.dev">Home</@Link></div>
```

Simple literal expressions (`"..."`, integer literals) are type-checked
against the declared prop type at `dreego generate` time. Non-literal expressions
are accepted unchecked because the transpiler does not evaluate Go scope.

```dreego
Component Card (title string)
<div><@Card title={42}/></div>
```

Error:

```text
routes/index.dreego:4:18: Card title: expected string, got int
```

The HTML attribute expression is escaped before emission. Escaping prevents
attribute injection, but it does not make an arbitrary URL trustworthy. Before
using untrusted input in `href`, `src`, or similar attributes, parse the URL and
allow only the schemes and forms the application expects, such as relative URLs
and `https`. Reject dangerous schemes such as `javascript`.

## Rules

1. **`Component` line** — Always line 1 of the file.
2. **`<@Name>`** — Component call. `@` prefix distinguishes from HTML tags.
3. **File-based Discovery** — `dreego/components/Card.dreego` → `<@Card>`.
4. **Scoped Styles** — `data-scope` per component. No leak to parent.
5. **Self-closing** — `<@Icon name="star"/>` when no body.
6. **Slots** — `{#slot}` in component template = child content.
7. **Expressions** — `{{ expression }}` renders escaped text or an HTML attribute value.
8. **Typed props** — `prop={expression}` passes a Go expression value without converting it to a string.
9. **Named prop contract** — order-independent, extra/missing props fail at `dreego generate`.

## Named Prop Contract

Component props are **named** and **order-independent**. The set of props passed in a call is validated against the component declaration at `dreego generate`.

- Missing required props are errors.
- Unknown props are errors.
- Duplicate props are errors.
- Prop order in the call does not have to match the declaration order.

**Component:**

```
Component Card (title string)
```

**Valid call:**

```dreego
<div><@Card title="Items"/></div>
```

**Invalid calls:**

```dreego
<div><@Card title="Items" count={3}/></div>
<div><@Card title="Items" title="Again"/></div>
<div><@Card/></div>
```

Error examples from `dreego generate`:

```text
routes/index.dreego:4:3: Card "count": unknown prop "count"
routes/index.dreego:5:3: Card "title": duplicate prop "title"
routes/index.dreego:6:3: Card "title": missing required prop "title"
```

The error includes the source file, line, column, component name, and prop name so the failure is easy to locate without running `go build`.

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
        <h2>{{ title }}</h2>
        {#slot}
    </article>
</div>
```

**Route:**
```html
import Card "components/Card.dreego"

<div>
<@Card title="Hi">
    {#slot header}<strong>HEADER</strong>{/slot}
    <p>Default content here</p>
</@Card>
</div>
```

- `{#slot header}{/slot}` — Placeholder in component (empty body)
- `{#slot header}content{/slot}` — Definition in the route (with body)
- `{#slot}` — Default slot (no `{/slot}` needed)
- Multiple named slots per component possible
- Unknown named slots fail at `dreego generate` with error `routes/index.dreego:4:3: Card: unknown slot "footer"`
- Nested slot declarations inside slot content are not allowed and fail with `nested slot declaration is not allowed`
- Slot content is scoped to exactly one component invocation; sibling calls never inherit another call's slot content

## Generated Go Code

```go
func Card(title string) dreego.Component {
    return dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
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

## Context Variable

Inside a component, the SSRContext is available as **`ctx`** — in routes it is called **`c`** (see [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md)). The generated render function always receives it as `ctx`:

```dreego
Component Greeting (name string)
<go>
    greeting := "Hello, " + ctx.Query("lang")
</go>
<div>
    <h1>{{ greeting }}, {{ name }}!</h1>
</div>
```

All SSRContext methods (`ctx.Param`, `ctx.Query`, `ctx.Set`, `ctx.Get`, …) are available under this name. Using `c` inside a component body produces a compile error (`undefined: c`).
