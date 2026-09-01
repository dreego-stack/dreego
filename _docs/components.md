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

<body>
    <article class="card">
        <h2>{{ title }}</h2>
        <div>{#slot}</div>
    </article>
</body>
<style>
.card { border: 1px solid var(--border); padding: 1rem; }
</style>
```

## Usage

Component declarations and imports are header directives. They are the only
content allowed outside the five root sections: `<server>`, `<head>`, `<body>`,
`<style>`, and `<client>`. Free text, HTML, and component calls at the root are
generation errors.

```dreego
import Card "components/Card.dreego"

<body><@Card title="Hello" /></body>
```

**Self-closing:**

```html
<body><@Card title="Hello"/></body>
```

A self-closing call renders empty default and named slots. Content after the
tag belongs to the parent and remains a normal sibling.

**With children (default slot):**

```html
<body>
    <@Card title="Welcome">
        <p>Slot content goes here</p>
    </@Card>
</body>
```

## Self-closing Calls and Slot Fallback

`<@Card/>` is allowed and renders an empty default slot.

The following paragraph is a sibling of the component:

```html
<body><@Card/><p>Sibling content</p></body>
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
<body><@Card count={count} /></body>
```

**In the component body:**

```html
Component Link (url string)
<body><a href="{{ url }}">{#slot}</a></body>
```

```html
<body><@Link url="https://dreego.dev">Home</@Link></body>
```

Simple literal expressions (`"..."`, integer literals) are type-checked
against the declared prop type at `dreego generate` time. Non-literal expressions
are accepted unchecked because the transpiler does not evaluate Go scope.

```dreego
Component Card (title string)
<body><@Card title={42}/></body>
```

Error:

```text
routes/index.dreego:4:18: Card title: expected string, got int
```

The HTML attribute expression is escaped before emission. Escaping prevents
attribute injection, but it does not make an arbitrary URL trustworthy. URL
attributes (`href`, `src`, `action`, …) are additionally scheme-validated:
values with unsafe schemes such as `javascript:` are replaced with `#`. See
[Output Safety](https://github.com/dreego-stack/dreego/blob/main/_docs/security.md)
for the exact context rules and the `|raw` opt-in.

## Rules

1. **`Component` line** — Always line 1 of the file.
2. **`<@Name>`** — Component call. `@` prefix distinguishes from HTML tags.
3. **File-based Discovery** — `www/components/Card.dreego` → `<@Card>`.
4. **Scoped Styles** — `data-scope` per component. No leak to parent.
5. **Self-closing** — `<@Icon name="star"/>` when no body.
6. **Slots** — `{#slot}` in component template = child content.
7. **Expressions** — `{{ expression }}` renders escaped text or an HTML attribute value, with context-aware rules (see [Output Safety](https://github.com/dreego-stack/dreego/blob/main/_docs/security.md)).
8. **Typed props** — `prop={expression}` passes a Go expression value without converting it to a string.
9. **Named prop contract** — order-independent, extra/missing props fail at `dreego generate`.

## Accessibility

Component markup remains responsible for semantic intent. Use labeled
navigation landmarks, ordered headings, associated form labels, alternative
text for informative images, and `aria-hidden="true"` for decorative content.
Use buttons for actions and links for navigation.

A page shell should provide one `<main id="main">` landmark and a skip link as
its first focusable element:

```dreego
Component PageShell (title string)

<body>
    <a href="#main" class="skip-link">skip to content</a>
    <header>...</header>
    <main id="main">{#slot}</main>
</body>
```

Visually hide the link until it receives focus:

```css
.skip-link { position: absolute; left: -9999px; }
.skip-link:focus { position: fixed; top: 0.5rem; left: 0.5rem; }
```

Dreego does not make arbitrary user applications automatically accessible.
See
[Accessibility](https://github.com/dreego-stack/dreego/blob/main/_docs/accessibility.md)
for the complete requirements and testing guidance.

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
<body><@Card title="Items"/></body>
```

**Invalid calls:**

```dreego
<body><@Card title="Items" count={3}/></body>
<body><@Card title="Items" title="Again"/></body>
<body><@Card/></body>
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
<body class="spinner"></body>
<style>
@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }
.spinner { animation: spin 1s linear infinite; }
</style>
```

## Named Slots (v0.0.8)

**Component:**
```
Component Card (title string)

<body>
    <article>
        {#slot header}{/slot}
        <h2>{{ title }}</h2>
        {#slot}
    </article>
</body>
```

**Route:**
```html
import Card "components/Card.dreego"

<body>
<@Card title="Hi">
    {#slot header}<strong>HEADER</strong>{/slot}
    <p>Default content here</p>
</@Card>
</body>
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
    return dreego.ComponentFunc(func(ctx dreego.RenderContext) (dreego.Result, error) {
        var b strings.Builder
        b.WriteString("<div data-scope=\"abc123\">")
        b.WriteString(`<article class="card">`)
        b.WriteString(`<h2>`)
        b.WriteString(dreego.SafeText(fmt.Sprintf("%v", title)))
        b.WriteString(`</h2>`)
        b.WriteString(ctx.Get("slot"))
        b.WriteString(`</article>`)
        b.WriteString("</div>")
        return dreego.Result{HTML: []byte(b.String())}, nil
    })
}
```

Call `<@Card title="x"/>` → `Card("x").Render(c)`.

Hand-written `Component` and `ComponentFunc` implementations must use this
`Render(dreego.RenderContext) (dreego.Result, error)` signature; they are not
regenerated, so update them manually.

## Context Variable

Inside a component, the SSRContext is available as **`ctx`** — in routes it is called **`c`** (see [Runtime API](https://github.com/dreego-stack/dreego/blob/main/_docs/runtime.md)). The generated render function always receives it as `ctx`:

```dreego
Component Greeting (name string)
<server>
    greeting := "Hello, " + ctx.Query("lang")
</server>
<body>
    <h1>{{ greeting }}, {{ name }}!</h1>
</body>
```

All SSRContext methods (`ctx.Param`, `ctx.Query`, `ctx.Set`, `ctx.Get`, …) are available under this name. Using `c` inside a component body produces a compile error (`undefined: c`).
