# Components

Component = wiederverwendbare `.dreego`-Datei mit Props und eigenem Scope.

## Declaration

```html
Component Name (prop Type, prop Type = default)
```

Die erste Zeile einer Component-Datei **muss** die `Component`-Deklaration sein.

| Element | Beschreibung |
|---------|-------------|
| `Name` | Component-Name (PascalCase). Wird zu `<@Name>`. |
| `(prop Type)` | Props mit Go-Typ. `= default` optional. |
| `{#slot}` | Default-Slot (immer verfügbar). |

**Beispiel:**

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

Components werden automatisch aus `dreego/components/` entdeckt. Kein `import` nötig.

**Self-closing:**

```html
<div><@Card title="Hello"/></div>
```

**Mit Children (Default-Slot):**

```html
<div>
    <@Card title="Welcome">
        <p>Slot content goes here</p>
    </@Card>
</div>
```

## Regeln

1. **`Component`-Zeile** — Immer Zeile 1 der Datei.
2. **`<@Name>`** — Component-Aufruf. `@`-Prefix unterscheidet von HTML-Tags.
3. **File-based Discovery** — `dreego/components/Card.dreego` → `<@Card>`.
4. **Scoped Styles** — `data-scope` pro Component. Kein Leak in Parent.
5. **Self-closing** — `<@Icon name="star"/>` wenn kein Body.
6. **Slots** — `{#slot}` im Component-Template = Kinder-Inhalt.

## Generierter Go-Code

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

Aufruf `<@Card title="x"/>` → `Card("x").Render(c)`.
