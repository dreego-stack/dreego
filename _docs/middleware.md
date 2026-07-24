# Middleware

## Core-Middleware

Dreego hat eine Middleware-Chain mit fester Reihenfolge:

```
[Recovery → RequestID → RealIP → RequestLogging*]
  → Redirect/Rewrite
    → Router → Handler
```

\* `RequestLogging` ist Core-Conditional: default aktiv, abschaltbar via `dreego/config.json`.

## RequestLogging

Loggt jeden Request als JSONL-Zeile:

```jsonl
{"time":"2026-07-24T21:42:12","method":"GET","path":"/","status":200,"ip":"[::1]:56365","duration":"13.584µs"}
```

Felder: `time`, `method`, `path`, `status`, `ip`, `duration`.

Konfiguration:

```json
{ "logging": { "enabled": true } }
```

## Redirect/Rewrite

Wird vor dem Router ausgefuhrt. Redirects leiten um (301/302), Rewrites andern den Pfad transparent.

Konfiguriert in `dreego/config.json` → `redirects` und `rewrites`.

## Plugin-Middleware

Plugins implementieren `MiddlewareProvider` und injizieren eigene Middleware in die Chain. Reihenfolge = `app.Use()`-Reihenfolge.

## Geplant (V1)

- `Recovery`: Panic → 500 (Core-Fixed)
- `RequestID`: X-Request-ID Header (Core-Fixed)
- `RealIP`: X-Forwarded-For Auswertung (Core-Fixed)
- `CSRF`: CSRF-Schutz (Core-Conditional)
