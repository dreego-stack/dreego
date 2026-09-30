---
type: Decision
title: Global defaults with local overrides
description: Root routes/, static/, layouts/, and components/ are shared defaults for every app; an app-local file with the same relative path wins
tags: [architecture, routing, codegen, static, layouts]
timestamp: 2026-09-30T00:00:00Z
---
# Global defaults with local overrides

**Date:** 2026-09-30
**Status:** Accepted

## Context

The multi-app website root introduced shared `layouts/`/`components/`, but kept
`routes/` and `static/` strictly per app: a root-level `routes/` was rejected and
every app had to repeat the common routes. The configuration model already works
differently: the root `dreego.config.json` provides defaults and an app file
overrides field by field. This decision makes the tree layout follow the same
instinct — **the root is the default, an app overrides**.

## Decision

- The website root may own `routes/`, `static/`, `layouts/`, and `components/`.
  All four are **global defaults** inherited by every app.
- Merging is **per relative file path**, not per directory. An app-local file
  with the same relative path wins; every other global file stays available.
  A local `static/favicon.ico` shadows the global one, while
  `static/shared.css` remains inherited.
- `layouts/` and `components/` are imported by explicit path, so a global and an
  app-local file coexist under one Go package unless they collide. The app-local
  default layout still wins in the cascade; a route may select a layout by path
  with `LAYOUT "…"`.
- An app is any immediate subdirectory carrying one of those trees or its own
  `dreego.config.json`. Root-owned directory names (`routes`, `static`,
  `layouts`, `components`, `locales`) are never apps. A config directory below
  another config directory is an app, not a new root.
- A `LAYOUT` directive names a layout path and is resolved against the website
  root and the working directory. `layouts/` may hold explicitly named layouts
  (for example `admin.dreego`) beside the default; `default.dreego` and
  `layout.dreego` in one directory remain an ambiguity error.

## Rationale

- The tree layout now mirrors the config model ("root default, app override")
  instead of contradicting it.
- Per-file shadowing lets an app replace a single route or asset without copying
  the whole shared tree, and without hiding global files it still wants.
- Explicit component and layout imports keep both variants addressable; nothing
  is silently dropped.

## Consequences

- Root-level `routes/` is no longer an error; the legacy-layout migration hint
  for that case is removed.
- Global route-local layouts live under `dreego/routes/<sub>/layouts/` and use
  the `root/<sub>` scope (`DefaultRootAdmin`).
- `generate`/discovery merge global and local trees; new unit and integration
  tests cover inheritance, per-file override, and app isolation.

## See also

- [One website root with multiple apps](multi-app-website-root.md)
- [Routing](../routing.md), [Layouts](../layouts.md)
