---
type: Decision
title: Monorepo Plugin Layout
description: Superseded decision that previously kept official plugins in this repository
tags: [v0.0.21]
timestamp: 2026-07-31T07:00:00Z
---
# Monorepo Plugin Layout

**Date:** 2026-07-31
**Status:** Superseded — official optional plugins now use separate repositories and modules
**Supersedes:** separate-repo statements in [name-dreego](name-dreego.md) and [routing-and-components](routing-and-components.md)

## Context

Dreego originally planned one repository per plugin (`github.com/dreego-stack/dreego-auth`, `github.com/dreego-stack/dreego-db`, etc.). As a one-person project this is unwieldy: every core update requires bumping versions and re-testing across dozens of separate repos, and the split provides little benefit when there are no external plugin authors yet.

Large frameworks (React, Svelte, Phoenix) keep all official packages in a single repository with independent dependency graphs. Go's equivalent is `go.work` plus per-plugin `go.mod` files.

## Decision

**Official plugins live in this repository under `plugins/`.**

```
dreego/
├── go.work             ← links all local modules for development
├── core/
│   └── go.mod          ← module github.com/dreego-stack/dreego/core (stdlib only)
├── cmd/dreego/
│   └── go.mod          ← module github.com/dreego-stack/dreego/cmd/dreego (requires core)
├── plugins/
│   └── sample/
│       └── go.mod      ← module github.com/dreego-stack/dreego/plugins/sample (requires core)
└── demo/
    └── go.mod          ← module demo (requires core)
```

### Rules

1. **Core never imports a plugin package.** This invariant keeps the root module dependency-free.
2. **Plugins import Core** via `github.com/dreego-stack/dreego/core`.
3. **Plugins with external dependencies get their own `go.mod`** inside `plugins/<name>/`. A `replace` directive points back to the root module.
4. **Dependency-free plugins can be plain packages** in the root module or in `plugins/` with or without their own module.
5. **`go.work` links all local modules** so development works without publishing tags.
6. **One repo, many modules.** Releases for self-contained plugin modules use directory-prefix tags (e.g. `plugins/auth/v0.0.1`).

## Why not the alternatives

| Option | Rejected because... |
|--------|----------------------|
| One repo, single `go.mod` (no per-plugin modules) | Anyone who imports dreego pulls every plugin's transitive dependencies into `go.sum`, even unused ones |
| Separate repo per plugin (original plan) | Unmanageable for a solo maintainer; no external authors yet; version churn across many repos |
| Single `go.mod` + Go 1.17 module pruning | Pruning reduces the build graph but not `go.sum` noise or the release/versioning overhead |

## Relationship to prior decisions

- **name-dreego.md** listed separate `codeberg.org/dreego/<plugin>` repos. This decision supersedes that for official plugins. Community plugins may still live in separate repos as templates.
- **routing-and-components.md** described plugin route/component discovery via `go list -m -json` against external module repos. Under the monorepo model, plugin `.dreego` sources and components live under `plugins/<name>/` and are discoverable by filesystem scan in the same repo.
- **plugin-interface.md** (capability-based contract) is unchanged — it describes the Go interface, not the repository layout.

## Consequences

- A single `git clone` gives access to core + all official plugins.
- One CI pipeline, one issue tracker, one PR-review context.
- Plugin dependency isolation is preserved: importing only core never pulls plugin deps.
- Local development uses `go.work`; consumers see only the modules they import.
- The root module stays dependency-free (stdlib only).
