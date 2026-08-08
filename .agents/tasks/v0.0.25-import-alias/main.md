# Task v0.0.25-import-alias — Rename `core` → `dreego` import alias

## Goal

Make the Dreego core package import alias consistent across the entire repo:

```go
import dreego "codeberg.org/dreego/dreego/core"
```

replacing the previous `core "codeberg.org/dreego/dreego/core"` everywhere, so
`dreego.UsePlugin`, `dreego.Listen`, etc. read naturally. Pure rename — no
behavior change.

## Files changed by category

### 1. Generator (source of truth)
The code generator emits the import alias into `gen/*.go` files. Changed it so
all generated output uses `dreego`:

- `core/generate.go` — import emission for `routes.go`, `components.go`,
  `dree.go`; emitted `SetLogging`, `RegisterRedirect`, `RegisterRewrite`,
  `RegisterStatic` calls.
- `core/codegen.go` — emitted `SSRContext`, `NewSSR`, `Register`,
  `SetErrorHandler`, `Component`, `ComponentFunc`.
- `core/codegen_form.go` — emitted `NewSSR`, `BindForm`, `ValidateForm`,
  `SaveErrors`, `SaveOld`, `ErrRedirect`.
- `core/codegen_template.go` / `core/codegen_component.go` — emitted `EachLoop`.

### 2. Golden testdata
- `core/testdata/golden/simple_route.golden`
- `core/testdata/golden/route_with_layout.golden`
- `core/testdata/golden/component_with_style.golden`

### 3. Generated/check-stale fixtures
- `_tests/core/CLI/check-stale/dreego/gen/routes.go`
- `_tests/core/CLI/check-stale/dreego/gen/components.go`
- `_tests/core/CLI/check-stale/dreego/gen/dree.go`

### 4. Blueprints + demo
- `cmd/dreego/blueprints/default/main.go.tmpl`
- `cmd/dreego/blueprints/landing/main.go.tmpl`
- `demo/main.go`
- root `main.go` (repo scratch demo)

### 5. CLI internals
- `cmd/dreego/main.go`
- `cmd/dreego/dev.go`
- `cmd/dreego/fmt.go`

### 6. Test scripts (59 heredoc Go snippets)
All `_tests/core/**/test.sh` files that embed a Go snippet importing core:
import alias + all `core.` uses → `dreego.`. Shell logic untouched. One test
(`Routing/servemux-cache/test.sh`) used a non-aliased import; it now uses the
`dreego` alias explicitly.

### 7. dreegotest + plugins
- `dreegotest/request.go`, `dreegotest/render.go`,
  `dreegotest/dreegotest_test.go`
- `plugins/sample/sample.go`

### 8. Documentation
- `_docs/runtime.md`, `_docs/getting-started.md`, `_docs/plugins.md`,
  `_docs/deployment.md`, `_docs/frontmatter.md`, `_docs/session-encryption.md`,
  `_docs/forms.md`, `_docs/middleware.md`, `_docs/components.md`
- `README.md`
- `.agents/guides/coding-standards.md`
- `cmd/dreego/embedded/**` — regenerated via `_scripts/sync-embedded-docs.sh`

## Not changed

- `core/` package internals (`package core` — no import alias used).
- Historical `CHANGELOG.md` prose that refers to the `core` package name
  (not the import alias).

## Verification

- `go build ./core/...` — ok
- `go test ./core/...` — ok
- `go build ./cmd/dreego` — ok
- `go test ./cmd/dreego/...` — ok
- `go test ./dreegotest/...` — ok
- Full suite `sh _tests/test.sh` — **164 Passed / 0 Failed**
- `run-timer-sigterm` (previous known flake) passed both in full run and
  isolated.
- grep confirms `core "codeberg.org/dreego/dreego/core"` no longer appears.

## Result

All generated, hand-written, blueprint, CLI, test, and documentation code now
consistently uses the `dreego` import alias for the core package.
