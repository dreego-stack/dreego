# CLI Reference

## dreego new

```bash
dreego new <name> [-t <template>] [-l|--list]
```

Scaffolds a new project from the `web-minimal` template in a new directory. It
writes `main.go`, `Taskfile.yml`, `.gitignore`, and the `www/` tree, then runs
`go mod init` and `go mod tidy`.

- `-t <template>`, `--template <template>`: select a template; `web-minimal` is
  the default, and `web-app` is the other shipped template
- `-l`, `--list`: list the available templates and exit without scaffolding

An unknown template name fails with an error that lists the valid names. A
missing `-t` value also fails. Both exit non-zero.

## dreego init

```bash
dreego init <path> [-t <template>] [-l|--list]
```

Scaffolds the `web-minimal` template into an existing or new path. It accepts
the same `-t`/`--template` and `-l`/`--list` flags as `dreego new`, with
`web-minimal` as the default.

## dreego task

```bash
dreego task [args...]
```

Forwards its arguments to the external `task` binary and forwards its exit
code. With no arguments it runs `task --list`. Dreego does not embed a task
runner; when `task` is not on `PATH`, the command fails with an error pointing
to the Task installation page.

## dreego generate

```bash
dreego generate [--force] [--check]
```

Transpiles `.dreego` files in the website root (any directory with a `dreego.config.json`) to Go code. Produces one `dree.go` per directory with sources, plus `dree.go` at the root (config + static assets + Register). Files are only written when content changes.

- `--force`: Forces complete regeneration (ignores cache)
- `--check`: CI mode — regenerates the expected output in memory and compares it byte-for-byte against the files on disk. No working-tree modification. Exits non-zero with a path-level diff (`missing:`, `extra:`, `stale:`) when any generated file (routes, components, layouts, static assets, config) is missing, extra, or stale. Timestamp manipulation cannot produce a false pass.

## dreego i18n extract

```bash
dreego i18n extract
```

Emits deterministic, versioned JSON for translation-management adapters. The
stream contains the default-locale messages, argument and formatting contracts,
structured variants, and configured target locales.

## dreego tools install typescript

```bash
dreego tools install typescript
```

Downloads and verifies Dreego's pinned native TypeScript compiler for the
current platform. Installation is explicit; `generate` and `build` never
download tools automatically. Projects without TypeScript sections do not need
this tool. See [TypeScript Client Code](client-typescript.md).

## dreego build

```bash
dreego build [--target <os/arch>] [--yes]
```

Runs `generate`, approved plugin build hooks, and `go build`. The binary lands
in `build/bin/<name>`.

- `--target <os/arch>` cross-compiles for a target such as `linux/amd64` or
  `darwin/arm64` and includes the target in the output name.
- `--yes` approves every declared plugin build hook for this run. Without prior
  approval, interactive builds ask first and non-interactive builds fail safely.

See [Build Hooks](build-hooks.md) for declarations, approvals, and CI behavior.

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

Reads documentation from the **local module store** — no HTTP, no embedded copy. Dreego's own documentation lives in one central `_docs/` tree at the repository root, so a feature is documented in a single place. External plugins keep their own `_docs/`.

Resolution follows Go itself: the current module, workspaces, replacements,
vendor trees, and downloaded modules are resolved with `go list -m`. For the
installed CLI and its dependencies, build information supplies the exact
version to `go mod download -json` when the project does not require that
module directly.

Without arguments, `dreego docs` shows the documentation index `/_docs/index.md`. Flags:

- `-p <name>`: read an external plugin's docs, such as `plugin-sse`
- `--list`: list the documentation index plus every plugin's `_docs/sitemap.json`
- `--dump`: print all sitemap pages (or a comma-separated list of paths) in one output
- `--json`: structured JSON (headings, code blocks, links) for AI agents
- `--web`: open the docs page in a browser

Examples:
```bash
dreego docs                    show the documentation index
dreego docs /README.md         show the readme
dreego docs /_docs/cli.md      show the CLI reference
dreego docs -p plugin-sse /_docs/index.md   show a plugin's docs
dreego docs --list             list the index and plugin pages
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
