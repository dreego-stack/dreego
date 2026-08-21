# Testing Strategy

Tests run as Go integration tests in `_tests/go/` via Docker (`make test`), using the `dreegotest` helpers. Every area covers positive (confirm behavior) and negative (detect errors early) cases.

## Test Layout

- `internal/transpiler/*_test.go` — unit tests for the lexer, parser, and codegen (run with `go test ./internal/transpiler/...`).
- `core/*_test.go` and `core/internal/*/*_test.go` — unit tests for the runtime framework facade and its internal packages (run with `go test ./core/...`).
- `_tests/go/*_test.go` — integration tests that build a real project, run the CLI, and assert on generated code and HTTP behavior.
- `dreegotest/` — shared helpers: `ProjectDir`, `RunCLI`, `Build`, `MustBuild`, `NewApp`, `RenderComponent`.

## Areas Covered

### Transpiler
Basic page with all sections, routes without `<go>`, unclosed `<div>`, mismatched closing tags, XSS escaping, output contexts (text, attribute, URL, script, style), duplicate sections, empty templates, large templates, unicode, comments, verbatim blocks.

### Template Expressions
`{#if}` true/false, `{#each}` loops, nested control flow, `{#else}` and `{#else if}`, empty lists, expressions with missing variables (build-time failure), function expressions, filters.

### Layout
`{#slot}` and `{#head}` merging, routes without a layout, layout application bugs (regression), route head merging.

### Routing
GET/POST/PUT/DELETE, dynamic `[id]` segments, catch-all `[...path]`, invisible `(group)/` segments, custom 404 and 500 pages, multi-segment and deep nesting.

### Middleware
Recovery, CSRF (token issue, POST validation, disable), gzip compression, security headers, health and readiness, request logging, request ID.

### Session
Set and read values, delete, destroy, cookie store setup, encryption (AES-256-GCM).

### Components
Props, self-closing calls, default and named slots, scoped CSS, nested components, expression props, named-prop contract checking, import aliases.

### CLI
`init`, `new`, `generate` (including `--force` and `--check`), `build`, `run`, `dev`, `docs`, `fmt`, `version`. Stale detection and no-argument help. Accessibility of CLI output (no color, screen-reader-linear help, actionable error format).

### Config
`dreego.config.json` redirects, rewrites, logging toggle, invalid JSON.

### Form Actions
`<form g-action>` generation, int/bool binding, validation, PRG redirect, error re-render with `c.Errors` and `c.Old`.

### Bugs (Regression)
Every fixed bug keeps a regression test in `_tests/go/bug_*_test.go`, `core/*_test.go`, `core/internal/*/*_test.go`, or `internal/transpiler/*_test.go`.

## Accessibility Tests

- CLI output is color-free and screen-reader-linear (`_tests/go/cli_accessibility_test.go`).
- Generator diagnostics lead with `file:line:col`, the cause, and a practical `Fix:` action.
- The landing blueprint uses semantic HTML (`<main>`, `<nav>`, skip link, `{#slot}`) and gives every `<img>` an `alt`. The minimal `init` blueprint is tested as a minimal route, not as a complete accessible application shell.
- The transpiler emits a11y diagnostics for missing image alternatives and unassociated form labels (`internal/transpiler/a11y_check_test.go`).

## Running Tests

```bash
make test                              # Docker-based full suite
go test ./core/...                     # runtime unit tests only
go test ./internal/transpiler/...     # transpiler unit tests only
go test ./_tests/go/ -parallel 1 -p 1  # integration tests (no parallelism for CLI builds)
```

## See Also

- [Accessibility](accessibility.md) — Framework accessibility guarantees
- [CLI](cli.md) — CLI reference
- [Getting Started](getting-started.md) — Tutorial
