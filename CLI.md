# Dreego CLI

## Phase 0 Commands

### `dreego init <path>`
Scaffolding. Erstellt ein neues Dreego-Projekt aus embedded Blueprint.

```
dreego init myapp

myapp/
├── main.go
└── dreego/
    └── routes/
        └── get.dreego
```

Danach: `cd myapp && go mod init myapp && dreego generate && go run .`

### `dreego generate`
Transpiliert alle `.dreego`-Dateien. Output: `dreego/gen/routes.go` + `dreego/gen/dree.go`.

## Phase 1 Commands (geplant)

- `dreego dev` — Dev-Server mit Hot Reload (air)
- `dreego build` — Production Binary

## Phase 2 Commands (geplant)

- `dreego add <plugin>` — Plugin installieren
- `dreego routes` — Alle Routen anzeigen
- `dreego tinker` — Go-REPL

## Phase 3 Commands (geplant)

- `dreego build --static` — SSG
- `dreego build --wails` — Wails Desktop
- `dreego build --mobile` — Wails Mobile
