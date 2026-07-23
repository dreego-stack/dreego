# Dreego CLI

## Phase 0 Commands

### `dreego init <name>`
Scaffolding. Erstellt ein neues Dreego-Projekt mit allem was man zum Starten braucht.

```
dreego init myapp

myapp/
├── routes/
│   └── index.dreego       ← Hello-World Startseite
├── main.go                 ← Minimaler Server
├── go.mod                  ← codeberg.org/dreego/dreego
└── .gitignore
```

Danach: `cd myapp && dreego generate && go run .` → Server läuft.

### `dreego generate`
Transpiliert alle `.dreego`-Dateien in `routes/` zu `_dreego.go`.

```
routes/
├── index.dreego       → index_dreego.go
├── about.dreego       → about_dreego.go
```

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
