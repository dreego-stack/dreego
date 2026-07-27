
---
type: Guide
title: Architektur-Guide
description: Project structure, module boundaries, and architectural rules for Dreego
tags: [v0.0.1]
timestamp: 2026-07-23T00:00:00Z
---
# Architektur-Guide

## Projektstruktur

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go              # CLI Entry Point (max 120 Zeilen)
├── dreego-core/                  # Core library (single package)
│   ├── lexer.go
│   ├── parser.go
│   ├── ast.go
│   ├── codegen.go
│   ├── router.go
│   ├── routes.go
│   ├── plugin.go
│   ├── context.go
│   └── middleware.go
├── dreego-plugin/                # Plugins (future)
├── internal/                     # Nicht-öffentliche Pakete
│   └── ...
├── testdata/                     # Test-Fixtures (.dreego-Dateien)
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── .kilo/
```

## Modul-Grenzen

| Modul       | Verantwortung                                    | Abhängigkeiten        |
|-------------|-------------------------------------------------|----------------------|
| core        | .dreego → Go-Code, Router, Context, Middleware    | net/http, chi        |
| plugin      | Plugin-Interface, Registry (future)              | core                 |

## Regeln

- **Jedes Package ist eigenständig testbar**
- **Keine zirkulären Abhängigkeiten**
- **`internal/` für Implementierungsdetails, die nicht Teil der Public API sind**
- **`dreego-core/` für stabile, öffentliche APIs**
- **`dreego-plugin/` für Plugins (future)**
- **`cmd/` nur für Einstiegspunkte**
