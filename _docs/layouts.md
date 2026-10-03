# Layouts

Layouts are shared shells rendered around route content. A layout lives in a
`layouts` directory and uses two special placeholders:

- `{#slot}` — where the route content is injected.
- `{#head}` — where the route's `<head>` markup is merged.

The shared layout lives at the website root (`dreego/layouts/`); an app may add
its own `layouts/` (`dreego/www/layouts/`) to override the shared one for that
app, and a route may add a route-local scope under `<app>/routes/<sub>/layouts/`
or under the shared `dreego/routes/<sub>/layouts/`. Each `layouts/` directory
may hold a `default.dreego`, a legacy `layout.dreego`, and any number of
explicitly named layouts (for example `admin.dreego`) that a route selects by
path.

A layout may reference components from the shared root `components/` tree and
from the app-local `components/` tree of the app that uses it.

## Layout Discovery

Layout discovery is restricted to a website root identified by
`dreego.config.json`. Layout files outside that root, including vendored modules
and nested applications, are ignored. Layouts are found at the root, on an app,
and inside a `routes/` tree (shared or app-local).

A layout file is named `default.dreego` (or the legacy `layout.dreego`) and is
selected by a cascade:

1. The route's own scope, deepest first (e.g.
   `dreego/www/routes/blog/layouts/default.dreego` for `dreego/www/routes/blog/…`).
2. A shared route-local scope under the root routes tree
   (`dreego/routes/blog/layouts/default.dreego`).
3. The app's layout (`dreego/www/layouts/default.dreego`).
4. The shared website-root layout (`dreego/layouts/default.dreego`).

The first matching layout in the cascade wins; the app-local layout therefore
always beats the shared root layout. Only one default layout file per scope is
allowed: `default.dreego` and `layout.dreego` in the same `layouts` directory is
an ambiguous-layout error and fails `dreego generate` with a diagnostic naming
both files. Explicitly named layout files may sit beside a default one.

## Explicit Layout Selection

A route may select a layout by path with the `LAYOUT` header directive instead
of relying on the cascade:

```dreego
LAYOUT "dreego/layouts/admin.dreego"

<body><p>Admin</p></body>
```

The path is resolved relative to the website root (or the working directory);
a missing target fails `dreego generate` with a `file:line:col` diagnostic
naming the path.


## Syntax

A layout file may declare its kind explicitly with the `DREEFILE layout` header
directive. The declaration is optional; layouts are still resolved by the
directory cascade described above, and layout files keep living under a
`layouts` directory.

**`dreego/layouts/default.dreego`:**

```html
DREEFILE layout

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

The doctype and the `<html>` element belong **inside** the layout's `<body>`
section — the file starts with a root section, and the document skeleton is its
content. Putting `<!DOCTYPE html>` or `<html>` before the first root section
fails with `expected root section`; the error now points to the `<body>`
placement. A layout that omits the doctype renders pages in quirks mode, so
keep `<!DOCTYPE html>` as the first line inside the `<body>` section.

## Route Body Attributes

The `<body>` section tag may carry `lang` (section language) and `method`
(dreego directives). Any other attribute on that tag (for example
`<body x-data="app()">`) sits on the route's body wrapper, but the layout
supplies the real document `<body>`, so the attribute never reaches the
rendered element. `dreego generate` warns and names the attribute. Put such
attributes on an element inside `<body>`, or add them to the layout's own body
tag.

## Route Head Behavior

- **With layout**: the route's `<head>` content (e.g. `<title>{{ doc.Title }}</title>`) is injected into the layout's `{#head}` placeholder. Expressions in the head are resolved and escaped. Both layout shapes work: a root-level `<head>` section and a body-level `<body><html><head>…{#head}…</head>…` skeleton.
- **Route overrides single-value tags**: when the route supplies a `<title>`, a `<meta name="description">`, a `<meta name="viewport">`, or a `<meta charset>`, the layout's copy is dropped so exactly one of each survives. Layout tags are kept when the route defines none. Tag and attribute names are matched case-insensitively.
- **Body-level literal requirement**: dedupe captures the layout markup around `{#head}` as a literal. If a body-level layout places a component, expression, or message between its head tag and `{#head}`, Dreego emits a generate-time warning naming the layout, skips dedupe for that markup, and the layout copy is emitted unchanged. Keep `{#head}` in a plain text node after the literal head markup to retain dedupe.
- **Without layout**: the rendered head fragment is emitted before the body
  wrapper. Dreego does not invent an `<html>` document or outer `<head>` element.
- **No `<head>` in route**: when the route declares no `<head>`, nothing is injected into `{#head}`.

## Generated Go Code

Generated route rendering passes page content and head content directly to the
selected layout renderer. Layout composition does not depend on mutable
request-context keys.

Layout renderers live in one `dreego/layouts` package and are named after their
scope so several scopes can coexist: the shared root layout is `Default` (or
`Layout` for a root `layout.dreego`), and an app or route-local layout is
suffixed with its PascalCased scope, e.g.
`dreego/www/routes/registrierung/layouts/default.dreego` becomes
`DefaultWwwRegistrierung` and a shared
`dreego/routes/registrierung/layouts/default.dreego` becomes
`DefaultRootRegistrierung`. A named layout file adds its own base name, e.g.
`dreego/layouts/admin.dreego` becomes `LayoutAdmin`. The generated
`<app>/routes/<route>/layouts/` directory is not treated as a route, so
route-local layouts do not create routes.

## Rules

1. `{#slot}` — required to render route content; always available.
2. `{#head}` — optional; collects the route's `<head>` sections.
3. Route `<head>` works with or without a layout.
4. One default layout per `layouts` directory; `default.dreego` and `layout.dreego` together is an error. Explicitly named layouts may sit beside them.
5. Layout lookup is route-local and cascades through documented parent directories; a `LAYOUT "path"` directive selects explicitly.
6. An app-local layout wins over the shared root layout.
