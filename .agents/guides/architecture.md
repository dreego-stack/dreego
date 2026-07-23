# Architektur-Guide

## Projektstruktur

```
dreego/
├── cmd/
│   └── dreego/
│       └── main.go              # CLI Entry Point (max 120 Zeilen)
├── pkg/
│   ├── transpiler/               # Lexer, Parser, AST, Code-Generator
│   │   ├── lexer.go
│   │   ├── parser.go
│   │   ├── ast.go
│   │   └── codegen.go
│   ├── router/                   # Chi-Wrapper, File-based Routing
│   │   ├── router.go
│   │   └── routes.go
│   ├── plugin/                   # Plugin-Interface & Registry
│   │   └── plugin.go
│   ├── context/                  # Request-Context, Session
│   │   └── context.go
│   └── middleware/                # CSRF, Session, Auth
│       └── middleware.go
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
| transpiler  | .dreego → Go-Code                                | Keine externen       |
| router      | HTTP-Routing, Chi-Integration                    | chi, transpiler      |
| plugin      | Plugin-Interface, Registry                       | Keine                |
| context     | Request-Kontext, User, Session                   | net/http             |
| middleware   | CSRF, Auth, Session, Logging                     | context              |

## Regeln

- **Jedes Package ist eigenständig testbar**
- **Keine zirkulären Abhängigkeiten**
- **`internal/` für Implementierungsdetails, die nicht Teil der Public API sind**
- **`pkg/` für stabile, öffentliche APIs**
- **`cmd/` nur für Einstiegspunkte**
