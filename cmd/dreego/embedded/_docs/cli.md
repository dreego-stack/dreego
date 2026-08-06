# CLI Reference

## dreego generate

```bash
dreego generate [--force]
```

Transpiles `.dreego` files in `dreego/routes/` and `dreego/components/` to Go code. Produces split-gen output: `gen/routes.go` + `gen/components.go` + `gen/dree.go` (config + static assets). Files are only written when content changes.

- `--force`: Forces complete regeneration (ignores cache)
- `--check`: CI mode — exits non-zero when generated files are stale

## dreego build

```bash
dreego build
```

Runs `generate`, then `go build`. The binary lands in `build/bin/<name>`.

## dreego run

```bash
dreego run [-d] [-t <seconds>]
```

Runs `build` and starts the server.

- `-d`: Debug mode. Writes request logs (JSONL) to `build/logs/<utc>.log`
- `-t <seconds>`: Timer. Server stops automatically after N seconds

### Examples

```bash
dreego run                  # build + start (foreground)
dreego run -d               # build + start + log to file
dreego run -t 60            # build + start + stop after 60s
dreego run -d -t 60         # debug log + 60s timer
```

## dreego dev

```bash
dreego dev
```

Runs `generate` + `build`, starts the server, then watches `.dreego` files (500 ms poll). On any change it regenerates, rebuilds, and gracefully restarts the server (SIGTERM + reap). Build errors do **not** kill the watcher — the previous server keeps running. Stop with `Ctrl-C`.

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.

## dreego docs

```bash
dreego docs [path]
```

Displays repo documentation from the embedded copy (`cmd/dreego/embedded/`), so it works **offline**. Without arguments, shows `/_docs/index.md`. Local plugin docs in `plugins/<name>/_docs/` take priority over the embedded copy.

Examples:
```bash
dreego docs                    show docs index
dreego docs /README.md         show readme
dreego docs /_docs/cli.md      show CLI docs
```

## dreego help

```bash
dreego help
dreego --help
```

Shows all available commands and flags.

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.

## See Also

- [Docs Index](https://codeberg.org/dreego/dreego/src/branch/main/_docs/index.md)
- [Getting Started](https://codeberg.org/dreego/dreego/src/branch/main/_docs/getting-started.md)
