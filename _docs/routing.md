# File-based Routing

## Ordnerstruktur

```
dreego/routes/
├── get.dreego              → GET /
├── about.dreego              → GET /about
├── users/
│   └── [id]/get.dreego       → GET /users/{id}
├── blog/
│   └── [...catchall].dreego  → GET /blog/{catchall}
└── (admin)/
    └── dashboard.dreego      → GET /dashboard (group ignored in path)
```

## Dynamische Segmente

| Syntax               | Pfad                    | Go-Param              |
|----------------------|-------------------------|-----------------------|
| `[id]`               | `/users/{id}`           | `c.Param("id")`       |
| `[...catchall]`      | `/blog/{catchall}`      | `c.Param("catchall")` |
| `[[optional]]`       | `/docs/{optional}`      | `c.Param("optional")` |
| `(group)/`           | unsichtbar im Pfad      | —                     |

## HTTP-Methoden

Aus Dateinamen abgeleitet:

```
users.get.dreego   → GET
users.post.dreego  → POST
users.put.dreego   → PUT
users.delete.dreego → DELETE
```

Alternativ via `<go method="post">` in der `.dreego`-Datei.

## API-Routen

Pfade mit `api/` rendern kein Layout (nur `<div>`-Fragment):

```
dreego/routes/api/users.get.dreego → GET /api/users
```

## Plugin-Routen

Plugins registrieren Routen via `init()` im eigenen Go-Package. `dreego generate` fugt den Import automatisch in `dreego/gen/dree.go` ein.
