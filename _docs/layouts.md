# Layouts

Layouts are shared shells rendered around route content. A layout lives in `dreego/layouts/default.dreego` and uses two special placeholders:

- `{#slot}` — where the route content is injected.
- `{#head}` — where the route's `<head>` markup is merged.

## Syntax

**`dreego/layouts/default.dreego`:**

```html
<!DOCTYPE html>
<html>
<head>
    {#head}
</head>
<body>
    {#slot}
</body>
</html>
```

The layout defines the outer `<html>`/`<head>`/`<body>` skeleton. At codegen time the route's page content is placed into `{#slot}` and the route's `<head>` sections are merged into `{#head}`.

## Route Head Behavior

- **With layout**: the route's `<head>` content (e.g. `<title>{doc.Title}</title>`) is injected into the layout's `{#head}` placeholder. Expressions in the head are resolved and escaped.
- **Without layout**: if no layout file exists, the route's `<head>` is emitted standalone as a full `<html>` fragment, so the page still renders with its title and meta tags.
- **No `<head>` in route**: when the route declares no `<head>`, nothing is injected into `{#head}`.

## Generated Go Code

The layout wrapping is emitted as `c.Set("slot", pageContent)` / `c.Set("head", headContent)`, and the layout template reads them back with `c.Get("slot")` / `c.Get("head")`.

## Rules

1. `{#slot}` — required to render route content; always available.
2. `{#head}` — optional; collects the route's `<head>` sections.
3. Route `<head>` works with or without a layout.
4. Only one layout file is used per route group (`default.dreego`).
