# Head Section

`<head>` contains HTML metadata for a rendered page. `lang="html"` is the
default and only supported head language.

```html
<server>pageTitle := "Account"</server>
<head>
<title>{{ pageTitle }}</title>
<meta name="description" content="Manage your account">
</head>
```

Expressions use the same context-aware escaping as body templates. Head content
from a route is combined with its layout instead of being discarded. For the two
single-value tags a page must own, the route wins: when the route supplies a
`<title>` or a `<meta name="description">`, the layout copy is dropped so the
merged document contains exactly one of each. This holds whether the layout
declares `<head>` at root level or inside a body-level `<body>` skeleton. Tag and
attribute names are matched case-insensitively, as HTML requires. Dreego does not
deduplicate any other metadata, so the application remains responsible for
intentional canonical-link and social-card policies.

Runtime dedupe captures the layout markup around `{#head}` as a literal. When a
body-level layout places a component, expression, or message between its head
tag and `{#head}`, Dreego cannot capture that markup and emits a generate-time
warning naming the layout; the layout copy is then emitted unchanged and a route
head may duplicate it. Keep `{#head}` in a plain text node after the literal
head markup to retain dedupe.

The section contains head *children*, not another `<head>` element. A layout may
place `{#head}` where route metadata belongs. Without a layout, generated output
contains the rendered head fragment followed by the body wrapper; Dreego does
not invent a complete HTML document around it.

See [Layouts](layouts.md) for head composition and [Output Safety](../../_docs/security.md)
for expression contexts.
