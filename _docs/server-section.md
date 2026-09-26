# Server Section

`<server>` contains request-time Go. `lang="go"` is the default and only
supported server language.

```html
<server>
name := c.Param("name")
if name == "" {
    name = "World"
}
</server>

<body><h1>Hello {{ name }}</h1></body>
```

Route server code receives `c`, the generated route context. It exposes request
data, route parameters, response helpers, sessions, validation data, and other
documented SSR operations. Values declared in a server section can be used by
the matching head and body templates.

## Scope

Each route directory under `www/routes` compiles into **its own Go package**
(named after the folder, sanitized to a valid Go identifier). The top-level
`www/routes/dree.go` collects the `Register(app)` call of every sub-package, so
the website entry point still calls one `routes.Register`. A `<server>` section
is split at generation time:

- The **leading declaration block** — `type`, `func`, `var`, and `const`
  declarations at the top of the section, plus any top-level `func` — is emitted
  at package level. It is visible to the other route files **in the same
  folder** and must be unique within that folder. The same name may be reused in
  a different route folder, because each folder is a separate package.
- **Statements** and any `var` that follows a statement stay inside the render
  function and are local to one request.

This makes shared types, helpers, and stores possible from `<server>` alone:

```html
GOIMPORT { sync }

<server>
type Store struct {
    mu    sync.Mutex
    items []string
}

var store = &Store{}

func (s *Store) Add(item string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.items = append(s.items, item)
}
</server>
```

A shared store is package-level mutable state across requests; synchronize it
with a mutex. Keep request-local `var`s after a statement so they stay inside
the render function.

Because folders are separate packages, two route folders may each declare their
own `type Product`. Within **one** folder the name must be declared once;
declare it again with the same name in the same folder and generation fails. To
share a type or store across folders, put it in a package both folders import
with `GOIMPORT` (see the app-code recipe). Directories that cannot be Go
packages — dynamic segments such as `[id]` and groups such as `(group)` — fold
into their nearest valid ancestor package. See [Routing](routing.md) for the
per-route file contract.

Component server code receives the component's typed props and render context;
it is not a second HTTP handler. See [Components](components.md) for the smaller
component boundary.

## Methods

Sections without `method` belong to GET. A route file can define other methods
explicitly:

```html
<server method="post">
result := saveForm(c)
</server>
<body method="post"><p>{{ result }}</p></body>
```

The server and body sections for one method form one route response. See
[Routing](routing.md) and [Forms](forms.md).

## Typed responses

`type="json"`, `type="xml"`, and custom response helpers allow a route to
respond without an HTML body. Content negotiation and response methods are
documented in [Routing](routing.md) and [Runtime API](runtime.md).
