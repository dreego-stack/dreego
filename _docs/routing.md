# File-based Routing

> **Current implementation:** `+page.dreego` and `index.dreego` define a
> directory URL. Every other `.dreego` filename defines a literal URL segment.
> Routes support method-specific `<server>` and `<body>` sections.

Route discovery is restricted to the website root's `routes/` tree. The
website root is any directory containing `dreego.config.json`. Directories
named `routes` outside a website root (e.g. `vendor/…/www/routes`,
`node_modules/…/www/routes`, `subapp/www/routes`) are ignored.

Directories below `www/routes/` define the URL path. `+page.dreego` and
`index.dreego` define the URL of the directory itself. A different `.dreego`
filename defines the final static path segment. A directory containing both
index filenames is rejected as a duplicate route.

## Directory Structure

```
www/routes/
├── +page.dreego                 → GET /
├── 404.dreego                  → GET /* (catch-all)
├── 500.dreego                  → Panic → 500
├── about/
│   └── +page.dreego             → GET /about
├── users/
│   ├── 404.dreego              → GET /users/* (catch-all)
│   └── [id]/
│       └── +page.dreego         → GET /users/{id}
├── blog/
│   └── [...catchall]/
│       └── +page.dreego         → GET /blog/{catchall...}
└── (group)/
    └── demo/
        └── +page.dreego         → GET /demo  (group ignored)
```

## Dynamic Segments

| Syntax               | URL-Pattern              | Go-Param              |
|----------------------|--------------------------|-----------------------|
| `[id]`               | `/users/{id}`            | `c.Param("id")`       |
| `[...catchall]`      | `/blog/{catchall...}`    | `c.Param("catchall")` |

Optional segments are not supported. Define each route explicitly so one
method file always owns one route pattern.

## Route Groups `(name)/`

Group directories do **not** appear in the URL. They serve code organization:

```
(admin)/            → Layout + middleware only for admin area
(auth)/             → Auth check for login/register
```

## HTTP Methods

One route file may define multiple methods. Sections without `method` default
to GET. A request renders only the sections matching its method:

```dreego
<server>
    page := loadPage(c)
</server>
<body>{{ page.Title }}</body>

<server method="post">
    result := savePage(c)
</server>
<body method="post">Saved: {{ result }}</body>
```

Components, imports, layouts, styles, and scripts remain route-level resources.
The method controls only the route logic and rendered `<body>` section.

The HTTP method is declared with `method` attributes on `<server>`/`<body>`
sections (default GET):

```
+page.dreego      → directory URL, GET by default
index.dreego      → directory URL, GET by default
page.dreego       → /page, GET by default
account.dreego    → /account, GET by default
```

Filenames never select an HTTP method. Named routes default to GET. Add
`method="post"` sections to handle POST requests.

Multiple methods per route possible in one `+page.dreego`:

```
users/
└── [id]/
    └── +page.dreego      → GET /users/{id} and DELETE /users/{id}
```

## Error Pages

- `404.dreego` in a route directory → Catch-all for this path
- Go mux chooses the most specific catch-all: `users/404.dreego` before `routes/404.dreego`
- `500.dreego` (only one global) → Recovery middleware renders on panic
- Error pages get **no** layout (avoids infinite recursion)

## Content-Type Routing (v0.0.15)

A single route can serve multiple content types via `<server type="...">` blocks:

```dreego
<server>
    user := db.GetUser(c.Param("id"))
</server>

<server type="json">
    c.JSON(200, user)
</server>

<body>
    <h1>{{ user.Name }}</h1>
</body>
```

| type | MIME | Behavior |
|------|------|----------|
| `json` | `application/json` | `c.JSON()`, `c.Bind()`, auto-detect via `Accept` header |
| `xml` | `application/xml` | `c.XML()`, auto-detect via `Accept` header |
| *(none)* | `text/html` | Default — renders `<body>` template |

- `<server>` without `method` runs for GET; method-specific `<server>` blocks run only
  for their matching method
- Typed `<server>` blocks run conditionally based on `Accept` header
- Pure JSON/XML routes (no `<body>`) skip template rendering entirely
- Raw content: `c.Write(status, contentType, body)` for FlatBuffers/Protobuf/etc.
- Content negotiation: `c.Wants(mime)` available in any `<server>` block
