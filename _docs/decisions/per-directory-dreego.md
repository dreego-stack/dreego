---
type: Decision
title: Per-directory dree.go output
description: dreego/gen/ is gone; each directory with .dreego sources gets its own dree.go, the website root is marked by dreego.config.json
tags: [architecture]
timestamp: 2026-08-21T00:00:00Z
---
# Per-directory dree.go output

**Date:** 2026-08-21
**Status:** Accepted (superseded in part by [One website root with multiple apps](multi-app-website-root.md): the website root now holds one or more app subdirectories, and each app emits `<app>/dree.go` with `var App dreego.Registrar` instead of a root-level `Register(app)`)

## Context

The generated output lived in `dreego/gen/` (routes.go, components.go,
dree.go) and the website was implicitly rooted at a `dreego/` directory with
a `config.json`. The user wants:

- no `dreego/gen/` import in user code — only `import "…/www"` and
  `www.Register(app)`
- the generated output next to its sources: every directory that contains
  `.dreego` files gets its own `dree.go` in that directory
- multiple websites in one repo — the website root must be freely named, so
  a marker file (instead of the fixed `dreego/` name) defines the root

## Decision

- The website root is any directory containing `dreego.config.json`
  (renamed from `config.json`). A root holds apps as immediate subdirectories,
  each with its own `routes/`.
- `dreego generate` walks all website roots and, per app, writes (shown for the
  scaffolded `www` app under `dreego/`):
  - `dreego/www/dree.go` — package www, `var App dreego.Registrar` wiring
    config, static assets and every route-package `Register`
  - `dreego/www/routes/dree.go` — package routes: the root collector that imports
    every sub-route package and calls its `Register`
  - `dreego/www/routes/<folder>/dree.go` — one package per route folder, holding
    that folder's route handlers, layouts are called via
    `layouts.Default(c, pageContent, head)`
  - `dreego/components/dree.go` — shared components package: all component functions
  - `dreego/layouts/dree.go` — shared layouts package: `Default`/`Layout` functions
- Each route folder is its own Go package, so a package-level declaration is
  visible only to the route files in the same folder and must be unique there.
  Shared state that must cross folders lives in a package imported via
  `GOIMPORT`. Dynamic segment directories (`[id]`) and groups (`(group)`) cannot
  be Go packages and fold into their nearest valid ancestor package.
- Components keep their own package; nested component directories use the
  longest valid Go-name prefix as package and are merged into the nearest
  valid parent package.

## Consequences

- User code imports the app package (`myapp/dreego/www` in the scaffolded
  layout) and passes its `App` registrar to `dreego.New(www.App)`; no `gen`
  import exists anymore
- The generated `dree.go` files are gitignored (`dree.go`)
- Multiple websites per repo: each root with `dreego.config.json`
- Integration tests and fixtures were migrated from `dreego/gen` to the
  new layout; the reference apps build and serve
