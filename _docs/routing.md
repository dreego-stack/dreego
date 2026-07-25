# File-based Routing

## Ordnerstruktur

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
        └── get.dreego          → GET /demo  (group ignoriert)
```

## Dynamische Segmente

| Syntax               | URL-Pattern              | Go-Param              |
|----------------------|--------------------------|-----------------------|
| `[id]`               | `/users/{id}`            | `c.Param("id")`       |
| `[...catchall]`      | `/blog/{catchall...}`    | `c.Param("catchall")` |
| `[[optional]]`       | `/{optional}`            | `c.Param("optional")` |

## Route-Groups `(name)/`

Gruppen-Ordner erscheinen **nicht** in der URL. Sie dienen der Code-Organisation:

```
(admin)/            → Layout + Middleware nur fur Admin-Bereich
(auth)/             → Auth-Check fur Login/Register
```

## HTTP-Methoden

Jede Route hat eine Methoden-Datei im Verzeichnis:

```
get.dreego     → GET
post.dreego    → POST
put.dreego     → PUT
delete.dreego  → DELETE
```

Mehrere Methoden pro Verzeichnis moglich:

```
users/
└── [id]/
    ├── get.dreego      → GET /users/{id}
    └── delete.dreego   → DELETE /users/{id}
```

## Error-Pages

- `404.dreego` in einem Route-Verzeichnis → Catch-All fur diesen Pfad
- Go Mux wahlt den spezifischsten Catch-All: `users/404.dreego` vor `routes/404.dreego`
- `500.dreego` (nur eine global) → Recovery-Middleware rendert bei Panic
- Error-Pages bekommen **kein** Layout (Endlos-Rekursion vermeiden)
