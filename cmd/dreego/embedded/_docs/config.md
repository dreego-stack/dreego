# dreego/config.json

The configuration file is located at `dreego/config.json` in the project root.

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

| Field    | Type   | Description                                 |
|----------|--------|---------------------------------------------|
| `from`   | string | Source path                                  |
| `to`     | string | Target path or external URL                  |
| `status` | int    | HTTP status: `301` (permanent), `302` (temp) |

## rewrites

| Field  | Type   | Description                               |
|--------|--------|-------------------------------------------|
| `from` | string | Source path pattern (e.g. `/api/v1/*`)    |
| `to`   | string | Target path (e.g. `/api/v2/*`)            |
