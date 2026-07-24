# dreego/config.json

Die Konfigurationsdatei liegt unter `dreego/config.json` im Projekt-Root.

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

| Feld      | Typ    | Default | Beschreibung                         |
|-----------|--------|---------|--------------------------------------|
| `enabled` | bool   | `true`  | RequestLogging-Middleware aktivieren |

Wenn `false`, wird kein Request-Logging durchgefuhrt. Ein Plugin (`dreego-logging`) kann die Logging-Funktion in V2 ubernehmen.

## redirects

| Feld     | Typ    | Beschreibung                                  |
|----------|--------|-----------------------------------------------|
| `from`   | string | Quell-Pfad                                    |
| `to`     | string | Ziel-Pfad oder externe URL                    |
| `status` | int    | HTTP-Status: `301` (permanent), `302` (temp)  |

## rewrites

| Feld   | Typ    | Beschreibung                                  |
|--------|--------|-----------------------------------------------|
| `from` | string | Quell-Pfad-Pattern (z.B. `/api/v1/*`)         |
| `to`   | string | Ziel-Pfad (z.B. `/api/v2/*`)                  |
