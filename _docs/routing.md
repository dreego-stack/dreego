# File-based Routing

> **Current implementation:** This page documents the released pre-v0.1
> filename-based router. The accepted v0.1 migration will use one route file per
> URL, with flat files or `+page.dreego` and method-specific `<go>` and `<div>`
> sections. That target is not available until `routing-correctness.1` and the
> App migration are implemented.

Route discovery is restricted to the project's `dreego/routes/` tree.
Directories named `routes` outside the project root (e.g. `vendor/…/dreego/routes`,
`node_modules/…/dreego/routes`, `subapp/dreego/routes`) are ignored. Only the
top-level `dreego/` directory is treated as the project root.

Directories below `dreego/routes/` define the URL path. The filename defines
the HTTP method. Keeping one method per file prevents a route file from growing
into a combined implementation for every operation on the same URL.

## Directory Structure

```
dreego/routes/
├── get.dreego                  → GET /
├── 404.dreego                  → GET /* (catch-all)
├── 500.dreego                  → Panic → 500
├── about/
│   └── get.dreego              → GET /about
├── users/
│   ├── 404.dreego              → GET /users/* (catch-all)
│   └── [id]/
│       └── get.dreego          → GET /users/{id}
├── blog/
│   └── [...catchall]/
│       └── get.dreego          → GET /blog/{catchall...}
└── (group)/
    └── demo/
        └── get.dreego          → GET /demo  (group ignored)
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

Each route has a method file in the directory:

```
get.dreego     → GET
post.dreego    → POST
put.dreego     → PUT
delete.dreego  → DELETE
```

Multiple methods per directory possible:

```
users/
└── [id]/
    ├── get.dreego      → GET /users/{id}
    └── delete.dreego   → DELETE /users/{id}
```

## Error Pages

- `404.dreego` in a route directory → Catch-all for this path
- Go mux chooses the most specific catch-all: `users/404.dreego` before `routes/404.dreego`
- `500.dreego` (only one global) → Recovery middleware renders on panic
- Error pages get **no** layout (avoids infinite recursion)

## Content-Type Routing (v0.0.15)

A single route can serve multiple content types via `<go type="...">` blocks:

```dreego
<go>
    user := db.GetUser(c.Param("id"))
</go>

<go type="json">
    c.JSON(200, user)
</go>

<div>
    <h1>{{ user.Name }}</h1>
</div>
```

| type | MIME | Behavior |
|------|------|----------|
| `json` | `application/json` | `c.JSON()`, `c.Bind()`, auto-detect via `Accept` header |
| `xml` | `application/xml` | `c.XML()`, auto-detect via `Accept` header |
| *(none)* | `text/html` | Default — renders `<div>` template |

- `<go>` without `type` runs **always** (shared logic)
- Typed `<go>` blocks run conditionally based on `Accept` header
- Pure JSON/XML routes (no `<div>`) skip template rendering entirely
- Raw content: `c.Write(status, contentType, body)` for FlatBuffers/Protobuf/etc.
- Content negotiation: `c.Wants(mime)` available in any `<go>` block
