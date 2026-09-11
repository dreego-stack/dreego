# Anatomy of a `.dreego` File

A `.dreego` file groups code by purpose. Root sections select where code runs;
the optional `lang` attribute selects how that section is compiled.

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
