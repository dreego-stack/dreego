# CLI Reference

## dreego generate

```bash
dreego generate [--force]
```

Transpiliert `.dreego`-Dateien in `dreego/routes/` zu Go-Code. Erzeugt `dree.go` pro Route-Verzeichnis und `dreego/gen/dree.go` als zentrale Import-Datei.

- `--force`: Erzwingt vollstandige Regeneration (ignoriert Cache)

## dreego build

```bash
dreego build
```

Fuhrt `generate` aus, dann `go build`. Das Binary landet unter `build/bin/<name>`.

## dreego run

```bash
dreego run [-d] [-t <seconds>]
```

Fuhrt `build` aus und startet den Server.

- `-d`: Debug-Modus. Schreibt Request-Logs (JSONL) nach `build/logs/<utc>.log`
- `-t <seconds>`: Timer. Server stoppt automatisch nach N Sekunden

### Beispiele

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

Zeigt alle verfugbaren Commands und Flags.

> **Hinweis:** `dreego build` und `dreego run` sind Dev-Tools, nicht fur Production.
