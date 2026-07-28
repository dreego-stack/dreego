# Dreego CLI

## Phase 0 Commands

### `dreego init <path>`
Scaffolds a new Dreego project from embedded blueprints.

```
dreego init myapp

myapp/
├── main.go
└── dreego/
    └── routes/
        └── get.dreego
```

Then: `cd myapp && go mod init myapp && dreego generate && go run .`

### `dreego generate [--force] [--check]`
Transpiles all `.dreego` files. Output: `dreego/gen/routes.go` + `dreego/gen/dree.go`.

Flags:
- `--force` — Force regeneration of all files
- `--check` — CI validation: exit non-zero if .dreego files are newer than gen output (stale)

## Phase 1 Commands (planned)

- `dreego dev` — Dev server with hot reload
- `dreego build` — Production binary

## Phase 2 Commands (planned)

- `dreego add <plugin>` — Install plugin
- `dreego routes` — Show all routes
- `dreego tinker` — Go REPL

## Phase 3 Commands (planned)

- `dreego build --static` — SSG
- `dreego build --wails` — Wails Desktop
- `dreego build --mobile` — Wails Mobile
