# File-based Routing

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
| `[[optional]]`       | `/{optional}`            | `c.Param("optional")` |

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
