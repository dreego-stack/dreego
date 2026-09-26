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
  is still the directory cascade described in [Layouts](layouts.md).
- `COMPONENT "www/components" IMPORT { Card, Card as ProductCard }` imports
  components from a path; `as` creates an alias.
- `GOIMPORT { path1, alias "path2" }` declares Go imports for the file's
  generated package. Any package that resolves in the application is accepted:
  the standard library, packages in your own module, and dependencies listed in
  `go.mod`. An entry may carry an alias (`alias "path"`); the base name is used
  when no alias is given. A package that is not in `go.mod` fails at
  `dreego generate` with a diagnostic naming the path, `not in go.mod`, and the
  matching `go get` command. Two imports whose aliases or base names collide
  without an explicit alias fail with a diagnostic naming both paths. Dreego
  never runs `go get` automatically. `strings`, `net/http`, and `fmt` are also
  detected automatically from the code.
- `PROFILE "name"` binds the route folder to a named profile registered with
  `app.Profile(name, dreego.Profile{…})` in `main.go`. A profile selects its own
  session store, CSRF switch, and cookie policy for the routes in that folder.
  The nearest ancestor `PROFILE` wins; a `PROFILE` in a `(group)/` folder applies
  to every descendant route folder. Without any `PROFILE` directive the app keeps
  the global session and CSRF behavior. A route bound to a profile that was
  never registered fails at start-up with an error naming the profile.

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

- [`server`](server-section.md): request-time Go and typed responses.
- [`head`](head-section.md): document metadata merged with layouts.
- [`body`](body-html.md): HTML templates and Dreego template constructs.
- [`style`](style-section.md): route or component CSS.
- [`client`](client-javascript.md): browser JavaScript.
- [TypeScript client code](client-typescript.md): checked and compiled to JavaScript.
- [Browser Lua](lua.md): compiled to JavaScript with a feature-linked runtime.
- [Markdown bodies and `<md>` regions](markdown.md): compiled into HTML IR.

Only documented section/language pairs are accepted. Language names describe
inputs, not additional output targets: Markdown becomes HTML IR, while
TypeScript and Lua become JavaScript.

## Routes, layouts, and components

Route files live below the configured website root's `routes/` directory.
Directories define URL segments, while `+page.dreego` or `index.dreego` owns
the directory URL. Any other route filename adds a literal URL segment.
Layouts and components use the same semantic sections but have different
composition rules. Continue with [Routing](routing.md), [Layouts](layouts.md),
and [Components](components.md).
