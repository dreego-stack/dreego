
---
type: Guide
title: Coding Standards für Dreego
description: Code conventions: file size limits, imports, error handling, testing for Dreego
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Coding Standards für Dreego

## Allgemein

- **Max 120 Zeilen pro Datei** — hard limit
- **Eine logische Sache pro Datei**
- **Keine Kommentare** — Code spricht für sich
- **Package names** kurz, sauber, ohne Hyphen
- **Go 1.26+**, Standard Library bevorzugen

## Build & Run

- **Niemals `go build` direkt** — Build via `make build` oder `dreego build`
- Dev-Server: `dreego dev`
- Tests: `make test`
- Generierte Dateien (`*_dreego.go`) werden nicht committed

## Imports

Standard-Library zuerst, dann Externe, dann Interne. Mit Leerzeilen gruppieren:

```go
import (
    "context"
    "log"

    "github.com/go-chi/chi/v5"

    core "codeberg.org/dreego/dreego/dreego-core"
)
```

## Fehlerbehandlung

- Immer explizit: `if err != nil { return err }`
- Keine `panic()` außer in `init()` und Tests
- `fmt.Errorf` mit `%w` für Wrapping

## Tests

- Testdateien nebem dem zu testenden Code (`foo_test.go`)
- Table-driven Tests bevorzugen
- Test-Fixtures in `testdata/`
- Bei fehlschlagendem Test: Code korrigieren, nicht den Test
