# dreego.config.json

The configuration file is located at `dreego.config.json` in the project root.

## Schema

```json
{
  "logging": {
    "enabled": true
  },
  "redirects": [
    { "from": "/old", "to": "/new", "status": 301 }
  ],
  "rewrites": [
    { "from": "/api/v1/*", "to": "/api/v2/*" }
  ]
}
```

## logging

| Field     | Type   | Default | Description                         |
|-----------|--------|---------|-------------------------------------|
| `enabled` | bool   | `true`  | Enable RequestLogging middleware    |

When `false`, no request logging is performed. A plugin (`plugin-logging`) can take over the logging function in a future version.

## redirects

| Field    | Type   | Description                                          |
|----------|--------|------------------------------------------------------|
| `from`   | string | Source path (exact) or `/*` wildcard                 |
| `to`     | string | Target path (exact) or `/*` wildcard                 |
| `status` | int    | Redirect status: `301`, `302`, `303`, `307`, `308`   |

### Semantics

- **Exact** (`from` without `/*`): matches only the exact path. `/api`
  matches `/api` but not `/api/v1` and not `/apiary`.
- **Wildcard** (`from` ending in `/*`): matches the base path and any
  path below it at a **segment boundary**. `/api/*` matches `/api` and
  `/api/users/1` but does **not** match `/apiary` (near-prefix collision).
- The suffix after the matched prefix is appended to the target prefix.
  `/api/*` → `/v2/*` rewrites `/api/users/1` to `/v2/users/1`.

### Validation (configuration time)

Invalid rules fail during `RegisterRedirect` (at generation/startup), never
at request time:

- Empty or non-`/`-prefixed `from`/`to`
- Trailing `/` (except root `/`), double slashes `//`
- `*` outside a trailing `/*`, or multiple `/*`
- Status outside `{301, 302, 303, 307, 308}`
- Self-loops (`from == to`) and wildcard loops (`/api/*` → `/api/v2/*`,
  which redirects back under the same prefix)

## rewrites

| Field  | Type   | Description                               |
|--------|--------|-------------------------------------------|
| `from` | string | Source path (exact) or `/*` wildcard      |
| `to`   | string | Target path (exact) or `/*` wildcard      |

Rewrites share the exact/wildcard semantics and validation rules of
redirects (minus the status code). A rewrite changes the request path
transparently before the router sees it; redirects send an HTTP redirect
response.

## Ordering vs. middleware

Rewrites are applied just before routing, after user middleware registered via
`app.Use`. Middleware therefore sees the **original** request path, not the
rewritten one — match middleware patterns against the source path (for example
match `/api/*` even when `/api/*` rewrites to `/v2/*`). Access logs record the
pre-rewrite path.
