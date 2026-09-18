# Anatomy of a `.dreego` File

A `.dreego` file groups code by purpose. Root sections select where code runs;
the optional `lang` attribute selects how that section is compiled.

## File header

Before the first root section a file may carry Dreefile header directives. The
surface rule is: **uppercase is a keyword, lowercase is a value.**

- `DREEFILE component (props)` declares a component file; the component name
  comes from the filename (`Card.dreego` becomes `<@Card>`).
- `DREEFILE layout` declares a layout file.
- `DREEFILE page` declares a page explicitly; omitting the `DREEFILE` line also
  means page.
- `LAYOUT "www/layouts/admin.dreego"` records an explicit layout path. It is
  parsed and reserved for the upcoming layout-chaining slice; layout selection
  is still the directory cascade described in [Layouts](../core/_docs/layouts.md).
- `COMPONENT "www/components" IMPORT { Card, Card as ProductCard }` imports
  components from a path; `as` creates an alias.
- `GOIMPORT { sync, encoding/json }` declares Go imports for the file's
  generated package. Only allow-listed standard-library packages are accepted;
  an unknown package fails at `dreego generate` with a diagnostic naming the
  supported set. The list is fixed (`bytes`, `context`, `encoding/base64`,
  `encoding/hex`, `encoding/json`, `errors`, `fmt`, `html`, `io`, `log`,
  `maps`, `math`, `net/http`, `net/url`, `path`, `path/filepath`, `regexp`,
  `slices`, `sort`, `strconv`, `strings`, `sync`, `time`, `unicode`,
  `unicode/utf8`), so arbitrary or dynamic imports stay impossible. `strings`,
  `net/http`, and `fmt` are also detected automatically from the code.

Header directives are the only content allowed alongside the five root sections.

```html
<server lang="go">
title := "Dashboard"
</server>

<head lang="html">
<title>{{ title }}</title>
</head>

<body lang="html">
<main><h1>{{ title }}</h1></main>
</body>

<style lang="css">
main { max-width: 64rem; margin-inline: auto; }
</style>

<client lang="js">
console.log("ready")
</client>
```

The defaults are Go, HTML, HTML, CSS, and JavaScript respectively. A route may
omit sections it does not need. A rendered HTML route normally needs a body;
typed JSON, XML, or custom responses can be produced by server sections alone.

## Section reference

- [`server`](../core/_docs/server-section.md): request-time Go and typed responses.
- [`head`](../core/_docs/head-section.md): document metadata merged with layouts.
- [`body`](../core/_docs/body-html.md): HTML templates and Dreego template constructs.
- [`style`](../core/_docs/style-section.md): route or component CSS.
- [`client`](../cmd/dreego/_docs/client-javascript.md): browser JavaScript.
- [TypeScript client code](../cmd/dreego/_docs/client-typescript.md): checked and compiled to JavaScript.
- [Browser Lua](../cmd/dreego/_docs/lua.md): compiled to JavaScript with a feature-linked runtime.
- [Markdown bodies and `<md>` regions](../cmd/dreego/_docs/markdown.md): compiled into HTML IR.

Only documented section/language pairs are accepted. Language names describe
inputs, not additional output targets: Markdown becomes HTML IR, while
TypeScript and Lua become JavaScript.

## Routes, layouts, and components

Route files live below the configured website root's `routes/` directory.
Directories define URL segments, while `+page.dreego` or `index.dreego` owns
the directory URL. Any other route filename adds a literal URL segment.
Layouts and components use the same semantic sections but have different
composition rules. Continue with [Routing](../core/_docs/routing.md), [Layouts](../core/_docs/layouts.md),
and [Components](../core/_docs/components.md).
