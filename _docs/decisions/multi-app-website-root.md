---
type: Decision
title: One website root with multiple apps
description: A website root holds shared layouts/components and multiple app subdirectories; each app is its own Go package with a Registrar, and hosts/ports live in main.go
tags: [architecture, routing, codegen, cli]
timestamp: 2026-09-30T00:00:00Z
---
# One website root with multiple apps

**Date:** 2026-09-30
**Status:** Accepted (breaking change, pre-1.0)

## Context

Today a *website* is any directory containing `dreego.config.json`; it owns
`routes/`, `layouts/`, `components/`, and `static/` (`internal/dreefile/discovery.go`).
Several sites in one module are therefore several roots, and `main.go` wires them
by hand — `demo/demo-ssr/main.go` builds a `hostRouter` that dispatches on the
`Host` header and path prefix. There is no shared `layouts/`/`components/` across
sites, one root maps to exactly one `Register(app)` entry point, and the server
listens on a single global port.

The new model is **one root, many apps**: a shared shell for a brand with several
subdomain surfaces (www, blog, app, …), explicit Go ownership of hosts and ports,
and per-app Go packages. It builds on the released per-folder route packages and
generated-file markers rather than replacing them.

## Decision

- **Website root** = a directory containing `dreego.config.json` (the marker).
  Its name is free; templates and examples use `dreego/`.
- The root holds **optional, shared** `layouts/` and `components/`.
- **Apps are subdirectories** that contain `routes/` (and optionally `static/`,
  their own `dreego.config.json`, and app-local `layouts/`/`components/`). Each
  app is its **own Go package**; `dreego generate` emits `<app>/dree.go` with
  `var App dreego.Registrar`.
- The **minimal app** is `dreego.config.json` + `<app>/routes/+page.dreego`.
  Nothing else is required — no layouts, no components.
- **Shared** `layouts/`/`components/` generate into the root's `layouts` and
  `components` packages; **app-local** `layouts/`/`components/` generate into
  their own uniquely named packages and take precedence for that app.
- **Config**: the root `dreego.config.json` provides defaults; an app's
  `dreego.config.json` overrides **field by field** (no deep merge). `host` and
  `port` are **not** config — they live in Go.
- **Hosts and ports live in `main.go`** and are wired explicitly:

  ```go
  wwwApp := dreego.New(www.App)
  blogApp := dreego.New(blog.App)
  go ssr.Listen(wwwApp, ":8080")
  err := ssr.Listen(blogApp, ":8081")
  ```

  `core` gains `type Registrar func(*App) error` and `New(registrars ...Registrar) *App`.
  `New` **panics** on a registration error, because route, static, and
  configuration conflicts (including invalid or looping redirect/rewrite rules)
  are already validated by `dreego generate`; a runtime registration failure
  cannot be produced by a correctly generated app.
- **Breaking**: the old layout (root-level `routes/`, multiple roots, a `main.go`
  that imports `<root>` and calls `Register`) is rejected with a clear error and
  a migration hint. No silent behavior, no transition period.

## Rationale

- One import per app; two apps are two Go packages — no artificial "collector"
  package whose only job would be to forward registrations.
- Shared `layouts/`/`components/` directly model "one brand, many surfaces".
- Ports/hosts in Go are typed and explicit and let an external proxy map
  `blog.example.com → :8081`.
- JSON config stays; no new dependency (TOML would add one to a module that is
  standard-library + `golang.org/x/` only).

## Consequences

- `generate`/discovery rewrite; new golden tests; the package-name rule
  (`sanitizePkgName`) now applies to app directories too, and it preserves `_` so
  path-derived package names cannot collide.
- CLI templates (`_common/main.go.tmpl`, `web-minimal`, `web-app`), demos,
  `_tests/fixtures`, docs, and `CHANGELOG` are updated; release is a **minor**
  bump (pre-1.0 breaking change).
- Several website roots per repo remain allowed (several `dreego.config.json`).
- `DREEGO_PORT` still applies to single-app projects; with several apps the ports
  from `main.go` win.

## Alternatives considered

- **`NAME.routes/` + `NAME.static/` siblings** — scattered; two directories per app.
- **Config-declared apps** — two places to maintain for one app.
- **A single collector package** — rejected; apps are separate packages.
- **`ssr.ServeAll` as the required start form** — rejected; plain `go ssr.Listen`
  is the norm. A collective helper may be offered later, never required.
- **TOML config** — rejected; JSON avoids a new dependency.

## Open / not built

- `dreego migrate` command (deferred).
- In-Go host dispatch helper (deferred).
- Transition period with warnings (rejected; hard break with a hint).
