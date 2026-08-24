# Layouts

Layouts are shared shells rendered around route content. A layout lives in a
`layouts` directory inside the project root and uses two special placeholders:

- `{#slot}` — where the route content is injected.
- `{#head}` — where the route's `<head>` markup is merged.

## Layout Discovery

Layout discovery is restricted to the project's `dreego/` tree. Layout files
outside the project root (e.g. `vendor/…/www/layouts`, `subapp/www/layouts`)
are ignored.

A layout file is named `default.dreego` (or the legacy `layout.dreego`). Layouts
resolve per route by a route-local cascade:

1. The route's own scope (e.g. `www/routes/blog/layouts/default.dreego` for
   `www/routes/blog/…`).
2. Each parent route scope up to the root (`www/layouts/default.dreego`).

The first matching layout in the cascade wins. Only one layout file per scope
is allowed: `default.dreego` and `layout.dreego` in the same `layouts`
directory is an ambiguous-layout error and fails `dreego generate` with a
diagnostic naming both files.

## Syntax

**`www/layouts/default.dreego`:**

```html
<body>
<!DOCTYPE html>
<html>
<head>
    {#head}
</head>
<body>
    {#slot}
</body>
</html>
</body>
```

The layout defines the outer `<html>`/`<head>`/`<body>` skeleton. At codegen time the route's page content is placed into `{#slot}` and the route's `<head>` sections are merged into `{#head}`.

## Route Head Behavior

- **With layout**: the route's `<head>` content (e.g. `<title>{{ doc.Title }}</title>`) is injected into the layout's `{#head}` placeholder. Expressions in the head are resolved and escaped.
- **Without layout**: if no layout file exists, the route's `<head>` is emitted standalone as a full `<html>` fragment, so the page still renders with its title and meta tags.
- **No `<head>` in route**: when the route declares no `<head>`, nothing is injected into `{#head}`.

## Generated Go Code

The layout wrapping is emitted as `c.Set("slot", pageContent)` / `c.Set("head", headContent)`, and the layout template reads them back with `c.Get("slot")` / `c.Get("head")`.

## Rules

1. `{#slot}` — required to render route content; always available.
2. `{#head}` — optional; collects the route's `<head>` sections.
3. Route `<head>` works with or without a layout.
4. One layout file per `layouts` directory; `default.dreego` and `layout.dreego` together is an error.
5. Layout lookup is route-local and cascades through documented parent directories.
