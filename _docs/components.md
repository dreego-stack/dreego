# Components

Component = wiederverwendbare `.dreego`-Datei mit Props, Slots und eigenem Scope.

## Declaration

```html
Component Name (prop Type, prop Type = default) (slot1, slot2)
```

Die erste Zeile einer Component-Datei **muss** die `Component`-Deklaration sein.

| Element | Beschreibung |
|---------|-------------|
| `Name` | Component-Name (PascalCase). Wird zu `<@Name>`. |
| `(prop Type)` | Props mit Go-Typ. `= default` optional. |
| `(slot1, slot2)` | Benannte Slots. Default-Slot immer via `{#slot}`. |

**Beispiel:**

```html
Component Card (title string, image url = "") (header, footer)

<head></head>
<go>
    formatted := strings.Title(title)
</go>
<div>
    <article class="card">
        {#slot header}
            <h2>{formatted}</h2>
            {#if image}
                <img src="{image}" alt="{title}"/>
            {/if}
        {/slot}
        {#slot}
        {#slot footer}
            <small>{#slot}</small>
        {/slot}
    </article>
</div>
<style>
.card { border: 1px solid var(--border); padding: 1rem; }
</style>
```

## Usage

```html
import Card "./components/Card"
import ui "github.com/dreego-ecosystem/dreego-ui"

<go>
    user := LoadUser(ctx.Param("id"))
</go>
<div>
    <@Card title={user.Name} image={user.Avatar}>
        {#slot header}
            <@ui.Badge variant="pro">PRO</@ui.Badge>
        {/slot}
        <p>{user.Bio}</p>
        {#slot footer}
            <small>Member since {user.Joined}</small>
        {/slot}
    </@Card>
</div>
```

## Regeln

1. **`Component`-Zeile** — Immer Zeile 1 der Datei. Keine Leerzeile davor.
2. **`import`** — Nur am Datei-Anfang, nach `Component`-Zeile (falls Component) oder vor allen Tags (falls Route).
3. **`<go>`** — Reiner Go-Code, keine `func`- oder `return`-Statements. Props sind als Go-Variablen verfügbar.
4. **`<@Name>`** — Component-Aufruf. `@`-Prefix unterscheidet von HTML-Tags.
5. **`{#slot name}`** — Benannter Slot. Ohne Namen = Default-Slot.
6. **File-based Discovery** — `dreego/components/Card.dreego` → `<@Card>`.
7. **External Components** — `import ui "url"` → `<@ui.Button>`. Alias + `.` + Name.
8. **Scoped Styles** — `data-scope` pro Component (wie bei Routes). Kein Leak in Parent.
9. **Self-closing** — `<@Icon name="star"/>` wenn kein Body.
10. **Kein `return`** — Component implizit. Der `<div>`-Inhalt ist der Output.

## Wie es kompiliert

Aus `Card.dreego` wird eine Go-Funktion:

```go
func Card(title string, image url) core.Component {
    return core.NewComponent(func(ctx *core.SSRContext) (string, error) {
        formatted := strings.Title(title)
        // ... generated template code ...
        return core.RenderSlots(ctx, map[string]core.Component{
            "header": headerSlot,
            "footer": footerSlot,
            "":       defaultSlot,
        })
    })
}
```

Aufruf `<@Card title="x">` wird zu `Card("x").Render(ctx, w)`.
