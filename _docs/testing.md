# Testing Strategy

The suite has three layers:

1. **Unit tests** next to the code — `internal/dreefile/**/*_test.go`,
   `core/**/*_test.go`, `adapter/**/*_test.go`, `dreegotest/*_test.go`,
   `cmd/dreego/**/*_test.go`. Run with `go test ./…`.
2. **End-to-end tests** under `_tests/e2e/` — they build a real project, run the
   CLI, and assert on generated code and HTTP behavior.
3. **Invariant gates** under `_tests/e2e/invariants/` — repo-wide rules
   (layering, dependencies, no binaries, coverage, fuzzing) that used to be
   separate CI steps.

## Test Layout

`_tests/` is one recursive tree. One test lives in one folder. A folder holds a
`test.sh`, a Go `<name>_test.go`, or both. See `_tests/how-to-test.md` for the
contract.

```
_tests/
  test.sh                 recursive orchestrator (the only entry point)
  run-one.sh              runs a single folder
  how-to-test.md          the contract for every test
  e2e/
    <group>/<test-name>/   the test folder
      test.sh              shell test (optional)
      <name>_test.go       Go test, package tests (optional)
      www/…                real input files the test needs, if any
```

Go and shell tests are equal: both are discovered the same way, both run with
their own folder as the working directory, and both signal success with exit
code 0.

## Areas Covered

### Transpiler
Basic page with all sections, routes without `<server>`, unclosed `<body>`, mismatched closing tags, XSS escaping, output contexts (text, attribute, URL, script, style), duplicate sections, empty templates, large templates, unicode, comments, verbatim blocks.

### Template Expressions
`{#if}` true/false, `{#each}` loops, nested control flow, `{#else}` and `{#else if}`, empty lists, expressions with missing variables (build-time failure), function expressions, filters.

Dedicated i18n tests cover `[[ message.key ]]` expressions, compile-time key
and argument validation, catalog parity, locale negotiation, plural and select
variants, currency display without conversion, and HTML language metadata.

### Layout
`{#slot}` and `{#head}` merging, routes without a layout, route-local layout
cascades, layout application bugs (regression), route head merging.

### Routing
GET/POST/PUT/DELETE, dynamic `[id]` segments, catch-all `[...path]`, invisible `(group)/` segments, custom 404 and 500 pages, multi-segment and deep nesting, one Go package per route folder.

### Middleware
Recovery, CSRF (token issue, POST validation, disable), gzip compression, security headers, health and readiness, request logging, request ID.

### Session
Set and read values, delete, destroy, cookie store setup, encryption (AES-256-GCM).

### Components
Props, self-closing calls, default and named slots, scoped CSS, nested components, expression props, named-prop contract checking, import aliases.

### CLI
`new`, `generate` (including `--force` and `--check`), `build`, `run`, `dev`, `docs`, `fmt`, `version`. Stale detection and no-argument help. Accessibility of CLI output (no color, screen-reader-linear help, actionable error format).

### Config
`dreego.config.json` redirects, rewrites, logging toggle, invalid JSON.

### Form Actions
`<form g-action>` generation, int/bool binding, validation, PRG redirect, error re-render with `c.Errors` and `c.Old`.

### Bugs (Regression)
Every fixed bug keeps a regression test in its own `_tests/e2e/bugs/<name>/`
folder (or a unit test next to the code for very small fixes).

### Invariants (`_tests/e2e/invariants/`)
Repo-wide gates, each a normal `test.sh`:

- `layering` — the internal layering and dreefile dependency rule.
- `core-deps` — root, core, and the adapters use only the standard library and `golang.org/x/`.
- `no-binaries` — no unexpected binary files in the repo.
- `coverage` — core coverage gate (35% minimum).
- `fuzz` — the Markdown-to-HTML renderer stays safe under fuzzing.

## Accessibility Tests

- CLI output is color-free and screen-reader-linear (`_tests/e2e/cli/cli_accessibility/`).
- Generator diagnostics lead with `file:line:col`, the cause, and a practical `Fix:` action.
- The `web-minimal` template layout ships `<html lang="en">`, a skip link, and a `<main id="main">` landmark. The route is tested as a minimal page, not as a complete accessible application shell.
- The compiler emits a11y diagnostics for missing image alternatives and unassociated form labels (`internal/dreefile/a11y_check_test.go`).

## Running Tests

```bash
task test                          # unit tests + the whole _tests/e2e tree
sh _tests/test.sh                  # only the _tests/e2e tree (recursive)
DREEGO_FILTER=layouts sh _tests/test.sh   # only matching folders
DREEGO_RACE=0 sh _tests/test.sh    # skip -race for Go folders
go test ./internal/dreefile/...    # compiler unit tests only
go test ./core/...                 # runtime unit tests only
```

`_tests/test.sh` builds the CLI once, installs the pinned TypeScript compiler
once, then runs every folder in parallel (`DREEGO_JOBS`, default `nproc`).

## See Also

- [Accessibility](accessibility.md) — Framework accessibility guarantees
- [CLI](cli.md) — CLI reference
- [Getting Started](getting-started.md) — Tutorial
