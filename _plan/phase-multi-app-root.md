# Phase: Multi-app website root

Breaking refactor of the website model: one root, many apps. See the accepted
decision in [`_docs/decisions/multi-app-website-root.md`](../_docs/decisions/multi-app-website-root.md).

## Target structure

```
myproj/
├─ main.go
└─ dreego/                     # website root (name free; templates use dreego/)
   ├─ dreego.config.json       # REQUIRED marker + root defaults
   ├─ layouts/                 # optional, shared
   ├─ components/              # optional, shared
   ├─ www/                     # app = own Go package
   │   ├─ dreego.config.json   # optional override
   │   ├─ dree.go              # GENERATED: var App dreego.Registrar
   │   ├─ routes/
   │   └─ static/
   └─ blog/
       ├─ routes/
       └─ static/
```

Minimal app: `dreego.config.json` + `www/routes/+page.dreego`.

## Locked decisions

- Root marker: `dreego.config.json` (required). Name free.
- Apps: subdirectories with `routes/`; each its own Go package with an exported
  `var App dreego.Registrar`.
- Shared `layouts/`/`components/` in the root; app-local ones override.
- Config: root defaults, app file overrides field-by-field. host/port in Go.
- Serving: `dreego.New(www.App)` + `go ssr.Listen(app, ":8080")` — no required
  collective helper.
- `core.New(...Registrar)` panics on registration error.
- Hard break: old layout errors with a migration hint.

## Tasks (one PR)

- **T1 core** — `type Registrar func(*App) error`; `New(registrars ...Registrar) *App`
  (panic on error); keep `Register` method. Tests: registrar applied, panic on error.
- **T2 ssr** — keep `Listen`; ensure multiple concurrent hosts + coordinated
  shutdown work and are tested. (No `ServeAll` requirement.)
- **T3 dreefile** — discovery (root = config dir; apps = subdirs with `routes/`;
  shared layouts/components) and generation (per-app `dree.go` with `var App`,
  root `dree.go` for shared only; no collector). Golden tests.
- **T4 config** — root→app field-by-field merge with tests.
- **T5 templates** — `_common/main.go.tmpl`, `web-minimal`, `web-app` on the new
  layout under `dreego/`.
- **T6 migrate** — `demo/demo-ssr`, `demo/demo-wailsv3`, `_tests/fixtures/*`.
- **T7 docs** — ADR, README, `_docs/config.md`, `_docs/routing.md`,
  `_docs/getting-started.md`, `_docs/reference-apps.md`, `CHANGELOG` (minor).
- **T8 integration** — CLI→build→HTTP with two apps on two ports.

## Phase gate

- Black-box integration coverage for the two-app/two-port workflow.
- Generated code compiles through the supported CLI workflow.
- Docs and migration guidance match released behavior.
- Accessibility, security, race, and dependency checks pass.
- `demo/demo-ssr` exercises shared layouts/components across apps.
