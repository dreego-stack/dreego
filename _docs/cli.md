# CLI Reference

## dreego generate

```bash
dreego generate [--force]
```

Transpiles `.dreego` files in `dreego/routes/` to Go code. Creates `dree.go` per route directory and `dreego/gen/dree.go` as central import file.

- `--force`: Forces complete regeneration (ignores cache)

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

## dreego help

```bash
dreego help
dreego --help
```

Shows all available commands and flags.

> **Note:** `dreego build` and `dreego run` are dev tools, not for production.
