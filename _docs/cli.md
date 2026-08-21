# CLI Reference

## dreego generate

```bash
dreego generate [--force] [--check]
```

Transpiles `.dreego` files in the website root (any directory with a `dreego.config.json`) to Go code. Produces one `dree.go` per directory with sources, plus `dree.go` at the root (config + static assets + Register). Files are only written when content changes.

- `--force`: Forces complete regeneration (ignores cache)
- `--check`: CI mode — regenerates the expected output in memory and compares it byte-for-byte against the files on disk. No working-tree modification. Exits non-zero with a path-level diff (`missing:`, `extra:`, `stale:`) when any generated file (routes, components, layouts, static assets, config) is missing, extra, or stale. Timestamp manipulation cannot produce a false pass.

## dreego build

```bash
dreego build
```

Runs `generate`, then `go build`. The binary lands in `build/bin/<name>`.

## dreego run

```bash
dreego run [-d]
```

Runs `build` and starts the server.

- `-d`: Debug mode. Writes request logs (JSONL) to `build/logs/<utc>.log`

### Examples

```bash
dreego run                  # build + start (foreground)
dreego run -d               # build + start + log to file
```

## dreego dev

```bash
dreego dev
```

Runs `generate` + `build`, starts the server, then watches `.dreego` files (500 ms poll). On any change it regenerates, rebuilds, and gracefully restarts the server (SIGTERM + reap). Build errors do **not** kill the watcher — the previous server keeps running. Stop with `Ctrl-C`.

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.

## dreego docs

```bash
dreego docs [-p <name>] [--web] [--json] [--dump] [--list] [path]
```

Reads documentation from the **local module store** — no HTTP, no embedded copy. The docs live next to the source in each module's `_docs/` directory, so there is a single source of truth per module.

Resolution uses the project's `go.mod` to locate each module on disk, in priority order:
1. the module itself (when run inside the dreego repo)
2. the `vendor/` directory (when present)
3. the Go module cache (`go env GOMODCACHE`)

Without arguments, `dreego docs` shows the core index `/_docs/index.md`. Flags:

- `-p <name>`: read a plugin's docs from `github.com/dreego-stack/<name>` (must be required in `go.mod`)
- `--list`: list every core + plugin page from each module's `_docs/sitemap.json`
- `--dump`: print all sitemap pages (or a comma-separated list of paths) in one output
- `--json`: structured JSON (headings, code blocks, links) for AI agents
- `--web`: open the docs page in a browser

Examples:
```bash
dreego docs                    show core docs index
dreego docs /README.md         show core readme
dreego docs /_docs/cli.md      show core CLI docs
dreego docs -p plugin-sse /_docs/index.md   show a plugin's docs
dreego docs --list             list all core + plugin pages
```

> **Note:** `dreego docs` reads the version installed in your project's `go.mod`. If a module is not downloaded yet, run `go mod download` first.

## dreego help

```bash
dreego help
dreego --help
```

Shows all available commands and flags.

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.

## See Also

- [Docs Index](https://github.com/dreego-stack/dreego/blob/main/_docs/index.md)
- [Getting Started](https://github.com/dreego-stack/dreego/blob/main/_docs/getting-started.md)
